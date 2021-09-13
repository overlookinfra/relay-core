package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WebhookTrigger represents a definition of a webhook to receive events.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
type WebhookTrigger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WebhookTriggerSpec `json:"spec"`

	// +optional
	Status WebhookTriggerStatus `json:"status,omitempty"`
}

type WebhookTriggerSpec struct {
	// TenantRef selects the tenant to apply this trigger to.
	TenantRef corev1.LocalObjectReference `json:"tenantRef"`

	// Name is a friendly name for this webhook trigger used for authentication
	// and reporting.
	//
	// +optional
	Name string `json:"name,omitempty"`

	// Container defines the properties of the Docker container to run.
	Container `json:",inline"`
}

type WebhookTriggerStatus struct {
	// ObservedGeneration is the generation of the resource specification that
	// this status matches.
	//
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Namespace is the Kubernetes namespace containing the target resources of
	// this webhook trigger.
	//
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// URL is the endpoint for the webhook once provisioned.
	//
	// +optional
	URL string `json:"url,omitempty"`

	// Conditions are the observations of this resource's tate.
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []WebhookTriggerCondition `json:"conditions,omitempty"`
}

type WebhookTriggerConditionType string

const (
	WebhookTriggerServiceReady WebhookTriggerConditionType = "ServiceReady"
	WebhookTriggerReady        WebhookTriggerConditionType = "Ready"
)

type WebhookTriggerCondition struct {
	Condition `json:",inline"`

	// Type is the identifier for this condition.
	//
	// +kubebuilder:validation:Enum=ServiceReady;Ready
	Type WebhookTriggerConditionType `json:"type"`
}

// WebhookTriggerList enumerates many WebhookTrigger resources.
//
// +kubebuilder:object:root=true
type WebhookTriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebhookTrigger `json:"items"`
}
