package dns

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/openmcp-project/openmcp-testing/pkg/clusterutils"
	"github.com/openmcp-project/openmcp-testing/pkg/resources"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	kindcluster "sigs.k8s.io/kind/pkg/cluster"
)

type DNSConfig struct {
	// Namespace defines the namespace where the core components like the openmcp-operator is deployed
	// Defaults to openmcp-system if no value is provided.
	Namespace string
	// ClusterPurpose defines the purpose that will be used to create the cluster for the dedicated dns deployment (if DedicatedDNS = true).
	ClusterPurpose string
	// DedicatedDNS allows to choose between deploying a dedicated dns and platform service dns
	// or a simple host alias config.
	DedicatedDNS bool
	// TargetContainer defines which kube-apiserver will be patched as part of the the initial setup.
	// Defaults to the onboarding cluster container if no value is provided.
	TargetContainer string
	// TLSRouteKey defines where to get the hostname from
	TLSRouteKey types.NamespacedName
}

const defaultNamespace = "openmcp-system"

func SetupDNS(config DNSConfig) features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		gatewayv1.Install(c.Client().Resources().GetScheme())
		gatewayv1alpha2.Install(c.Client().Resources().GetScheme())
		openmcpSystemNamespace := config.Namespace
		if openmcpSystemNamespace == "" {
			openmcpSystemNamespace = defaultNamespace
		}
		targetContainer := config.TargetContainer
		if targetContainer == "" {
			onboardingClusterContainer, err := onboardingClusterContainer()
			if err != nil {
				t.Fatalf("failed to determine onboarding container name: %v", err)
				return ctx
			}
			targetContainer = onboardingClusterContainer
		}
		if !config.DedicatedDNS {
			// inject host into kube-apiserver
			gwIP := getGatewayIP(ctx, t, c, "default", openmcpSystemNamespace)
			// opencontrolplane-gen:replace foo=KIND_LOWER
			wbHostname := getHostname(ctx, t, c, config.TLSRouteKey.Name, config.TLSRouteKey.Namespace)
			klog.Infof("add host %s with ip %s to /etc/hosts of the (%s) kube-apiserver", wbHostname, gwIP, targetContainer)
			if err := addHostToKubeAPIServer(targetContainer, wbHostname, gwIP); err != nil {
				t.Errorf("failed to add host to kube-apiserver: %v", err)
				return ctx
			}
			if err := waitForKubeAPIServerRestart(targetContainer, 3*time.Minute); err != nil {
				t.Errorf("kube-apiserver didn't restart properly: %v", err)
			}
			return ctx
		}
		// create dedicate dns deployment to use with platform service dns
		if err := createCluster(ctx, c, ClusterRequest{
			Name:      "dns",
			Namespace: openmcpSystemNamespace,
			Purpose:   config.ClusterPurpose,
		}); err != nil {
			t.Errorf("failed to create dns cluster: %v", err)
		}
		dnsClusterConfig, err := clusterutils.ConfigByPrefix("dns", "default")
		if err != nil {
			t.Errorf("failed to retrieve dns cluster config: %v", err)
			return ctx
		}
		if _, err := resources.CreateObjectsFromDir(ctx, dnsClusterConfig, "dns/etcd"); err != nil {
			t.Errorf("failed to deploy etcd for dns: %v", err)
			return ctx
		}
		if _, err := resources.CreateObjectsFromDir(ctx, c, "dns/coredns"); err != nil {
			t.Errorf("failed to deploy coredns: %v", err)
			return ctx
		}
		// install ps-dns with the etcd IP for the external-dns config
		etcdIP := getLoadBalancerIP(ctx, t, dnsClusterConfig, "etcd-external", "default")
		psDNSConfig := PlatformServiceDNSConfig{
			Version:                 "v0.1.0",
			EtcdIP:                  etcdIP,
			ExternalDNSChartVersion: "v0.21.0",
		}
		if err := createPlatformServiceDNS(ctx, t, c, psDNSConfig); err != nil {
			t.Errorf("failed to create platform service dns config: %v", err)
		}
		// inject additional nameserver into kube-apiserver
		nameserverIP := getLoadBalancerIP(ctx, t, dnsClusterConfig, "coredns", "default")
		klog.Infof("add nameserver with ip %s to dns config of (%s) kube-apiserver", nameserverIP, targetContainer)
		if err := addNameserverToKubeAPIServer(targetContainer, nameserverIP); err != nil {
			t.Errorf("failed to add host to kube-apiserver: %v", err)
			return ctx
		}
		if err := waitForKubeAPIServerRestart(targetContainer, 3*time.Minute); err != nil {
			t.Errorf("kube-apiserver didn't restart properly: %v", err)
		}
		return ctx
	}
}

func onboardingClusterContainer() (string, error) {
	kind := kindcluster.NewProvider()
	clusters, err := kind.List()
	if err != nil {
		return "", err
	}
	for _, clusterName := range clusters {
		if strings.HasPrefix(clusterName, "onboarding") {
			nodes, err := kind.ListNodes(clusterName)
			if err != nil {
				return "", fmt.Errorf("failed to retrieve onboarding cluster nodes: %w", err)
			}
			return nodes[0].String(), nil
		}
	}
	return "", errors.New("onboarding cluster not found")
}
