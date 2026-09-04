//go:generate opencontrolplane-gen
package e2e

import (
	"context"
	"testing"
	"time"

	// opencontrolplane-gen:if SAMPLECODE=true
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	// opencontrolplane-gen:fi

	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/openmcp-project/openmcp-testing/pkg/clusterutils"
	"github.com/openmcp-project/openmcp-testing/pkg/clusterutils/apiserver"
	"github.com/openmcp-project/openmcp-testing/pkg/providers"

	// opencontrolplane-gen:replace github.com/openmcp-project/service-provider-template=MODULE
	apiv1alpha1 "github.com/openmcp-project/service-provider-template/api/v1alpha1"
	// opencontrolplane-gen:replace github.com/openmcp-project/service-provider-template=MODULE

	// opencontrolplane-gen:if SAMPLECODE=true
	openmcpconditions "github.com/openmcp-project/openmcp-testing/pkg/conditions"
	// opencontrolplane-gen:fi
)

func TestServiceProvider(t *testing.T) {
	basicProviderTest := features.New("provider test").
		Setup(func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			apiv1alpha1.AddToScheme(c.Client().Resources().GetScheme())
			config := &apiv1alpha1.ProviderConfig{}
			// opencontrolplane-gen:replace configname=PROVIDER_NAME
			config.SetName("configname")
			if err := c.Client().Resources().Create(ctx, config); err != nil {
				t.Errorf("failed to create ProviderConfig object: %v", err)
			}
			return ctx
		}).
		Setup(providers.CreateMCP("test-controlplane")).
		// opencontrolplane-gen:if SAMPLECODE=true
		Setup(prepareWebhookExecution()).
		Assess("verify provider can be successfully consumed",
			func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
				config := c
				config, err := clusterutils.OnboardingConfig()
				if err != nil {
					t.Error(err)
					return ctx
				}
				apiv1alpha1.AddToScheme(config.Client().Resources().GetScheme())
				// opencontrolplane-gen:replace Foo=KIND
				api := &apiv1alpha1.Foo{}
				api.SetName("test-controlplane")
				api.SetNamespace("default")
				if err := config.Client().Resources().Create(ctx, api); err != nil {
					// opencontrolplane-gen:replace Foo=KIND
					t.Errorf("failed to create Foo object: %v", err)
				}
				if err := wait.For(openmcpconditions.Match(api, config, "Ready", corev1.ConditionTrue)); err != nil {
					t.Error(err)
				}
				return ctx
			},
		).
		Assess("verify domain objects can be created",
			func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
				mcpConfig, err := clusterutils.MCPConfig(ctx, c, "test-controlplane")
				if err != nil {
					t.Error(err)
					return ctx
				}
				domainObj := &unstructured.Unstructured{}
				domainObj.SetName("test-domain-object")
				domainObj.SetNamespace("default")
				domainObj.SetAPIVersion("example.domain/v1alpha1")
				domainObj.SetKind("Foo")
				if err := mcpConfig.Client().Resources().Create(ctx, domainObj); err != nil {
					t.Errorf("failed to create domain object on controlplane: %v", err)
				}
				return ctx
			},
		).
		Assess("verify service deletion is blocked due to existing domain service object",
			func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
				config := c
				config, err := clusterutils.OnboardingConfig()
				if err != nil {
					t.Error(err)
					return ctx
				}
				apiv1alpha1.AddToScheme(config.Client().Resources().GetScheme())
				// opencontrolplane-gen:replace Foo=KIND
				api := &apiv1alpha1.Foo{}
				api.SetName("test-controlplane")
				api.SetNamespace("default")
				if err := config.Client().Resources().Delete(ctx, api); err != nil {
					// opencontrolplane-gen:replace Foo=KIND
					t.Errorf("failed to delete Foo object: %v", err)
				}
				// verify object is stuck in Terminating with UserResourcesPresent reason
				if err := wait.For(func(ctx context.Context) (bool, error) {
					if err := config.Client().Resources().Get(ctx, api.GetName(), api.GetNamespace(), api); err != nil {
						return false, nil
					}
					if c := meta.FindStatusCondition(api.Status.Conditions, "DeletionBlocked"); c != nil && c.Status == metav1.ConditionTrue {
						return true, nil
					}
					return false, nil
				}); err != nil {
					t.Errorf("expected deletion to be blocked with reason UserResourcesPresent: %v", err)
				}
				return ctx
			},
		).
		Assess("delete domain object",
			func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
				mcpConfig, err := clusterutils.MCPConfig(ctx, c, "test-controlplane")
				if err != nil {
					t.Error(err)
					return ctx
				}
				domainObj := &unstructured.Unstructured{}
				domainObj.SetName("test-domain-object")
				domainObj.SetNamespace("default")
				domainObj.SetAPIVersion("example.domain/v1alpha1")
				domainObj.SetKind("Foo")
				if err := mcpConfig.Client().Resources().Delete(ctx, domainObj); err != nil {
					t.Errorf("failed to create domain object on controlplane: %v", err)
				}
				return ctx
			},
		).
		Assess("verify service is deleted",
			func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
				config := c
				config, err := clusterutils.OnboardingConfig()
				if err != nil {
					t.Error(err)
					return ctx
				}
				apiv1alpha1.AddToScheme(config.Client().Resources().GetScheme())
				// opencontrolplane-gen:replace Foo=KIND
				api := &apiv1alpha1.Foo{}
				api.SetName("test-controlplane")
				api.SetNamespace("default")
				if err := wait.For(conditions.New(config.Client().Resources()).ResourceDeleted(api)); err != nil {
					// opencontrolplane-gen:replace Foo=KIND
					t.Errorf("expected Foo to be deleted after domain object removal, but it still exists: %v", err)
				}
				return ctx
			},
		).
		// opencontrolplane-gen:fi
		// opencontrolplane-gen:if SAMPLECODE=false
		// TODO add assess steps
		// opencontrolplane-gen:fi
		Teardown(providers.DeleteMCP("test-controlplane", wait.WithTimeout(5*time.Minute)))
	testenv.Test(t, basicProviderTest.Feature())
}

