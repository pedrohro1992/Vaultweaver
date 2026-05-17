/*
Copyright 2026.

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

// VaultKubernetesRoleBindingSpec defines the desired state of VaultKubernetesRoleBinding
type VaultKubernetesRoleBindingSpec struct {
	// authMount is the Vault auth mount point (e.g., "kubernetes").
	// +kubebuilder:validation:Required
	// +kubebuilder:default="kubernetes"
	AuthMount string `json:"authMount"`

	// roleName is the name of the Vault role to create/update.
	// +kubebuilder:validation:Required
	RoleName string `json:"roleName"`

	// boundServiceAccounts is a list of Kubernetes ServiceAccount names that are allowed to login.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	BoundServiceAccounts []string `json:"boundServiceAccounts"`

	// boundNamespaces is a list of Kubernetes namespaces that are allowed to login.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	BoundNamespaces []string `json:"boundNamespaces"`

	// tokenPolicies is a list of Vault policies to assign to tokens issued under this role.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	TokenPolicies []string `json:"tokenPolicies"`

	// tokenTTL is the TTL of the Vault token.
	// +optional
	TokenTTL string `json:"tokenTTL,omitempty"`

	// audience is the expected audience for the ServiceAccount JWT.
	// +optional
	Audience string `json:"audience,omitempty"`
}

// VaultKubernetesRoleBindingStatus defines the observed state of VaultKubernetesRoleBinding.
type VaultKubernetesRoleBindingStatus struct {
	// conditions represent the current state of the VaultKubernetesRoleBinding resource.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// observedGeneration is the last reconciled generation of the resource.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VaultKubernetesRoleBinding is the Schema for the vaultkubernetesrolebindings API
type VaultKubernetesRoleBinding struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of VaultKubernetesRoleBinding
	// +required
	Spec VaultKubernetesRoleBindingSpec `json:"spec"`

	// status defines the observed state of VaultKubernetesRoleBinding
	// +optional
	Status VaultKubernetesRoleBindingStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// VaultKubernetesRoleBindingList contains a list of VaultKubernetesRoleBinding
type VaultKubernetesRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []VaultKubernetesRoleBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VaultKubernetesRoleBinding{}, &VaultKubernetesRoleBindingList{})
}
