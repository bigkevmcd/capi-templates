package flavours

import (
	"os"
	"testing"
	"io/ioutil"

	"github.com/bigkevmcd/capi-templates/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
			Description: "this is test template 1",
			Params: []TemplateParam{
				{
					Name:        "CLUSTER_NAME",
					Description: "This is used for the cluster naming.",
				},
			},
			ResourceTemplates: []CAPIResourceTemplate{},
		},
	}
	if diff := cmp.Diff(want, c, cmpopts.IgnoreFields(CAPITemplateSpec{}, "ResourceTemplates")); diff != "" {
		t.Fatalf("failed to read the template:\n%s", diff)
	}
}

func TestParseFile_with_unknown_file(t *testing.T) {
	_, err := ParseFile("testdata/unknownyaml")
	test.AssertErrorMatch(t, "failed to read template", err)
}

func TestParseFileFromFS(t *testing.T) {
	c, err := ParseFileFromFS(os.DirFS("testdata"), "template2.yaml")
	if err != nil {
		t.Fatal(err)
	}

	want := &CAPITemplate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CAPITemplate",
			APIVersion: "capi.weave.works/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-template2",
		},
		Spec: CAPITemplateSpec{
			Description: "this is test template 2",
			Params: []TemplateParam{
				{
					Name:        "AWS_SSH_KEY_NAME",
					Description: "A description",
				},
			},

			ResourceTemplates: []CAPIResourceTemplate{},
		},
	}
	if diff := cmp.Diff(want, c, cmpopts.IgnoreFields(CAPITemplateSpec{}, "ResourceTemplates")); diff != "" {
		t.Fatalf("failed to read the template:\n%s", diff)
	}
}

func TestParseFileFromFS_with_unknown_file(t *testing.T) {
	_, err := ParseFileFromFS(os.DirFS("testdata"), "unknown.yaml")
	test.AssertErrorMatch(t, "failed to read template", err)
}


func TestParseConfigmap(t *testing.T) {
	cmBytes, err := ioutil.ReadFile("testdata/configmap1.yaml")
    if err != nil {
		t.Fatal(err)
	}

	var data unstructured.Unstructured
	cm, gvk, err := serializer.Decode(cmBytes, nil, &data)
	if err != nil {
		t.Fatal(err)
	}

	tm, err := ParseConfigmap(cm)
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
			Description: "this is test template 1",
			Params: []TemplateParam{
				{
					Name:        "CLUSTER_NAME",
					Description: "This is used for the cluster naming.",
				},
			},
			ResourceTemplates: []CAPIResourceTemplate{},
		},
	}
	if diff := cmp.Diff(want, tm["template1"], cmpopts.IgnoreFields(CAPITemplateSpec{}, "ResourceTemplates")); diff != "" {
		t.Fatalf("failed to read the template from the configmap:\n%s", diff)
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
			want: []string{
				"AWS_NODE_MACHINE_TYPE",
				"AWS_SSH_KEY_NAME",
				"CLUSTER_NAME",
			},
		},
	}

	for _, tt := range paramTests {
		t.Run(tt.filename, func(t *testing.T) {
			c := mustParseFile(t, tt.filename)
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
