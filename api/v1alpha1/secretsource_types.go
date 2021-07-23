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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type SecretSourcePrivateKeyRef struct {
	// Name of the secret to use
	Name string `json:"name"`

	// The key of the base64 encoded private key
	Key string `json:"key"`
}

type SecretSourceGitRepositoryRef struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// SecretSourceSpec defines the desired state of SecretSource
type SecretSourceSpec struct {
	// Reference to the Flux GitRepository tracking the pass repo
	GitRepository SecretSourceGitRepositoryRef `json:"gitRepository"`
	PrivateKey    SecretSourcePrivateKeyRef    `json:"privateKey"`
}

// SecretSourceStatus defines the observed state of SecretSource
type SecretSourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SecretSource is the Schema for the secretsources API
type SecretSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretSourceSpec   `json:"spec,omitempty"`
	Status SecretSourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SecretSourceList contains a list of SecretSource
type SecretSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretSource{}, &SecretSourceList{})
}
