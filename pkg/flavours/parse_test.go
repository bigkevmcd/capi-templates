package flavours

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestParseFile(t *testing.T) {
	c, err := ParseFile("testdata/template1.yaml")
	if err != nil {
		t.Fatal(err)
	}

	want := &CAPITemplate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CAPITemplate",
			APIVersion: "capi.weave.works/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-template",
		},
		Spec: CAPITemplateSpec{
			ResourceTemplates: []CAPIResourceTemplate{},
		},
	}
	if diff := cmp.Diff(want, c, cmpopts.IgnoreFields(CAPITemplateSpec{}, "ResourceTemplates")); diff != "" {
		t.Fatalf("failed to read the template:\n%s", diff)
	}
}

func TestParams(t *testing.T) {
	paramTests := []struct {
		filename string
		want     []string
	}{
		{
			filename: "testdata/template1.yaml",
			want:     []string{"CLUSTER_NAME"},
		},
		{
			filename: "testdata/template2.yaml",
			want:     []string{"CLUSTER_NAME", "AWS_NODE_MACHINE_TYPE", "AWS_SSH_KEY_NAME"},
		},
	}

	for _, tt := range paramTests {
		t.Run(tt.filename, func(t *testing.T) {
			c, err := ParseFile(tt.filename)
			if err != nil {
				t.Fatal(err)
			}

			params, err := Params(c.Spec)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, params); diff != "" {
				t.Fatalf("failed to extract params:\n%s", diff)
			}
		})
	}
}
