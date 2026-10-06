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

// ConfigSyncSpec defines the desired state of ConfigSync.
type ConfigSyncSpec struct {
    // TargetNamespace is where the config maps will be synced
    // +kubebuilder:validation:Required
    TargetNamespace string `json:"targetNamespace"`

    // ConfigData holds key-value pairs of configuration data
    // +kubebuilder:validation:Required
    ConfigData map[string]string `json:"configData,omitempty"`

    // TargetDeployments is a list of deployment names to restart when config changes
    TargetDeployments []string `json:"targetDeployments,omitempty"`
}

type ConfigSyncStatus struct {
	// Phase represents the current state of the sync (e.g., Pending, Synced, Failed)
	Phase string `json:"phase,omitempty"`

	// Message provides details about the last reconciliation attempt
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ConfigSync is the Schema for the configsyncs API.
type ConfigSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigSyncSpec   `json:"spec,omitempty"`
	Status ConfigSyncStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ConfigSyncList contains a list of ConfigSync.
type ConfigSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigSync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigSync{}, &ConfigSyncList{})
}