// opencontrolplane-gen:if SAMPLECODE=true
// prepareWebhookExecution updates the onboarding cluster kube-apiserver to resolve the platform cluster gateway IP when calling the service provider webhook.
func prepareWebhookExecution() features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		runtime.Must(gatewayv1.Install(c.Client().Resources().GetScheme()))
		runtime.Must(gatewayv1alpha2.Install(c.Client().Resources().GetScheme()))
		gwIP := getGatewayIP(ctx, t, c, "default", "openmcp-system")
		// opencontrolplane-gen:replace foo=KIND_LOWER
		wbHostname := getHostname(ctx, t, c, "foo-webhook", "openmcp-system")
		updater, err := apiserver.NewUpdater()
		if err != nil {
			t.Fatalf("failed to create api-server updater: %v", err)
		}
		if err := updater.AddHostAlias(wbHostname, gwIP); err != nil {
			t.Fatalf("failed to add host to kube-apiserver: %v", err)
		}
		return ctx
	}
}

// getHostname retrieves the first hostname defined in the TLSRoute with the given name and namespace.
func getHostname(ctx context.Context, t *testing.T, config *envconf.Config, name, namespace string) string {
	t.Helper()
	tlsRoute := &gatewayv1alpha2.TLSRoute{}
	tlsRoute.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "gateway.networking.k8s.io",
		Version: "v1alpha2",
		Kind:    "TLSRoute",
	})
	tlsRoute.SetName(name)
	tlsRoute.SetNamespace(namespace)
	if err := config.Client().Resources().Get(ctx, name, namespace, tlsRoute); err != nil {
		t.Fatalf("failed to get TLSRoute '%s/%s': %v", namespace, name, err)
	}
	if len(tlsRoute.Spec.Hostnames) == 0 {
		t.Fatalf("TLSRoute '%s/%s' does not have any hostnames defined", namespace, name)
	}
	return string(tlsRoute.Spec.Hostnames[0])
}

// getGatewayIP retrieves the first IP address of the Gateway with the given name and namespace.
func getGatewayIP(ctx context.Context, t *testing.T, config *envconf.Config, name, namespace string) string {
	t.Helper()
	gateway := &gatewayv1.Gateway{}
	gateway.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "gateway.networking.k8s.io",
		Version: "v1",
		Kind:    "Gateway",
	})
	gateway.SetName(name)
	gateway.SetNamespace(namespace)
	if err := config.Client().Resources().Get(ctx, name, namespace, gateway); err != nil {
		t.Fatalf("failed to get Gateway '%s/%s': %v", namespace, name, err)
	}
	for _, addr := range gateway.Status.Addresses {
		if addr.Type != nil && *addr.Type == gatewayv1.IPAddressType {
			return addr.Value
		}
	}
	t.Fatalf("Gateway '%s/%s' does not have any IP addresses exposed", namespace, name)
	return ""
}

// opencontrolplane-gen:fi
