package toolspec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func (s *SpecFile) GetTool(toolName string) (ToolSpec, error) {
	toolSpec, ok := s.Tools[toolName]
	if !ok {
		return ToolSpec{}, fmt.Errorf("tool %s was not found in the given specification file", toolName)
	}

	return toolSpec, nil
}

func LoadToolSpec(rawData []byte) (SpecFile, error) {
	var toolSpec SpecFile
	err := yaml.Unmarshal(rawData, &toolSpec)
	if err != nil {
		return SpecFile{}, err
	}

	return toolSpec, nil
}

