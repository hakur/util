package fsm

type Transition struct {
	Name       string
	Conditions map[string]ICondition
	From       State
	To         State
}

func (t *Transition) AddCondition(conditionName string, condition ICondition) {

}

func (t *Transition) RemoveCondition(conditionName string) {

}

// Transit 尝试检查状态装换条件，如果满足其中一个条件就认为是可以转换的，那么将返回要切换到的目标状态，否则返回空
// 如果没有设置任何条件则直接返回目标状态
func (t *Transition) Transit(parameters map[string]*Parameter) (to State) {
	if len(t.Conditions) > 0 {
		for _, c := range t.Conditions {
			if c.Compare(parameters) {
				return t.To
			}
		}
	} else {
		return t.To
	}
	return
}
