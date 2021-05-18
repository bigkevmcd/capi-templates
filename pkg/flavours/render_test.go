package flavours

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRender(t *testing.T) {
	parsed := mustParseFile(t, "testdata/template3.yaml")

	b, err := Render(parsed, map[string]string{"CLUSTER_NAME": "testing"})
	if err != nil {
		t.Fatal(err)
	}

	want := `---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: Cluster
metadata:
  name: testing
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSMachineTemplate
metadata:
  name: testing-md-0
`
	if diff := cmp.Diff(want, writeMultiDoc(t, b)); diff != "" {
		t.Fatalf("rendering failure:\n%s", diff)
	}
}

func TestRender_unknown_parameter(t *testing.T) {
	parsed := mustParseFile(t, "testdata/template3.yaml")

	_, err := Render(parsed, map[string]string{})
	assertErrorMatch(t, "value for variables.*CLUSTER_NAME.*is not set", err)
}

func writeMultiDoc(t *testing.T, objs [][]byte) string {
	t.Helper()
	var out bytes.Buffer
	for _, v := range objs {
		if _, err := out.Write([]byte("---\n")); err != nil {
			t.Fatal(err)
		}
		if _, err := out.Write(v); err != nil {
			t.Fatal(err)
		}
	}
	return out.String()
}

func assertErrorMatch(t *testing.T, s string, e error) {
	t.Helper()
	if s == "" && e == nil {
		return
	}
	if s != "" && e == nil {
		t.Fatalf("wanted error matching %s, got nil", s)
	}
	match, err := regexp.MatchString(s, e.Error())
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("failed to match error %s against %s", e, s)
	}
}
