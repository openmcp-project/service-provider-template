package app

import (
	"context"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openmcp-project/controller-utils/pkg/clusters"
	"github.com/openmcp-project/openmcp-operator/lib/clusteraccess"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	localaccess "github.com/openmcp-project/opencontrolplane-runtime/pkg/serviceprovider/clusteraccess"
	clustersv1alpha1 "github.com/openmcp-project/openmcp-operator/api/clusters/v1alpha1"
)

const (
	debugEnvVar         = "DEV_DEBUG"
	appInit     appMode = "-init"
	appRun      appMode = "-run"
)

type appMode string

func requestOnboardingClusterAccess(ctx context.Context, mgr clusteraccess.Manager, platformCluster *clusters.Cluster, onboardingScheme *runtime.Scheme, permissions []clustersv1alpha1.PermissionsRequest, providerName string, mode appMode) (*clusters.Cluster, error) {
	cluster, err := mgr.CreateAndWaitForCluster(ctx, clustersv1alpha1.PURPOSE_ONBOARDING+string(mode), clustersv1alpha1.PURPOSE_ONBOARDING, onboardingScheme, permissions)
	if err != nil {
		return cluster, err
	}
	if envFlagEnabled(debugEnvVar) {
		return patchOnboardingClient(ctx, platformCluster, cluster, providerName)
	}
	return cluster, nil
}

func patchOnboardingClient(ctx context.Context, platformCluster *clusters.Cluster, onboardingCluster *clusters.Cluster, providerName string) (*clusters.Cluster, error) {
	onboardingAr := &clustersv1alpha1.AccessRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusteraccess.StableRequestNameFromLocalName(providerName, "onboarding-run"),
			Namespace: os.Getenv("POD_NAMESPACE"),
		},
	}
	if err := platformCluster.Client().Get(ctx, client.ObjectKeyFromObject(onboardingAr), onboardingAr); err != nil {
		return onboardingCluster, err
	}
	return localaccess.MustPatchClusterClient(ctx, onboardingAr, onboardingCluster), nil
}
