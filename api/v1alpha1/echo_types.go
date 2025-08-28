/*
Copyright 2025.

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

// EchoPhase represents the current phase of Echo execution
// +kubebuilder:validation:Enum=Pending;Running;Completed;Failed
type EchoPhase string

const (
	// EchoPhasePending indicates the Echo is pending execution
	EchoPhasePending EchoPhase = "Pending"

	// EchoPhaseRunning indicates the Echo Job is currently running
	EchoPhaseRunning EchoPhase = "Running"

	// EchoPhaseCompleted indicates the Echo Job completed successfully
	EchoPhaseCompleted EchoPhase = "Completed"

	// EchoPhaseFailed indicates the Echo Job failed
	EchoPhaseFailed EchoPhase = "Failed"
)

// Condition types for Echo resources
const (
	// EchoConditionReady indicates whether the Echo is ready to be processed
	EchoConditionReady = "Ready"

	// EchoConditionProgressing indicates whether the Echo is progressing
	EchoConditionProgressing = "Progressing"

	// EchoConditionFailed indicates whether the Echo has failed
	EchoConditionFailed = "Failed"
)

// EchoSpec defines the desired state of Echo.
type EchoSpec struct {
	// Message is the text to echo. This field is required and cannot be empty.
	// The message will be printed by a Kubernetes Job running in a Pod.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1000
	Message string `json:"message"`
}

// EchoStatus defines the observed state of Echo.
type EchoStatus struct {
	// Phase represents the current phase of Echo execution
	// +optional
	Phase EchoPhase `json:"phase,omitempty"`

	// JobName is the name of the Kubernetes Job created to echo the message
	// +optional
	JobName string `json:"jobName,omitempty"`

	// Message provides human-readable information about the current state
	// +optional
	Message string `json:"message,omitempty"`

	// Conditions represent the latest available observations of Echo's current state
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// LastExecutionTime represents the last time the echo was executed
	// +optional
	LastExecutionTime *metav1.Time `json:"lastExecutionTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=echo
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".spec.message"
// +kubebuilder:printcolumn:name="Job",type="string",JSONPath=".status.jobName"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Echo is the Schema for the echoes API.
type Echo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EchoSpec   `json:"spec,omitempty"`
	Status EchoStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EchoList contains a list of Echo.
type EchoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Echo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Echo{}, &EchoList{})
}
