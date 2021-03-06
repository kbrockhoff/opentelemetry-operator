package opentelemetrycollector

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/open-telemetry/opentelemetry-operator/pkg/apis/opentelemetry"
	"github.com/open-telemetry/opentelemetry-operator/pkg/apis/opentelemetry/v1alpha1"
)

// reconcileDeployment reconciles the deployment(s) required for the instance in the current context
func (r *ReconcileOpenTelemetryCollector) reconcileDeployment(ctx context.Context) error {
	desired := deployments(ctx)

	// first, handle the create/update parts
	if err := r.reconcileExpectedDeployments(ctx, desired); err != nil {
		return fmt.Errorf("failed to reconcile the expected deployments: %v", err)
	}

	// then, delete the extra objects
	if err := r.deleteDeployments(ctx, desired); err != nil {
		return fmt.Errorf("failed to reconcile the deployments to be deleted: %v", err)
	}

	return nil
}

func deployments(ctx context.Context) []*appsv1.Deployment {
	instance := ctx.Value(opentelemetry.ContextInstance).(*v1alpha1.OpenTelemetryCollector)

	var desired []*appsv1.Deployment
	if len(instance.Spec.Mode) == 0 || instance.Spec.Mode == opentelemetry.ModeDeployment {
		desired = append(desired, deployment(ctx))
	}

	return desired
}

func deployment(ctx context.Context) *appsv1.Deployment {
	instance := ctx.Value(opentelemetry.ContextInstance).(*v1alpha1.OpenTelemetryCollector)
	logger := ctx.Value(opentelemetry.ContextLogger).(logr.Logger)
	name := resourceName(instance.Name)

	image := instance.Spec.Image
	if len(image) == 0 {
		image = viper.GetString(opentelemetry.OtelColImageConfigKey)
	}

	labels := commonLabels(ctx)
	labels["app.kubernetes.io/name"] = name

	specAnnotations := instance.Annotations
	if specAnnotations == nil {
		specAnnotations = map[string]string{}
	}

	specAnnotations["prometheus.io/scrape"] = "true"
	specAnnotations["prometheus.io/port"] = "8888"
	specAnnotations["prometheus.io/path"] = "/metrics"

	argsMap := instance.Spec.Args
	if argsMap == nil {
		argsMap = map[string]string{}
	}

	if _, exists := argsMap["config"]; exists {
		logger.Info("the 'config' flag isn't allowed and is being ignored")
	}

	// this effectively overrides any 'config' entry that might exist in the CR
	argsMap["config"] = fmt.Sprintf("/conf/%s", opentelemetry.CollectorConfigMapEntry)

	var args []string
	for k, v := range argsMap {
		args = append(args, fmt.Sprintf("--%s=%s", k, v))
	}

	configMapVolumeName := fmt.Sprintf("otc-internal-%s", name)
	volumeMounts := []corev1.VolumeMount{{
		Name:      configMapVolumeName,
		MountPath: "/conf",
	}}
	volumes := []corev1.Volume{{
		Name: configMapVolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: name},
				Items: []corev1.KeyToPath{{
					Key:  opentelemetry.CollectorConfigMapEntry,
					Path: opentelemetry.CollectorConfigMapEntry,
				}},
			},
		},
	}}

	if len(instance.Spec.VolumeMounts) > 0 {
		volumeMounts = append(volumeMounts, instance.Spec.VolumeMounts...)
	}

	if len(instance.Spec.Volumes) > 0 {
		volumes = append(volumes, instance.Spec.Volumes...)
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   instance.Namespace,
			Labels:      labels,
			Annotations: instance.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: specAnnotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: ServiceAccountNameFor(ctx),
					Containers: []corev1.Container{{
						Name:         "opentelemetry-collector",
						Image:        image,
						VolumeMounts: volumeMounts,
						Args:         args,
					}},
					Volumes: volumes,
				},
			},
		},
	}
}

func (r *ReconcileOpenTelemetryCollector) reconcileExpectedDeployments(ctx context.Context, expected []*appsv1.Deployment) error {
	logger := ctx.Value(opentelemetry.ContextLogger).(logr.Logger)
	for _, obj := range expected {
		desired := obj

		// #nosec G104 (CWE-703): Errors unhandled.
		r.setControllerReference(ctx, desired)

		deps := r.clientset.Kubernetes.AppsV1().Deployments(desired.Namespace)

		existing, err := deps.Get(desired.Name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			if existing, err = deps.Create(desired); err != nil {
				return fmt.Errorf("failed to create: %v", err)
			}

			logger.WithValues("deployment.name", desired.Name, "deployment.namespace", desired.Namespace).V(2).Info("created")
			continue
		} else if err != nil {
			return fmt.Errorf("failed to get: %v", err)
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		if updated.Annotations == nil {
			updated.Annotations = map[string]string{}
		}
		if updated.Labels == nil {
			updated.Labels = map[string]string{}
		}

		updated.Spec = desired.Spec
		updated.ObjectMeta.OwnerReferences = desired.ObjectMeta.OwnerReferences

		for k, v := range desired.ObjectMeta.Annotations {
			updated.ObjectMeta.Annotations[k] = v
		}
		for k, v := range desired.ObjectMeta.Labels {
			updated.ObjectMeta.Labels[k] = v
		}

		if updated, err = deps.Update(updated); err != nil {
			return fmt.Errorf("failed to apply changes: %v", err)
		}
		logger.V(2).Info("applied", "deployment.name", desired.Name, "deployment.namespace", desired.Namespace)
	}

	return nil
}

func (r *ReconcileOpenTelemetryCollector) deleteDeployments(ctx context.Context, expected []*appsv1.Deployment) error {
	instance := ctx.Value(opentelemetry.ContextInstance).(*v1alpha1.OpenTelemetryCollector)
	logger := ctx.Value(opentelemetry.ContextLogger).(logr.Logger)
	deps := r.clientset.Kubernetes.AppsV1().Deployments(instance.Namespace)

	opts := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			"app.kubernetes.io/instance":   fmt.Sprintf("%s.%s", instance.Namespace, instance.Name),
			"app.kubernetes.io/managed-by": "opentelemetry-operator",
		}).String(),
	}
	list, err := deps.List(opts)
	if err != nil {
		return fmt.Errorf("failed to list: %v", err)
	}

	for _, existing := range list.Items {
		del := true
		for _, keep := range expected {
			if keep.Name == existing.Name && keep.Namespace == existing.Namespace {
				del = false
			}
		}

		if del {
			if err := deps.Delete(existing.Name, &metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("failed to delete: %v", err)
			}
			logger.V(2).Info("deleted", "deployment.name", existing.Name, "deployment.namespace", existing.Namespace)
		}
	}

	return nil
}
