// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"

	v1alpha1 "github.com/open-telemetry/opentelemetry-operator/pkg/apis/opentelemetry/v1alpha1"
)

// FakeOpenTelemetryCollectors implements OpenTelemetryCollectorInterface
type FakeOpenTelemetryCollectors struct {
	Fake *FakeOpentelemetryV1alpha1
	ns   string
}

var opentelemetrycollectorsResource = schema.GroupVersionResource{Group: "opentelemetry.io", Version: "v1alpha1", Resource: "opentelemetrycollectors"}

var opentelemetrycollectorsKind = schema.GroupVersionKind{Group: "opentelemetry.io", Version: "v1alpha1", Kind: "OpenTelemetryCollector"}

// Get takes name of the openTelemetryCollector, and returns the corresponding openTelemetryCollector object, and an error if there is any.
func (c *FakeOpenTelemetryCollectors) Get(name string, options v1.GetOptions) (result *v1alpha1.OpenTelemetryCollector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(opentelemetrycollectorsResource, c.ns, name), &v1alpha1.OpenTelemetryCollector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenTelemetryCollector), err
}

// List takes label and field selectors, and returns the list of OpenTelemetryCollectors that match those selectors.
func (c *FakeOpenTelemetryCollectors) List(opts v1.ListOptions) (result *v1alpha1.OpenTelemetryCollectorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(opentelemetrycollectorsResource, opentelemetrycollectorsKind, c.ns, opts), &v1alpha1.OpenTelemetryCollectorList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.OpenTelemetryCollectorList{ListMeta: obj.(*v1alpha1.OpenTelemetryCollectorList).ListMeta}
	for _, item := range obj.(*v1alpha1.OpenTelemetryCollectorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested openTelemetryCollectors.
func (c *FakeOpenTelemetryCollectors) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(opentelemetrycollectorsResource, c.ns, opts))

}

// Create takes the representation of a openTelemetryCollector and creates it.  Returns the server's representation of the openTelemetryCollector, and an error, if there is any.
func (c *FakeOpenTelemetryCollectors) Create(openTelemetryCollector *v1alpha1.OpenTelemetryCollector) (result *v1alpha1.OpenTelemetryCollector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(opentelemetrycollectorsResource, c.ns, openTelemetryCollector), &v1alpha1.OpenTelemetryCollector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenTelemetryCollector), err
}

// Update takes the representation of a openTelemetryCollector and updates it. Returns the server's representation of the openTelemetryCollector, and an error, if there is any.
func (c *FakeOpenTelemetryCollectors) Update(openTelemetryCollector *v1alpha1.OpenTelemetryCollector) (result *v1alpha1.OpenTelemetryCollector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(opentelemetrycollectorsResource, c.ns, openTelemetryCollector), &v1alpha1.OpenTelemetryCollector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenTelemetryCollector), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeOpenTelemetryCollectors) UpdateStatus(openTelemetryCollector *v1alpha1.OpenTelemetryCollector) (*v1alpha1.OpenTelemetryCollector, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(opentelemetrycollectorsResource, "status", c.ns, openTelemetryCollector), &v1alpha1.OpenTelemetryCollector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenTelemetryCollector), err
}

// Delete takes name of the openTelemetryCollector and deletes it. Returns an error if one occurs.
func (c *FakeOpenTelemetryCollectors) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(opentelemetrycollectorsResource, c.ns, name), &v1alpha1.OpenTelemetryCollector{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOpenTelemetryCollectors) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(opentelemetrycollectorsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.OpenTelemetryCollectorList{})
	return err
}

// Patch applies the patch and returns the patched openTelemetryCollector.
func (c *FakeOpenTelemetryCollectors) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.OpenTelemetryCollector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(opentelemetrycollectorsResource, c.ns, name, pt, data, subresources...), &v1alpha1.OpenTelemetryCollector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OpenTelemetryCollector), err
}
