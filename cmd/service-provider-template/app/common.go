//go:generate opencontrolplane-gen
package app

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openmcp-project/controller-utils/pkg/clusters"
	"github.com/openmcp-project/openmcp-operator/lib/clusteraccess"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// opencontrolplane-gen:replace github.com/openmcp-project/service-provider-template=MODULE
	"github.com/openmcp-project/service-provider-template/api/providerscheme"

	localaccess "github.com/openmcp-project/opencontrolplane-runtime/pkg/serviceprovider/clusteraccess"
	clustersv1alpha1 "github.com/openmcp-project/openmcp-operator/api/clusters/v1alpha1"
)

const (
	debugEnvVar         = "DEV_DEBUG"
	appInit     appMode = "-init"
	appRun      appMode = "-run"
)

type appMode string

func requestOnboardingAccess(ctx context.Context, mgr clusteraccess.Manager, permissions []clustersv1alpha1.PermissionsRequest, opts *SharedOptions, mode appMode) (*clusters.Cluster, error) {
	onboardingScheme := providerscheme.OnboardingScheme(runtime.NewScheme())
	cluster, err := mgr.CreateAndWaitForCluster(ctx, clustersv1alpha1.PURPOSE_ONBOARDING+string(mode), clustersv1alpha1.PURPOSE_ONBOARDING, onboardingScheme, permissions)
	if err != nil {
		return cluster, err
	}
	if envFlagEnabled(debugEnvVar) {
		return patchOnboardingClient(ctx, cluster, opts, mode)
	}
	return cluster, nil
}

func patchOnboardingClient(ctx context.Context, onboardingCluster *clusters.Cluster, opts *SharedOptions, mode appMode) (*clusters.Cluster, error) {
	onboardingAr := &clustersv1alpha1.AccessRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusteraccess.StableRequestNameFromLocalName(opts.ProviderName, "onboarding"+string(mode)),
			Namespace: opts.ProviderNamespace,
		},
	}
	if err := opts.PlatformCluster.Client().Get(ctx, client.ObjectKeyFromObject(onboardingAr), onboardingAr); err != nil {
		return onboardingCluster, err
	}
	return localaccess.MustPatchClusterClient(ctx, onboardingAr, onboardingCluster), nil
}
