//go:generate opencontrolplane-gen
package e2e

import (
	"context"
	"testing"
	"time"

	// opencontrolplane-gen:if SAMPLECODE=true
	corev1 "k8s.io/api/core/v1"
	// opencontrolplane-gen:fi
	// opencontrolplane-gen:if SAMPLECODE=true
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	// opencontrolplane-gen:fi
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	// opencontrolplane-gen:if SAMPLECODE=true
	"github.com/openmcp-project/openmcp-testing/pkg/clusterutils"
	// opencontrolplane-gen:fi
	// opencontrolplane-gen:if SAMPLECODE=true
	"github.com/openmcp-project/openmcp-testing/pkg/conditions"
	// opencontrolplane-gen:fi
	"github.com/openmcp-project/openmcp-testing/pkg/providers"
	"github.com/openmcp-project/openmcp-testing/pkg/resources"
)

func TestServiceProvider(t *testing.T) {
	// opencontrolplane-gen:if SAMPLECODE=true
	var onboardingList unstructured.UnstructuredList
	// opencontrolplane-gen:fi
	basicProviderTest := features.New("provider test").
		Setup(func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			if _, err := resources.CreateObjectsFromDir(ctx, c, "platform"); err != nil {
				t.Errorf("failed to create platform cluster objects: %v", err)
			}
			return ctx
		}).
		Setup(providers.CreateMCP("test-controlplane")).
		// opencontrolplane-gen:if SAMPLECODE=true
		Assess("verify service can be successfully consumed",
			func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
				onboardingConfig, err := clusterutils.OnboardingConfig()
				if err != nil {
					t.Error(err)
					return ctx
				}
				objList, err := resources.CreateObjectsFromDir(ctx, onboardingConfig, "onboarding")
				if err != nil {
					t.Errorf("failed to create onboarding cluster objects: %v", err)
					return ctx
				}
				for _, obj := range objList.Items {
					if err := wait.For(conditions.Match(&obj, onboardingConfig, "Ready", corev1.ConditionTrue)); err != nil {
						t.Error(err)
					}
				}
				objList.DeepCopyInto(&onboardingList)
				return ctx
			},
		).
		Assess("verify domain objects can be created", providers.ImportDomainAPIs("test-controlplane", "controlplane")).
		Teardown(func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			onboardingConfig, err := clusterutils.OnboardingConfig()
			if err != nil {
				t.Error(err)
				return ctx
			}
			for _, obj := range onboardingList.Items {
				if err := resources.DeleteObject(ctx, onboardingConfig, &obj, wait.WithTimeout(time.Minute)); err != nil {
					t.Errorf("failed to delete onboarding object: %v", err)
				}
			}
			return ctx
		}).
		// opencontrolplane-gen:fi
		// opencontrolplane-gen:if SAMPLECODE=false
		// TODO add assess steps
		// opencontrolplane-gen:fi
		Teardown(providers.DeleteMCP("test-controlplane", wait.WithTimeout(5*time.Minute)))
	testenv.Test(t, basicProviderTest.Feature())
}
