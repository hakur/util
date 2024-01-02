package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateMachineTransit(t *testing.T) {
	var err error
	sm := NewStateMachine("APP")
	sm.SetStateUpdatedCallback(func(from, to State) {
		println(fmt.Sprintf("--- state updated from %s to %s", from, to))
	})

	sm.AddValidTransition("idle", []State{"walk", "run"})
	sm.AddValidTransition("walk", []State{"idle", "run"})
	sm.AddValidTransition("run", []State{"walk"})
	parameter := &Parameter{Name: "speed", Value: "0", Type: ParameterTypeFloat}
	sm.AddParameter(parameter)

	sm.AutoTransit()
	assert.Equal(t, "idle", sm.CurrentState)

	err = sm.AddAutoTransition(&Transition{Name: "idle_walk", From: "idle", To: "walk", Conditions: map[string]ICondition{
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
	assert.Equal(t, nil, err)

	err = sm.AddAutoTransition(&Transition{Name: "walk_run", From: "walk", To: "run", Conditions: map[string]ICondition{
		"random id": &Condition{
			CompareType:   CompareTypeGreater,
			Value:         "5",
			ParameterName: "speed",
		},
	}}, parameter)
	assert.Equal(t, nil, err)

	sm.SetParameterValue("speed", "4.9")
	assert.Equal(t, "walk", sm.CurrentState)

	sm.SetParameterValue("speed", "10")
	assert.Equal(t, "run", sm.CurrentState)
}

func BenchmarkTransit(b *testing.B) {
	sm := NewStateMachine("APP")

	sm.AddValidTransition("idle", []State{"walk", "run"})
	sm.AddValidTransition("walk", []State{"idle", "run"})
	sm.AddValidTransition("run", []State{"walk"})
	parameter := &Parameter{Name: "speed", Value: "0", Type: ParameterTypeFloat}
	sm.AddParameter(parameter)
	sm.AutoTransit()

	sm.AddAutoTransition(&Transition{Name: "idle_walk", From: "idle", To: "walk", Conditions: map[string]ICondition{
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
	sm.AddAutoTransition(&Transition{Name: "walk_run", From: "walk", To: "run", Conditions: map[string]ICondition{
		"random id": &Condition{
			CompareType:   CompareTypeGreater,
			Value:         "5",
			ParameterName: "speed",
		},
	}}, parameter)

	sm.SetParameterValue("speed", "4.9")
	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		sm.AutoTransit()
	}
}

func TestGetMachine(t *testing.T) {
	sm := NewStateMachine("Player")
	motionSM := NewStateMachine("Motion")
	sm.AddSubMachine(motionSM)
	groundSM := NewStateMachine("Ground")
	motionSM.AddSubMachine(groundSM)

	flySM := NewStateMachine("Fly")
	motionSM.AddSubMachine(flySM)

	s := sm.GetMachine("Player/Motion/Ground")
	assert.Equal(t, groundSM.Name, s.Name)

	s = sm.GetMachine("/Player/Motion/Fly")
	assert.Equal(t, flySM.Name, s.Name)
}
