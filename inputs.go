package toolspec

import (
	"encoding/json"
	"fmt"
)

type ToolInput struct {
	Parameters map[string]interface{} `json:"parameters"`
	Datasets   map[string]string      `json:"data"`
}

type InputFile map[string]ToolInput

func LoadInputs(rawData []byte) (InputFile, error) {
	var inputFile InputFile
	err := json.Unmarshal(rawData, &inputFile)
	if err != nil {
		return nil, err
	}

	return inputFile, nil
}

func (f InputFile) GetToolInput(toolName string) (ToolInput, error) {
	toolInput, ok := f[toolName]
	if !ok {
		return ToolInput{}, fmt.Errorf("tool input %s was not found in the inputs file", toolName)
	}

	return toolInput, nil
}

