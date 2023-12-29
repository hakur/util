package fsm

import (
	"strconv"
)

type ICondition interface {
	Compare(parameters map[string]*Parameter) bool
}

type CompareType = string

const (
	CompareTypeEqual        = "=="
	CompareTypeNotEqual     = "!="
	CompareTypeLess         = "<"
	CompareTypeLessEuqal    = "<="
	CompareTypeGreater      = ">"
	CompareTypeGreaterEqual = ">="
)

type Condition struct {
	CompareType   CompareType
	ParameterName string
	Value         string
}

func (t *Condition) Compare(parameters map[string]*Parameter) bool {
	parameter, ok := parameters[t.ParameterName]
	if !ok {
		return false
	}

	switch parameter.Type {
	case ParameterTypeBool:
		return t.CompareBool(parameter.Value)
	case ParameterTypeFloat:
		return t.CompareFloat(parameter.Value)
	case ParameterTypeInt:
		return t.CompareInt(parameter.Value)
	case ParameterTypeString:
		return t.CompareString(parameter.Value)
	}
	return false
}

func (t *Condition) CompareBool(value string) bool {
	pv, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	v, err := strconv.ParseBool(t.Value)
	if err != nil {
		return false
	}

	switch t.CompareType {
	case CompareTypeEqual:
		return pv == v
	case CompareTypeNotEqual:
		return pv != v
	}

	return false
}

func (t *Condition) CompareFloat(value string) bool {
	pv, err := strconv.ParseFloat(value, 32)
	if err != nil {
		println(err.Error())
		return false
	}
	v, err := strconv.ParseFloat(t.Value, 32)
	if err != nil {
		println(err.Error())
		return false
	}

	switch t.CompareType {
	case CompareTypeEqual:
		return pv == v
	case CompareTypeNotEqual:
		return pv != v
	case CompareTypeGreater:
		return pv > v
	case CompareTypeGreaterEqual:
		return pv >= v
	case CompareTypeLess:
		return pv < v
	case CompareTypeLessEuqal:
		return pv <= v
	}

	return false
}

func (t *Condition) CompareInt(value string) bool {
	pv, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	v, err := strconv.Atoi(t.Value)
	if err != nil {
		return false
	}

	switch t.CompareType {
	case CompareTypeEqual:
		return pv == v
	case CompareTypeNotEqual:
		return pv != v
	case CompareTypeGreater:
		return pv > v
	case CompareTypeGreaterEqual:
		return pv >= v
	case CompareTypeLess:
		return pv < v
	case CompareTypeLessEuqal:
		return pv <= v
	}

	return false
}

func (t *Condition) CompareString(value string) bool {
	switch t.CompareType {
	case CompareTypeEqual:
		return value == t.Value
	case CompareTypeNotEqual:
		return value != t.Value
	}

	return false
}
