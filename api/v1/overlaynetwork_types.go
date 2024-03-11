/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OverlayNetworkSpec defines the desired state of OverlayNetwork
type OverlayNetworkSpec struct {
	// Name is the overlay network identifier.
	// The actual overlay network is prefixed with the current namespace
	// to assure uniqueness across namespaces.
	Name string `json:"name"`

	// The maximum transmission unit (MTU).
	// The default value, 1300, is automatically set by the kernel.
	// +optional
	Mtu string `json:"mtu,omitempty"`

	// The subnet to use for the network across the cluster.
	// Only include the CIDR for the node. E.g. 10.100.200.0/24.
	//
	// IPv6 (2001:DBB::/64) and dual-stack (192.168.100.0/24,2001:DBB::/64) subnets are supported.
	//
	// When omitted, the logical switch implementing the network only provides layer 2 communication,
	// and users must configure IP addresses for the pods.
	// Port security only prevents MAC spoofing.
	// +optional
	Subnets string `json:"subnets,omitempty"`

	// A comma-separated list of CIDRs and IP addresses.
	// IP addresses are removed from the assignable IP address pool and are never passed to the pods.
	// +optional
	ExcludeSubnets string `json:"excludeSubnets,omitempty"`
}

//+kubebuilder:object:root=true

// OverlayNetwork is the Schema for the overlaynetworks API
type OverlayNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OverlayNetworkSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// OverlayNetworkList contains a list of OverlayNetwork
type OverlayNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OverlayNetwork `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OverlayNetwork{}, &OverlayNetworkList{})
}
