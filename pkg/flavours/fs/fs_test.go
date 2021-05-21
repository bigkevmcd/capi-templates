package fs

import (
	"os"
	"testing"

	"github.com/bigkevmcd/capi-templates/pkg/flavours"
	"github.com/bigkevmcd/capi-templates/test"
	"github.com/google/go-cmp/cmp"
)

func TestFlavours(t *testing.T) {
	f := New(os.DirFS("testdata"), "flavours")

	got, err := f.Flavours()
	if err != nil {
		t.Fatal(err)
	}
	want := []*flavours.Flavour{
		{
			Name:        "cluster-template1",
			Description: "this is test template 1",
			Version:     "1.2.3",
			Params: []flavours.Param{
				{
					Name:        "CLUSTER_NAME",
					Description: "This is used for the cluster naming.",
				},
			},
		},
		{
			Name:        "cluster-template1",
			Description: "this is test template 1",
			Version:     "2.1.0",
			Params: []flavours.Param{
				{
					Name:        "CLUSTER_NAME",
					Description: "This is used for the cluster naming.",
				},
			},
		},
		{
			Name:        "cluster-template2",
			Description: "this is test template 2",
			Version:     "1.2.3",
			Params: []flavours.Param{
				{
					Name: "AWS_NODE_MACHINE_TYPE",
				},
				{
					Name: "AWS_SSH_KEY_NAME",
				},
				{
					Name: "CLUSTER_NAME",
				},
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("failed to parse flavours:\n%s", diff)
	}
}

func TestFlavours_with_unknown_directory(t *testing.T) {
	f := New(os.DirFS("testdata"), "unknown")
	_, err := f.Flavours()
	test.AssertErrorMatch(t, "no such file or directory", err)
}

func TestFlavours_with_error_cases(t *testing.T) {
	errorTests := []struct {
		description string
		dirname     string
		errMsg      string
	}{
		{"invalid yaml", "badtemplates1", "failed to unmarshal badtemplates1/0.0.1/bad_template.yaml"},
		{"invalid params", "badtemplates2", "failed to get params from template"},
	}

	for _, tt := range errorTests {
		t.Run(tt.description, func(t *testing.T) {
			f := New(os.DirFS("testdata"), tt.dirname)
			_, err := f.Flavours()
			test.AssertErrorMatch(t, tt.errMsg, err)
		})
	}
}
