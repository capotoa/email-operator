package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EmailSpec defines the desired state of Email
type EmailSpec struct {
	SenderConfigRef string `json:"senderConfigRef"`
	RecipientEmail  string `json:"recipientEmail"`
	Subject         string `json:"subject"`
	Body            string `json:"body"`
}

// EmailStatus defines the observed state of Email
type EmailStatus struct {
	DeliveryStatus string `json:"deliveryStatus"`
	MessageID      string `json:"messageId"`
	Error          string `json:"error"`
}

// +kubebuilder:object:root=true

// Email is the Schema for the emails API
type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSpec   `json:"spec,omitempty"`
	Status EmailStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EmailList contains a list of Email
type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Email `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Email{}, &EmailList{})
}
