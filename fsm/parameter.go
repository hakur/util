package fsm

type ParameterType = string

const (
	ParameterTypeBool   ParameterType = "bool"
	ParameterTypeString ParameterType = "string"
	ParameterTypeFloat  ParameterType = "float"
	ParameterTypeInt    ParameterType = "int"
)

type Parameter struct {
	Name  string
	Value string
	Type  ParameterType
}
