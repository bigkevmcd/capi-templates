package fs

import (
	"os"
	"testing"

	"github.com/bigkevmcd/capi-templates/pkg/flavours"
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
