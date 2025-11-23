package toolspec

import (
	"github.com/alexander-lindner/go-cff"
	"gopkg.in/yaml.v3"
)

type SpecFile struct {
	Tools map[string]ToolSpec `yaml:"tools"`
}

type ToolSpec struct {
	ID          string                   `json:"id" yaml:"-"`
	Name        string                   `json:"name" yaml:"-"`
	Title       string                   `json:"title" yaml:"title"`
	Description string                   `json:"description" yaml:"description"`
	Parameters  map[string]ParameterSpec `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Data        map[string]DataSpec      `json:"data,omitempty" yaml:"data,omitempty"`
	Citation    cff.Cff                  `json:"citation,omitempty" yaml:"-"`
}

type ParameterSpec struct {
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	ToolType    string      `json:"type" yaml:"type"`
	IsArray     bool        `json:"array,omitempty" yaml:"array,omitempty" default:"false"`
	Default     interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	Values      []string    `json:"values,omitempty" yaml:"values,omitempty"`
	Min         float64     `json:"min,omitempty" yaml:"min,omitempty"`
	Max         float64     `json:"max,omitempty" yaml:"max,omitempty"`
	Optional    bool        `json:"optional,omitempty" yaml:"optional,omitempty"`
}

type DataSpec struct {
	Path        string      `json:"path" yaml:"path"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Example     string      `json:"example,omitempty" yaml:"example,omitempty"`
	Extension   interface{} `json:"-" yaml:"extension,omitempty"`
	Extensions  []string    `json:"extension,omitempty"`
}

func (d *DataSpec) UnmarshalYAML(value *yaml.Node) error {
	type dataSpecAlias DataSpec
	var alias dataSpecAlias
	if err := value.Decode(&alias); err != nil {
		return err
	}

	*d = DataSpec(alias)

	switch ext := d.Extension.(type) {
	case string:
		d.Extensions = []string{ext}
	case []interface{}:
		d.Extensions = make([]string, len(ext))
		for i, e := range ext {
			d.Extensions[i] = e.(string)
		}
	}

	d.Extension = nil
	return nil
}

