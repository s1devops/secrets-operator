/*
Copyright 2021.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretMappingItemType string

const (
	TYPE_STRING   SecretMappingItemType = "string"
	TYPE_PASS     SecretMappingItemType = "pass"
	TYPE_TEMPLATE SecretMappingItemType = "template"
)

type SecretSourceRef struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type SecretMappingItem struct {
	Name  string                `json:"name"`
	Value string                `json:"value"`
	Type  SecretMappingItemType `json:"type"`
}

// SecretMappingSpec defines the desired state of SecretMapping
type SecretMappingSpec struct {
	// Name of the secret, defaults to the name of the SecretMapping
	Name     string              `json:"name,omitempty"`
	Source   SecretSourceRef     `json:"source"`
	Mappings []SecretMappingItem `json:"mappings"`
}

// SecretMappingStatus defines the observed state of SecretMapping
type SecretMappingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Revision string `json:"revision"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SecretMapping is the Schema for the secretmappings API
type SecretMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretMappingSpec   `json:"spec,omitempty"`
	Status SecretMappingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SecretMappingList contains a list of SecretMapping
type SecretMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretMapping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretMapping{}, &SecretMappingList{})
}
