package fsm

import (
	"fmt"
)

type EventType string
type StateType string
type ReturnVal int
type ArgsType map[string]interface{}

// Callback newEvtCtxt将在下一个状态被处理
type Callback func(StateType, EventType, ArgsType) (newEvtCtxt *FsmEventContext, rv ReturnVal)
type LogFunc func(format string, a ...any)

type FsmEvent struct {
	HandleFunc Callback
	TransTable map[ReturnVal]StateType
}

type FsmState struct {
	DefaultHandle *FsmEvent
	EventTable    map[EventType]*FsmEvent
}

type Fsm struct {
	fsmId      string
	initState  StateType
	stateTable map[StateType]*FsmState
	logFunc    LogFunc
}

type FsmInst struct {
	fsmInstId string
	currState StateType
	sm        *Fsm
}

type FsmEventContext struct {
	Event EventType
	Args  ArgsType
}

func NewFsmEventContext(event EventType, args ArgsType) *FsmEventContext {
	evtCtxt := &FsmEventContext{
		Event: event,
		Args:  args,
	}
	return evtCtxt
}

func NewFsm(fsmId string, initState StateType, stateTable map[StateType]*FsmState) *Fsm {
	fsm := &Fsm{
		fsmId:      fsmId,
		initState:  initState,
		stateTable: stateTable,
		logFunc:    nil,
	}
	return fsm
}

func NewFsmInst(sm *Fsm, fsmInstId string) *FsmInst {
	fsmInst := &FsmInst{
		fsmInstId: fsmInstId,
		currState: sm.initState,
		sm:        sm,
	}
	return fsmInst
}

func (sm *Fsm) SetLogFunc(logFunc LogFunc) {
	sm.logFunc = logFunc
}

func (fsmInst *FsmInst) HandleEventContext(eventContext *FsmEventContext) error {
	return fsmInst.HandleEvent(eventContext.Event, eventContext.Args)
}

func (fsmInst *FsmInst) HandleEvent(event EventType, args ArgsType) error {
	eventTemp := event
	argsTemp := args
	for {
		newEvtCtxt, err := fsmInst.HandleEventImpl(eventTemp, argsTemp)
		if err != nil {
			return err
		}
		if newEvtCtxt == nil {
			break
		}
		eventTemp = newEvtCtxt.Event
		argsTemp = newEvtCtxt.Args
	}
	return nil
}

func (fsmInst *FsmInst) HandleEventImpl(event EventType, args ArgsType) (*FsmEventContext, error) {
	sm := fsmInst.sm
	currState := fsmInst.currState

	if sm.stateTable == nil {
		return nil, fmt.Errorf("FSM[%s] INST[%s]: state table is nil", sm.fsmId, fsmInst.fsmInstId)
	}

	smState, ok := sm.stateTable[currState]
	if !ok {
		return nil, fmt.Errorf("FSM[%s] INST[%s]: Can not find state[%s]", sm.fsmId, fsmInst.fsmInstId, currState)
	}
	if smState == nil {
		return nil, fmt.Errorf("FSM[%s] INST[%s]: state is nil", sm.fsmId, fsmInst.fsmInstId)
	}

	smEvent, ok := smState.EventTable[event]
	if !ok {
		smEvent = smState.DefaultHandle
	}
	if smEvent == nil {
		return nil, fmt.Errorf("FSM[%s] INST[%s] STATE[%s]: event[%s] handle and default handle is nil", sm.fsmId, fsmInst.fsmInstId, currState, event)
	}
	if smEvent.HandleFunc == nil {
		return nil, fmt.Errorf("FSM[%s] INST[%s] STATE[%s]: handle func is nil", sm.fsmId, fsmInst.fsmInstId, currState)
	}
	newEvtCtxt, rv := smEvent.HandleFunc(currState, event, args)

	newState, ok := smEvent.TransTable[rv]
	if !ok {
		return nil, fmt.Errorf("FSM[%s] INST[%s] STATE[%s]: Can not find new state by callback return value %d",
			sm.fsmId, fsmInst.fsmInstId, currState, rv)
	}
	fsmInst.currState = newState

	if sm.logFunc != nil {
		sm.logFunc("FSM[%s] INST[%s]: %s[%s] --> %s\n", sm.fsmId, fsmInst.fsmInstId, currState, event, newState)
	}

	return newEvtCtxt, nil
}

func (fsmInst *FsmInst) CurrState() StateType {
	return fsmInst.currState
}

func (fsmInst *FsmInst) ResetState() {
	fsmInst.currState = fsmInst.sm.initState
}
