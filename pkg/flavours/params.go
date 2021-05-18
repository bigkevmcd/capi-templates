package flavours

import "fmt"

// ParamsFromSpec extracts the named parameters from a CAPITemplate, finding all
// the named parameters in each of the resource templates, and enriching that
// with data from the params field in the spec (if found).
//
// Any fields in the templates, but not in the params will not be enriched, and
// only the name will be returned.
func ParamsFromSpec(s CAPITemplateSpec) ([]Param, error) {
	paramNames, err := Params(s)
	if err != nil {
		return nil, fmt.Errorf("failed to get params from template: %w", err)
	}
	paramsMeta := map[string]Param{}
	for _, v := range paramNames {
		paramsMeta[v] = Param{Name: v}
	}

	for _, v := range s.Params {
		if m, ok := paramsMeta[v.Name]; ok {
			m.Description = v.Description
			paramsMeta[v.Name] = m
		}
	}

	var params []Param
	for _, v := range paramsMeta {
		params = append(params, v)
	}
	return params, nil
}
