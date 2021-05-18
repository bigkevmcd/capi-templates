package flavours

import (
	"fmt"
	"io/fs"
	"os"
	"sort"

	processor "sigs.k8s.io/cluster-api/cmd/clusterctl/client/yamlprocessor"
	"sigs.k8s.io/yaml"
)

func ParseFile(fname string) (*CAPITemplate, error) {
	b, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}
	return ParseBytes(b, fname)
}

func ParseFileFromFS(fsys fs.FS, fname string) (*CAPITemplate, error) {
	b, err := fs.ReadFile(fsys, fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}
	return ParseBytes(b, fname)
}

func ParseBytes(b []byte, name string) (*CAPITemplate, error) {
	var t CAPITemplate
	err := yaml.Unmarshal(b, &t)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", name, err)
	}
	return &t, nil
}

// Params extracts the named parameters from resource templates in a spec.
func Params(s CAPITemplateSpec) ([]string, error) {
	proc := processor.NewSimpleProcessor()
	variables := map[string]bool{}
	for _, v := range s.ResourceTemplates {
		tv, err := proc.GetVariables(v.RawExtension.Raw)
		if err != nil {
			return nil, fmt.Errorf("processing template: %w", err)
		}
		for _, n := range tv {
			variables[n] = true
		}
	}
	var names []string
	for k := range variables {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, nil
}
