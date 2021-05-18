package fs

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/bigkevmcd/capi-templates/pkg/flavours"
)

// New creates and returns a new Filesystem based template library.
func New(fs fs.FS, base string) *FSLibrary {
	return &FSLibrary{FS: fs, base: base}
}

type FSLibrary struct {
	fs.FS
	base string
}

func (f FSLibrary) Flavours() ([]*flavours.Flavour, error) {
	dirs, err := fs.ReadDir(f.FS, f.base)
	if err != nil {
		return nil, fmt.Errorf("failed to ReadDir in Flavours(): %w", err)
	}

	var found []*flavours.Flavour
	for _, v := range dirs {
		if !v.IsDir() {
			t, err := flavours.ParseFileFromFS(f.FS, filepath.Join(f.base, v.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to parse: %w", err)
			}

			params, err := paramsFromSpec(t.Spec)
			if err != nil {
				return nil, err
			}
			found = append(found, &flavours.Flavour{
				Name:        t.ObjectMeta.Name,
				Description: t.Spec.Description,
				Params:      params,
			})
		}
	}
	sort.Slice(found, func(i, j int) bool { return found[i].Name < found[j].Name })
	return found, nil
}

func paramsFromSpec(s flavours.CAPITemplateSpec) ([]flavours.Param, error) {
	paramNames, err := flavours.Params(s)
	if err != nil {
		return nil, fmt.Errorf("failed to get params from template: %w", err)
	}
	paramsMeta := map[string]flavours.Param{}
	for _, v := range paramNames {
		paramsMeta[v] = flavours.Param{Name: v}
	}

	for _, v := range s.Params {
		if m, ok := paramsMeta[v.Name]; ok {
			m.Description = v.Description
			paramsMeta[v.Name] = m
		}
	}

	var params []flavours.Param
	for _, v := range paramsMeta {
		params = append(params, v)
	}
	return params, nil
}
