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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OverlayNetworkSpec defines the desired state of OverlayNetwork
type OverlayNetworkSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of OverlayNetwork. Edit overlaynetwork_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// OverlayNetworkStatus defines the observed state of OverlayNetwork
type OverlayNetworkStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OverlayNetwork is the Schema for the overlaynetworks API
type OverlayNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OverlayNetworkSpec   `json:"spec,omitempty"`
	Status OverlayNetworkStatus `json:"status,omitempty"`
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
