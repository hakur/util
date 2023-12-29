package fsm

import (
	"fmt"
	"strings"
	"sync"
)

type State = string

func NewStateMachine(name string) *StateMachine {
	t := new(StateMachine)
	t.Parameters = make(map[string]*Parameter)
	t.ParametersLink = make(map[*Parameter][]*Transition)
	t.Transitions = make(map[string]*Transition)
	t.SubMachines = make(map[string]*StateMachine)
	t.CurrentState = "Entry"
	t.Name = name
	return t
}

type StateMachine struct {
	lock     sync.RWMutex
	callback func(from State, to State)
	// Parameters 参数列表
	Parameters map[string]*Parameter
	// States 所有状态列表
	States []State
	// CurrentState 当前状态
	CurrentState State
	// Transitions 转换器列表
	Transitions map[string]*Transition
	// ParametersLink 用于自动触发转换
	ParametersLink map[*Parameter][]*Transition
	// SubMachines 内部子状态机
	SubMachines map[string]*StateMachine
	// Name 状态机的名称
	Name string
}

func (t *StateMachine) AddSubMachine(machine *StateMachine) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.SubMachines[machine.Name]; ok {
		return fmt.Errorf("sub state machine = %s already exists", machine.Name)
	}

	t.SubMachines[machine.Name] = machine
	return
}

// SetStateUpdatedCallback 设置状态发生切换时触发的回调函数
func (t *StateMachine) SetStateUpdatedCallback(callback func(from State, to State)) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.callback = callback
}

// GetCurrentState 返回状态机的当前状态
func (t *StateMachine) GetCurrentState() State {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.CurrentState
}

// GetParameter 取得状态切换参数对象
func (t *StateMachine) GetParameter(parameterName string) (parameter *Parameter) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	parameter = t.Parameters[parameterName]
	return
}

func (t *StateMachine) AddState(states ...State) (err error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for _, state := range states {
		for _, s := range t.States {
			if s == state {
				return fmt.Errorf("state=%s already exists", state)
			}
		}

		t.States = append(t.States, state)
	}

	return
}

// AddTransition 添加自动状态切换
func (t *StateMachine) AddTransition(trans *Transition, parameter *Parameter) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.Transitions[trans.Name]; ok {
		return fmt.Errorf("transition=%s already exists", trans.Name)
	}

	t.Transitions[trans.Name] = trans
	t.ParametersLink[parameter] = append(t.ParametersLink[parameter], trans)

	return
}

// Transit 手动检查所有的状态切换是否需要进行一次状态切换
func (t *StateMachine) Transit() {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.CurrentState == "Entry" {
		if len(t.States) > 0 {
			t.CurrentState = t.States[0]
		}
		return
	}

	var transitions []*Transition
	for _, trans := range t.Transitions {
		transitions = append(transitions, trans)
	}

	t.transit(transitions)
}

func (t *StateMachine) transit(transitions []*Transition) {
	for _, trans := range transitions {
		if toState := trans.Transit(t.Parameters); toState != "" && toState != t.CurrentState {
			var oldState = t.CurrentState
			t.CurrentState = toState
			if t.callback != nil {
				t.callback(oldState, toState)
			}
			break
		}
	}
}

// AddParameter 添加状态切换参数
func (t *StateMachine) AddParameter(parameter *Parameter) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.Parameters[parameter.Name]; ok {
		return fmt.Errorf("parameter name=%s already exists", parameter.Name)
	}

	t.Parameters[parameter.Name] = parameter

	return
}

// SetParameterValue 设置参数值并自动切换对应的状态
func (t *StateMachine) SetParameterValue(parameterName string, value string) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	parameter, ok := t.Parameters[parameterName]
	if !ok {
		return fmt.Errorf("parameter=%s not found", parameterName)
	}

	parameter.Value = value

	if transitions, ok := t.ParametersLink[parameter]; ok {
		t.transit(transitions)
	}

	return
}

func (t *StateMachine) RemoveParameter(parameterName string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	parameter := t.GetParameter(parameterName)
	delete(t.Parameters, parameterName)
	if parameter != nil {
		delete(t.ParametersLink, parameter)
	}
}

// GetMachine 取得状态机，参数形如 /App/Game/Match，如果不存在则会返回空指针
func (t *StateMachine) GetMachine(namepath string) (sm *StateMachine) {
	namepath = strings.TrimPrefix(namepath, "/")
	namepath = strings.TrimSuffix(namepath, "/")
	arr := strings.Split(namepath, "/")
	if len(arr) > 0 {
		if len(arr) == 1 {
			if t.Name == arr[0] {
				return t
			}
		} else {
			if t.Name == arr[0] {
				for _, ssm := range t.SubMachines {
					if ssm.Name == arr[1] {
						return ssm.GetMachine("/" + strings.Join(arr[1:], "/"))
					}
				}
			}
		}
	}
	return sm
}
