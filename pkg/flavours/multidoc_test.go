package flavours

import (
	"testing"

	"github.com/bigkevmcd/capi-templates/test"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

var (
	serializer  = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	metadataGVK = &schema.GroupVersionKind{Group: "capi.weave.works", Version: "templates", Kind: "TemplateMetadata"}
)

func TestParseMultiDoc_with_invalid_yaml(t *testing.T) {
	t.Skip()
}

func TestParseMultiDoc_with_unknown_file(t *testing.T) {
	_, err := ParseMultiDoc("testdata/unknownyaml")
	test.AssertErrorMatch(t, "failed to read template", err)
}

func TestParseMultiDoc_with_no_metadata(t *testing.T) {
	t.Skip()
}

func TestParseMultiDoc(t *testing.T) {
	params, err := ParseMultiDoc("testdata/multidoc.yaml")
	if err != nil {
		t.Fatal(err)
	}
	want := []Param{
		{Name: "AWS_CONTROL_PLANE_MACHINE_TYPE"},
		{Name: "AWS_NODE_MACHINE_TYPE"},
		{Name: "AWS_REGION"},
		{Name: "AWS_SSH_KEY_NAME"},
		{Name: "CLUSTER_NAME", Description: "This is used for the cluster naming."},
		{Name: "CONTROL_PLANE_MACHINE_COUNT"},
		{Name: "KUBERNETES_VERSION"},
		{Name: "WORKER_MACHINE_COUNT"},
	}
	if diff := cmp.Diff(want, params); diff != "" {
		t.Fatalf("parsing failed:\n%s", diff)
	}
}
