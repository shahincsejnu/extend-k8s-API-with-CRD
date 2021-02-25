package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus

// +groupName=shahin.oka.com
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=string
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Phase,type=string
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name=Age,type=date
// +kubebuilder:resource:path=teployments,singular=teployment,shortName=teploy,categories={}
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// Teployment describes a teployment. It is our root type, it describes the Teployment kind. It contains TypeMeta (which describes API version and Kind), and also contains ObjectMeta, which holds things like name, namespace, and labels.
type Teployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TeploymentSpec `json:"spec"`
	// +optional
	Status TeploymentStatus `json:"status"`
}

// TeploymentSpec is the spec for a teployment resource, it defines the desired state of Teployment
type TeploymentSpec struct {
	// +optional
	// +kubebuilder:default:=1
	Replicas *int32 `json:"replicas"`
	ServiceType string 	`json:"serviceType"`
	NodePort int `json:"nodePort,omitempty"`
	Image string `json:"image"`
	ContainerPort int `json:"containerPort"`
}

// TeploymentStatus defines the observed state of Teployment
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

// TeploymentList is a list of Teployment resources. In general, we never modify either of these -- all modifications go in either Spec or Status.
type TeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Teployment `json:"items"`
}