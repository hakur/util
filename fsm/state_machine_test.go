package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateMachineTransit(t *testing.T) {
	sm := NewStateMachine()
	sm.SetStateUpdatedCallback(func(from, to State) {
		println(fmt.Sprintf("--- state updated from %s to %s", from, to))
	})

	sm.AddState("idle", "walk", "run")
	parameter := &Parameter{Name: "speed", Value: "0", Type: ParameterTypeFloat}
	sm.AddParameter(parameter)

	sm.Transit()
	assert.Equal(t, "idle", sm.CurrentState)

	sm.AddTransition(&Transition{Name: "idle_walk", From: "idle", To: "walk", Conditions: map[string]ICondition{
		"group id": &ConditionGroup{
			Conditions: map[string]ICondition{
				"random id": &Condition{
					CompareType:   CompareTypeGreater,
					Value:         "0",
					ParameterName: "speed",
				},
				"random id 2": &Condition{
					CompareType:   CompareTypeLess,
					Value:         "5",
					ParameterName: "speed",
				},
			},
			CompareType: ConditionGroupCompareTypeAnd,
		},
	}}, parameter)
	sm.AddTransition(&Transition{Name: "walk_run", From: "walk", To: "run", Conditions: map[string]ICondition{
		"random id": &Condition{
			CompareType:   CompareTypeGreater,
			Value:         "5",
			ParameterName: "speed",
		},
	}}, parameter)

	sm.SetParameterValue("speed", "4.9")
	assert.Equal(t, "walk", sm.CurrentState)

	sm.SetParameterValue("speed", "10")
	assert.Equal(t, "run", sm.CurrentState)
}

func BenchmarkTransit(b *testing.B) {
	sm := NewStateMachine()

	sm.AddState("idle", "walk", "run")
	parameter := &Parameter{Name: "speed", Value: "0", Type: ParameterTypeFloat}
	sm.AddParameter(parameter)
	sm.Transit()

	sm.AddTransition(&Transition{Name: "idle_walk", From: "idle", To: "walk", Conditions: map[string]ICondition{
		"group id": &ConditionGroup{
			Conditions: map[string]ICondition{
				"random id": &Condition{
					CompareType:   CompareTypeGreater,
					Value:         "0",
					ParameterName: "speed",
				},
				"random id 2": &Condition{
					CompareType:   CompareTypeLess,
					Value:         "5",
					ParameterName: "speed",
				},
			},
		},
	}}, parameter)
	sm.AddTransition(&Transition{Name: "walk_run", From: "walk", To: "run", Conditions: map[string]ICondition{
		"random id": &Condition{
			CompareType:   CompareTypeGreater,
			Value:         "5",
			ParameterName: "speed",
		},
	}}, parameter)

	sm.SetParameterValue("speed", "4.9")
	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		sm.Transit()
	}
}
