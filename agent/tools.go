package agent

type Tools []Tool

func (tools Tools) Info() []ToolInfo {
	info := make([]ToolInfo, len(tools))
	for i, tool := range tools {
		info[i] = tool.Info()
	}
	return info
}

type ToolInfo struct {
	Type     string           `json:"type"`
	Function ToolInfoFunction `json:"function"`
}
type ToolInfoFunction struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Parameters  ToolInfoFunctionParameters `json:"parameters"`
}
type ToolInfoFunctionParameters struct {
	Type       string                                        `json:"type"`
	Properties map[string]ToolInfoFunctionParametersProperty `json:"properties"`
	Required   []string                                      `json:"required"`
}
type ToolInfoFunctionParametersProperty struct {
	Type        string            `json:"type"`
	Items       map[string]string `json:"items"`
	Description string            `json:"description"`
}
