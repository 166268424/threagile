/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type TrustBoundaryType int

const (
	NetworkOnPrem TrustBoundaryType = iota
	NetworkDedicatedHoster
	NetworkVirtualLAN
	NetworkCloudProvider
	NetworkCloudSecurityGroup
	NetworkPolicyNamespaceIsolation
	ExecutionEnvironment
)

func TrustBoundaryTypeValues() []TypeEnum {
	return []TypeEnum{
		NetworkOnPrem,
		NetworkDedicatedHoster,
		NetworkVirtualLAN,
		NetworkCloudProvider,
		NetworkCloudSecurityGroup,
		NetworkPolicyNamespaceIsolation,
		ExecutionEnvironment,
	}
}

var TrustBoundaryTypeDescription = [...]TypeDescription{
	{"network-on-prem", "The whole network is on prem"},
	{"network-dedicated-hoster", "The network is at a dedicated hoster"},
	{"network-virtual-lan", "Network is a VLAN"},
	{"network-cloud-provider", "Network is at a cloud provider"},
	{"network-cloud-security-group", "Cloud rules controlling network traffic"},
	{"network-policy-namespace-isolation", "Segregation in a Kubernetes cluster"},
	{"execution-environment", "Logical group of items (not a protective network boundary in that sense). More like a namespace or another logical group of items"},
}

func ParseTrustBoundary(value string) (trustBoundary TrustBoundaryType, err error) {
	value = strings.TrimSpace(value)
	for _, candidate := range TrustBoundaryTypeValues() {
		if candidate.String() == value {
			return candidate.(TrustBoundaryType), err
		}
	}
	return trustBoundary, errors.New("Unable to parse into type: " + value)
}

func (what TrustBoundaryType) String() string {
	// NOTE: maintain list also in schema.json for validation in IDEs
	return TrustBoundaryTypeDescription[what].Name
}

func (what TrustBoundaryType) Explain() string {
	return TrustBoundaryTypeDescription[what].Description
}

func (what TrustBoundaryType) IsNetworkBoundary() bool {
	return what == NetworkOnPrem || what == NetworkDedicatedHoster || what == NetworkVirtualLAN ||
		what == NetworkCloudProvider || what == NetworkCloudSecurityGroup || what == NetworkPolicyNamespaceIsolation
}

func (what TrustBoundaryType) IsWithinCloud() bool {
	return what == NetworkCloudProvider || what == NetworkCloudSecurityGroup
}

func (what TrustBoundaryType) MarshalJSON() ([]byte, error) {
	return json.Marshal(what.String())
}

func (what *TrustBoundaryType) UnmarshalJSON([]byte) error {
	for index, description := range TrustBoundaryTypeDescription {
		if strings.ToLower(what.String()) == strings.ToLower(description.Name) {
			*what = TrustBoundaryType(index)
			return nil
		}
	}

	return fmt.Errorf("unknown trust boundary type value %q", int(*what))
}
