//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MeteringReceiverSpec defines the desired state of MeteringReceiver
type MeteringReceiverSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Version         string                      `json:"version"`
	ImageRegistry   string                      `json:"imageRegistry,omitempty"`
	ImageTagPostfix string                      `json:"imageTagPostfix,omitempty"`
	ClusterIssuer   string                      `json:"clusterIssuer,omitempty"`
	MongoDB         MeteringReceiverSpecMongoDB `json:"mongodb"`
}

// MeteringStatus defines the observed state of each Metering service
type MeteringStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// PodNames are the names of the metering pods
	PodNames []string `json:"podNames"`
}

// MeteringSpecMongoDB defines the MongoDB configuration in all the Metering specs
type MeteringReceiverSpecMongoDB struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	UsernameSecret     string `json:"usernameSecret"`
	UsernameKey        string `json:"usernameKey"`
	PasswordSecret     string `json:"passwordSecret"`
	PasswordKey        string `json:"passwordKey"`
	ClusterCertsSecret string `json:"clustercertssecret"`
	ClientCertsSecret  string `json:"clientcertssecret"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MeteringReceiver is the Schema for the meteringreceivers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=meteringreceivers,scope=Namespaced
type MeteringReceiver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MeteringReceiverSpec `json:"spec,omitempty"`
	Status MeteringStatus       `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MeteringReceiverList contains a list of MeteringReceiver
type MeteringReceiverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MeteringReceiver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MeteringReceiver{}, &MeteringReceiverList{})
}
