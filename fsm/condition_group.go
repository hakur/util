package fsm

type ConditionGroupCompareType = string

const (
	ConditionGroupCompareTypeAnd = "and"
	ConditionGroupCompareTypeOr  = "or"
)

// ConditionGroup 条件组，组内条件必须都满足才能算是这个条件组满足了
type ConditionGroup struct {
	Conditions  map[string]ICondition
	CompareType ConditionGroupCompareType
}

func (t *ConditionGroup) Compare(parameters map[string]*Parameter) bool {
	var count int
	for _, c := range t.Conditions {
		if c.Compare(parameters) {
			if t.CompareType == ConditionGroupCompareTypeOr {
				return true
			}
			count++
		} else {
			return false
		}
	}

	return count == len(t.Conditions)
}
