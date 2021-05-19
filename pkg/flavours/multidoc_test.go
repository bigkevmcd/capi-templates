package flavours

import (
	"testing"

	"github.com/bigkevmcd/capi-templates/test"
	"github.com/google/go-cmp/cmp"
)

func TestParseMultiDoc_errors(t *testing.T) {
	parseTests := []struct {
		filename string
		errMsg   string
	}{
		{"testdata/bad_multidoc.yaml", "processing template.*bad substitution"},
		{"testdata/bad_multidoc_corrupt_metadata.yaml", "failed to unmarshal testdata/bad_multidoc_corrupt_metadata.yaml to unstructured"},
		{"testdata/bad_multidoc_invalid_metadata.yaml", "failed to parse metadata from TemplateMetadata"},
		{"testdata/unknownyaml", "failed to read template"},
	}
	for _, tt := range parseTests {
		t.Run(tt.filename, func(t *testing.T) {
			_, err := ParseMultiDoc(tt.filename)
			test.AssertErrorMatch(t, tt.errMsg, err)
		})
	}
}

func TestParseMultiDoc_with_no_metadata(t *testing.T) {
	params, err := ParseMultiDoc("testdata/multidoc_with_no_metadata.yaml")
	if err != nil {
		t.Fatal(err)
	}
	want := []Param{
		{Name: "AWS_CONTROL_PLANE_MACHINE_TYPE"},
		{Name: "AWS_NODE_MACHINE_TYPE"},
		{Name: "AWS_REGION"},
		{Name: "AWS_SSH_KEY_NAME"},
		{Name: "CLUSTER_NAME"},
		{Name: "CONTROL_PLANE_MACHINE_COUNT"},
		{Name: "KUBERNETES_VERSION"},
		{Name: "WORKER_MACHINE_COUNT"},
	}
	if diff := cmp.Diff(want, params); diff != "" {
		t.Fatalf("parsing failed:\n%s", diff)
	}
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
