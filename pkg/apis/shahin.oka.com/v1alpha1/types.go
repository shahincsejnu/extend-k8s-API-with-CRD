package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=teployments,singular=teployment,shortName=teploy,categories={}
// +kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=string
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Phase,type=string

// Teployment describes a teployment.
type Teployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TeploymentSpec `json:"spec"`
	// +optional
	Status TeploymentStatus `json:"status"`
}

// TeploymentSpec is the spec for a teployment resource
type TeploymentSpec struct {
	// +optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas"`
	ServiceType string 	`json:"serviceType"`
	NodePort int `json:"nodePort,omitempty"`
	Image string `json:"image"`
	ContainerPort int `json:"containerPort"`
}

type TeploymentStatus struct {
	// Specifies the current phase of the teployment
	// +optional
	Phase string `json:"phase"`

	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TeploymentList is a list of Teployment resources
type TeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Teployment `json:"items"`
}