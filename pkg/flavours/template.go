package flavours

import (
	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TemplateParam provides metadata for a template parameter.
type TemplateParam struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CAPITemplateSpec holds the desired state of CAPITemplate.
type CAPITemplateSpec struct {
	Description       string                 `json:"description,omitempty"`
	Params            []TemplateParam        `json:"params,omitempty"`
	ResourceTemplates []CAPIResourceTemplate `json:"resourcetemplates,omitempty"`
}

// CAPIResourceTemplate describes a resource to create.
type CAPIResourceTemplate struct {
	runtime.RawExtension `json:",inline"`
}

// CAPITemplateStatus describes the desired state of CAPITemplate
type CAPITemplateStatus struct{}

// CAPITemplate takes parameters and uses them to create CAPI templates.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type CAPITemplate struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec holds the desired state of the CAPITemplate from the client
	// +optional
	Spec CAPITemplateSpec `json:"spec"`
	// +optional
	Status CAPITemplateStatus `json:"status,omitempty"`
}

type CAPITemplateMetadata struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec holds the desired state of the CAPITemplate from the client
	// +optional
	Spec CAPITemplateMetadataSpec `json:"spec"`
}

// CAPITemplateMetadataSpec holds the desired state of CAPITemplateMetadata.
type CAPITemplateMetadataSpec struct {
	Description string          `json:"description,omitempty"`
	Params      []TemplateParam `json:"params,omitempty"`
}
