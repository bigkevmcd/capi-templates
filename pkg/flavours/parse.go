package flavours

import (
	"fmt"
	"io/fs"
	"os"
	"sort"

	processor "sigs.k8s.io/cluster-api/cmd/clusterctl/client/yamlprocessor"
	"sigs.k8s.io/yaml"
	corev1 "k8s.io/api/core/v1"
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

// ParseConfigmap returns a map of CAPITemplates indexed by their name.
// The name of the template is set to the key of the Configmap.Data map.
func ParseConfigmap(cm corev1.Configmap) (map[string]*CAPITemplate, error) {
	tm := map[string]*CAPITemplate{}

	for k, v := range(cm) {
		t, err := ParseBytes([]byte(v), k)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal template %s from configmap %s, err: %w", k, cm.Metadata.Name, err)
		}
		tm[k] = t
	}
	return tm, nil
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
