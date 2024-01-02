package fsm

import (
	"fmt"
	"slices"
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
	t.ValidTransition = make(map[string][]string)
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
	// ValidTransition 针对自动状态转换 Transitions 字段所做的约束
	ValidTransition map[State][]State
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

func (t *StateMachine) addState(states ...State) {
	for _, state := range states {
		if slices.Contains(t.States, state) {
			continue
		}

		t.States = append(t.States, state)
	}
}

// AddValidTransition 添加状态切换范围约束，即一个状态可以切换为哪些状态
func (t *StateMachine) AddValidTransition(fromState State, toStates []State) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.addValidTransition(fromState, toStates)
}

func (t *StateMachine) addValidTransition(fromState State, toStates []State) {
	if _, ok := t.ValidTransition[fromState]; !ok {
		t.ValidTransition[fromState] = make([]string, 1)
	}

	t.addState(fromState)

	for _, toState := range toStates {
		if slices.Contains(t.ValidTransition[fromState], toState) {
			continue
		}
		t.ValidTransition[fromState] = append(t.ValidTransition[fromState], toState)
		t.addState(toState)
	}
}

// checkTransitionValid 检查状态切换是否合法
func (t *StateMachine) checkTransitionValid(fromState State, toState State) bool {
	if _, ok := t.ValidTransition[fromState]; ok {
		if slices.Contains(t.ValidTransition[fromState], toState) {
			return true
		}
	}
	return false
}

// AddAutoTransition 添加自动状态切换
func (t *StateMachine) AddAutoTransition(trans *Transition, parameter *Parameter) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.checkTransitionValid(trans.From, trans.To) {
		return fmt.Errorf("transition=%s fromState=%s toState=%s was not registered in valid transition set", trans.Name, trans.From, trans.To)
	}

	if _, ok := t.Transitions[trans.Name]; ok {
		return fmt.Errorf("transition=%s already exists", trans.Name)
	}

	t.addValidTransition(trans.From, []State{trans.To})

	t.Transitions[trans.Name] = trans
	t.ParametersLink[parameter] = append(t.ParametersLink[parameter], trans)

	return
}

// RemoveAutoTransition 移除自动状态转换
func (t *StateMachine) RemoveAutoTransition(transitionName string) {
	transition, ok := t.Transitions[transitionName]
	if !ok {
		return
	}

	for parameter := range t.ParametersLink {
		for k, trans := range t.ParametersLink[parameter] {
			if trans == transition {
				slices.Delete(t.ParametersLink[parameter], k, k+1)
			}
		}
		if len(t.ParametersLink[parameter]) < 1 {
			delete(t.ParametersLink, parameter)
		}
	}

	delete(t.Transitions, transitionName)
}

// SetState 手动设置状态机状态，但会检查条件是否满足
func (t *StateMachine) SetState(toState State) (err error) {
	if t.checkTransitionValid(t.CurrentState, toState) {
		return fmt.Errorf("SetState fromState=%s toState=%s was not registered in valid transition set", t.CurrentState, toState)
	}

	// 查找条件约束
	var transSet []*Transition
	for _, trans := range t.Transitions {
		if trans.From == t.CurrentState && trans.To == toState {
			transSet = append(transSet, trans)
		}
	}

	if len(transSet) > 0 {
		t.autoTransit(transSet)
	} else {
		t.CurrentState = toState
	}

	return nil
}

// Transit 手动检查所有的状态切换是否需要进行一次状态切换
func (t *StateMachine) AutoTransit() {
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

	t.autoTransit(transitions)
}

func (t *StateMachine) autoTransit(transitions []*Transition) {
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

func (t *StateMachine) RemoveParameter(parameterName string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	parameter := t.GetParameter(parameterName)
	if parameter != nil {
		delete(t.Parameters, parameterName)
		delete(t.ParametersLink, parameter)
	}
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
		t.autoTransit(transitions)
	}

	return
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
