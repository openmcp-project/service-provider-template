//go:generate opencontrolplane-gen
package webhook

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/openmcp-project/controller-utils/pkg/logging"

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
func SetupFooWebhookWithManager(ctx context.Context, mgr ctrl.Manager) error {
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
func (p *FooWebhook) Default(ctx context.Context, obj *v1alpha1.Foo) error {
	log := logging.FromContextOrPanic(ctx).WithName(webhookName)
	// opencontrolplane-gen:replace Foo=KIND
	log.Info("Default Foo...")
	return nil
}

// opencontrolplane-gen:replace Foo=KIND
var _ admission.Validator[*v1alpha1.Foo] = &FooWebhook{}

// ValidateCreate implements admission.Validator[] so a webhook will be registered for the type
// opencontrolplane-gen:replace Foo=KIND
func (v *FooWebhook) ValidateCreate(ctx context.Context, obj *v1alpha1.Foo) (admission.Warnings, error) {
	log := logging.FromContextOrPanic(ctx).WithName(webhookName)
	// opencontrolplane-gen:replace Foo=KIND
	log.Info("Validate Foo create...")
	return admission.Warnings{}, nil
}

// ValidateUpdate implements admission.Validator[] so a webhook will be registered for the type
// opencontrolplane-gen:replace Foo=KIND
func (v *FooWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj *v1alpha1.Foo) (admission.Warnings, error) {
	log := logging.FromContextOrPanic(ctx).WithName(webhookName)
	// opencontrolplane-gen:replace Foo=KIND
	log.Info("Validate Foo update...")
	return admission.Warnings{}, nil
}

// ValidateDelete implements admission.Validator[] so a webhook will be registered for the type
// opencontrolplane-gen:replace Foo=KIND
func (v *FooWebhook) ValidateDelete(ctx context.Context, obj *v1alpha1.Foo) (admission.Warnings, error) {
	log := logging.FromContextOrPanic(ctx).WithName(webhookName)
	// opencontrolplane-gen:replace Foo=KIND
	log.Info("Validate Foo delete...")
	return admission.Warnings{}, nil
}
