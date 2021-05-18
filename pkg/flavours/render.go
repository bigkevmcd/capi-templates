package flavours

import (
	"fmt"

	processor "sigs.k8s.io/cluster-api/cmd/clusterctl/client/yamlprocessor"
	"sigs.k8s.io/yaml"
)

func Render(t *CAPITemplate, vars map[string]string) ([][]byte, error) {
	proc := processor.NewSimpleProcessor()
	var processed [][]byte
	for _, v := range t.Spec.ResourceTemplates {
		b, err := proc.Process(v.RawExtension.Raw, func(n string) (string, error) {
			if s, ok := vars[n]; ok {
				return s, nil
			}
			return "", fmt.Errorf("variable %s not found", n)
		})
		if err != nil {
			return nil, fmt.Errorf("processing template: %w", err)
		}
		data, err := yaml.JSONToYAML(b)
		if err != nil {
			return nil, fmt.Errorf("failed to convert back to YAML: %w", err)
		}
		processed = append(processed, data)
	}
	return processed, nil
}
