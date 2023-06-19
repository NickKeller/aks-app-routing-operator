package dns

import (
	"strings"
	"testing"

	"github.com/Azure/aks-app-routing-operator/pkg/config"
	"github.com/Azure/aks-app-routing-operator/pkg/manifests"
	"github.com/stretchr/testify/require"
)

var (
	privateZoneOne = "/subscriptions/test-private-subscription/resourceGroups/test-rg-private/providers/Microsoft.Network/privatednszones/test-one.com"
	privateZoneTwo = "/subscriptions/test-private-subscription/resourceGroups/test-rg-private/providers/Microsoft.Network/privatednszones/test-two.com"
	privateZones   = []string{privateZoneOne, privateZoneTwo}

	publicZoneOne = "/subscriptions/test-public-subscription/resourceGroups/test-rg-private/providers/Microsoft.Network/dnszones/test-one.com"
	publicZoneTwo = "/subscriptions/test-public-subscription/resourceGroups/test-rg-private/providers/Microsoft.Network/dnszones/test-two.com"
	publicZones   = []string{publicZoneOne, publicZoneTwo}

	zones = strings.Join(append(privateZones, publicZones...), ",")

	publicConfig = &config.Config{
		NS:              "test-ns",
		DisableKeyvault: false,
		PrivateZoneConfig: config.DnsZoneConfig{
			ZoneIds:       nil,
			Subscription:  "",
			ResourceGroup: "",
		},
		PublicZoneConfig: config.DnsZoneConfig{
			ZoneIds:       publicZones,
			Subscription:  "test-public-subscription",
			ResourceGroup: "test-public-rg",
		},
	}
	privateConfig = &config.Config{
		NS:              "test-ns",
		DisableKeyvault: false,
		PrivateZoneConfig: config.DnsZoneConfig{
			ZoneIds:       privateZones,
			Subscription:  "test-private-subscription",
			ResourceGroup: "test-private-rg",
		},
		PublicZoneConfig: config.DnsZoneConfig{
			ZoneIds:       nil,
			Subscription:  "",
			ResourceGroup: "",
		},
	}
	fullConfig = &config.Config{
		NS:              "test-ns",
		DisableKeyvault: false,
		PrivateZoneConfig: config.DnsZoneConfig{
			ZoneIds:       privateZones,
			Subscription:  "test-private-subscription",
			ResourceGroup: "test-private-rg",
		},
		PublicZoneConfig: config.DnsZoneConfig{
			ZoneIds:       publicZones,
			Subscription:  "test-public-subscription",
			ResourceGroup: "test-public-rg",
		},
	}
	zoneless = &config.Config{
		NS:              "test-ns",
		DisableKeyvault: false,
	}
)

// TODO: add tests for cleanup names
func TestGenerateZoneConfigs_PublicOnly(t *testing.T) {
	zoneConfigs, _ := generateZoneConfigs(publicConfig)

	require.Equal(t, 1, len(zoneConfigs))
	require.Equal(t, publicConfig.PublicZoneConfig.ZoneIds, zoneConfigs[0].DnsZoneResourceIDs)
	require.Equal(t, manifests.Provider(manifests.PublicProvider), zoneConfigs[0].Provider)
	require.Equal(t, publicConfig.PublicZoneConfig.Subscription, zoneConfigs[0].Subscription)
}

func TestGenerateZoneConfigs_PrivateOnly(t *testing.T) {
	zoneConfigs, _ := generateZoneConfigs(privateConfig)

	require.Equal(t, 1, len(zoneConfigs))
	require.Equal(t, privateConfig.PrivateZoneConfig.ZoneIds, zoneConfigs[0].DnsZoneResourceIDs)
	require.Equal(t, manifests.Provider(manifests.PrivateProvider), zoneConfigs[0].Provider)
	require.Equal(t, privateConfig.PrivateZoneConfig.Subscription, zoneConfigs[0].Subscription)
}

func TestGenerateZoneConfigs_All(t *testing.T) {
	zoneConfigs, _ := generateZoneConfigs(fullConfig)

	require.Equal(t, len(zoneConfigs), 2)

	prConfig := zoneConfigs[0]
	pbConfig := zoneConfigs[1]

	require.Equal(t, fullConfig.PrivateZoneConfig.ZoneIds, prConfig.DnsZoneResourceIDs)
	require.Equal(t, fullConfig.PublicZoneConfig.ZoneIds, pbConfig.DnsZoneResourceIDs)

	require.Equal(t, manifests.Provider(manifests.PrivateProvider), prConfig.Provider)
	require.Equal(t, manifests.Provider(manifests.PublicProvider), pbConfig.Provider)

	require.Equal(t, fullConfig.PrivateZoneConfig.Subscription, prConfig.Subscription)
	require.Equal(t, fullConfig.PublicZoneConfig.Subscription, pbConfig.Subscription)

}

func TestGenerateZoneConfigs_zoneless(t *testing.T) {
	zoneConfigs, _ := generateZoneConfigs(zoneless)
	require.Equal(t, len(zoneConfigs), 0)
}