package flavours

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	processor "sigs.k8s.io/cluster-api/cmd/clusterctl/client/yamlprocessor"
	sigsyaml "sigs.k8s.io/yaml"
)

var (
	serializer  = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	metadataGVK = &schema.GroupVersionKind{Group: "capi.weave.works", Version: "templates", Kind: "TemplateMetadata"}
)

// ParseMultiDoc takes a single template with multiple YAML documents, and
// parses the parameters from it.
//
// If a TemplateMetadata object is found, the parameters returned are enriched
// with the data.
func ParseMultiDoc(fname string) ([]Param, error) {
	b, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}
	proc := processor.NewSimpleProcessor()

	variables := map[string]bool{}
	metas := map[string]Param{}
	for _, v := range bytes.Split(b, []byte("---\n")) {
		if len(v) == 0 {
			continue
		}
		var data unstructured.Unstructured
		_, gvk, err := serializer.Decode(v, nil, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s to unstructured %s: %w", fname, v, err)
		}
		if gvk.String() == metadataGVK.String() {
			m, err := parseMetadata(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse metadata from TemplateMetadata in file %s: %w", fname, err)
			}
			if m != nil {
				metas = m
			}
		}
		tv, err := proc.GetVariables(v)
		if err != nil {
			return nil, fmt.Errorf("processing template in file %s: %w", fname, err)
		}
		for _, n := range tv {
			variables[n] = true
		}
	}

	var result []Param
	for k := range variables {
		p, ok := metas[k]
		if !ok {
			result = append(result, Param{Name: k})
			continue
		}
		result = append(result, p)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

// TODO This should really use a serializer with the correct scheme registered.
func parseMetadata(b []byte) (map[string]Param, error) {
	var ct CAPITemplateMetadata
	if err := sigsyaml.Unmarshal(b, &ct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	params := map[string]Param{}
	for _, v := range ct.Spec.Params {
		params[v.Name] = Param{
			Name:        v.Name,
			Description: v.Description,
		}
	}
	return params, nil
}
