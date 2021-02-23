package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/code-generator"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=teployments,singular=teployment,shortName=teploy,categories={}
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Teployment describes a teployment.
type Teployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TeploymentSpec `json:"spec"`
}

// TeploymentSpec is the spec for a teployment resource
type TeploymentSpec struct {
	Replicas *int32 `json:replicas`
	ServiceType string 	`json:serviceType`
	NodePort *int `json:nodePort`
	Image string `json:image`
	ContainerPort int `json:containerPort`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TeploymentList is a list of Teployment resources
type TeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Teployment `json:"items"`
}