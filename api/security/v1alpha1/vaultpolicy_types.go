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

// VaultPolicySpec defines the desired state of VaultPolicy
type VaultPolicySpec struct {
	// vaultPolicyName is the name of the policy in Vault.
	// +kubebuilder:validation:Required
	VaultPolicyName string `json:"vaultPolicyName"`

	// policy is the HCL policy string.
	// +kubebuilder:validation:Required
	Policy string `json:"policy"`
}

// VaultPolicyStatus defines the observed state of VaultPolicy
type VaultPolicyStatus struct {
	// conditions represent the current state of the VaultPolicy resource.
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

// VaultPolicy is the Schema for the vaultpolicies API
type VaultPolicy struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of VaultPolicy
	// +required
	Spec VaultPolicySpec `json:"spec"`

	// status defines the observed state of VaultPolicy
	// +optional
	Status VaultPolicyStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// VaultPolicyList contains a list of VaultPolicy
type VaultPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []VaultPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VaultPolicy{}, &VaultPolicyList{})
}
