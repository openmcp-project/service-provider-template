//go:generate opencontrolplane-gen
package webhook

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	// opencontrolplane-gen:replace github.com/openmcp-project/service-provider-template=MODULE
	"github.com/openmcp-project/service-provider-template/api/v1alpha1"
)

// opencontrolplane-gen:replace foo=KIND_LOWER
const webhookName = "foo-webhook"

// opencontrolplane-gen:replace Foo=KIND
type FooWebhook struct {
	client.Client
}

// opencontrolplane-gen:replace Foo=KIND
func SetupFooWebhookWithManager(mgr ctrl.Manager) error {
	// opencontrolplane-gen:replace Foo=KIND
	wh := &FooWebhook{
		Client: mgr.GetClient(),
	}

	// opencontrolplane-gen:replace Foo=KIND
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.Foo{}).
		WithDefaulter(wh).
		WithValidator(wh).
		Complete()
}

// opencontrolplane-gen:replace Foo=KIND
var _ admission.Defaulter[*v1alpha1.Foo] = &FooWebhook{}

// Default implements admission.Defaulter so a webhook will be registered for the type
// opencontrolplane-gen:replace Foo=KIND
func (wh *FooWebhook) Default(ctx context.Context, obj *v1alpha1.Foo) error {
	l := logf.FromContext(ctx).WithName(webhookName)
	// opencontrolplane-gen:replace Foo=KIND
	l.Info("Default Foo ...", "name", obj.GetName())
	return nil
}

// opencontrolplane-gen:replace Foo=KIND
var _ admission.Validator[*v1alpha1.Foo] = &FooWebhook{}

// ValidateCreate implements admission.Validator[] so a webhook will be registered for the type
// opencontrolplane-gen:replace Foo=KIND
func (wh *FooWebhook) ValidateCreate(ctx context.Context, obj *v1alpha1.Foo) (admission.Warnings, error) {
	l := logf.FromContext(ctx).WithName(webhookName)
	// opencontrolplane-gen:replace Foo=KIND
	l.Info("Validate Foo create...", "name", obj.GetName())
	return admission.Warnings{}, nil
}

// ValidateUpdate implements admission.Validator[] so a webhook will be registered for the type
// opencontrolplane-gen:replace Foo=KIND
func (wh *FooWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj *v1alpha1.Foo) (admission.Warnings, error) {
	l := logf.FromContext(ctx).WithName(webhookName)
	// opencontrolplane-gen:replace Foo=KIND
	l.Info("Validate Foo update...", "oldGeneration", oldObj.GetGeneration(), "newGeneration", newObj.GetGeneration())
	return admission.Warnings{}, nil
}

// ValidateDelete implements admission.Validator[] so a webhook will be registered for the type
// opencontrolplane-gen:replace Foo=KIND
func (wh *FooWebhook) ValidateDelete(ctx context.Context, obj *v1alpha1.Foo) (admission.Warnings, error) {
	l := logf.FromContext(ctx).WithName(webhookName)
	// opencontrolplane-gen:replace Foo=KIND
	l.Info("Validate Foo delete...", "name", obj.GetName())
	return admission.Warnings{}, nil
}
