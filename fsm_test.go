package fsm_test

import (
	"fmt"
	"github.com/todayVenturn/fsm"
	"testing"
)

var testTime = 0

func fsm_test_default_handler1(state fsm.StateType, event fsm.EventType, args fsm.ArgsType) (newEvtCtxt *fsm.FsmEventContext, rv fsm.ReturnVal) {
	return nil, 0
}
func fsm_test_idle_state_handler1(state fsm.StateType, event fsm.EventType, args fsm.ArgsType) (newEvtCtxt *fsm.FsmEventContext, rv fsm.ReturnVal) {
	return nil, 0
}

func fsm_test_wait_ack_state_handler1(state fsm.StateType, event fsm.EventType, args fsm.ArgsType) (newEvtCtxt *fsm.FsmEventContext, rv fsm.ReturnVal) {
	if testTime%2 == 0 {
		rv = 0
	} else {
		rv = -1
	}
	testTime += 1
	return nil, rv
}

func fsm_test_wait_close_state_handler1(state fsm.StateType, event fsm.EventType, args fsm.ArgsType) (newEvtCtxt *fsm.FsmEventContext, rv fsm.ReturnVal) {
	return nil, 0
}

func fsmLogFunc(format string, a ...any) {
	fmt.Printf(format, a...)
}

func TestFSM(t *testing.T) {
	sm := fsm.NewFsm("TestFsm", "STATE_IDLE", map[fsm.StateType]*fsm.FsmState{
		"STATE_IDLE": {
			DefaultHandle: nil,
			EventTable: map[fsm.EventType]*fsm.FsmEvent{
				"EVENT_CONNECT": {
					fsm_test_idle_state_handler1, map[fsm.ReturnVal]fsm.StateType{
						0:  "STATE_WAIT_ACK",
						-1: "STATE_IDLE"}}}},
		"STATE_WAIT_ACK": {
			DefaultHandle: &fsm.FsmEvent{
				fsm_test_default_handler1, map[fsm.ReturnVal]fsm.StateType{
					0: "STATE_IDLE"}},
			EventTable: map[fsm.EventType]*fsm.FsmEvent{
				"EVENT_RCV_ACK": {
					fsm_test_wait_ack_state_handler1, map[fsm.ReturnVal]fsm.StateType{
						0:  "STATE_WAIT_CLOSE",
						-1: "STATE_IDLE"}}}},
		"STATE_WAIT_CLOSE": {
			DefaultHandle: nil,
			EventTable: map[fsm.EventType]*fsm.FsmEvent{
				"EVENT_CLOSE": {
					fsm_test_wait_close_state_handler1, map[fsm.ReturnVal]fsm.StateType{
						0: "STATE_IDLE"}}}},
	})
	fmt.Println("sm: ", sm)
	sm.SetLogFunc(fsmLogFunc)

	smInst := fsm.NewFsmInst(sm, "TestFsmInst")
	fmt.Println("sm: ", smInst)

	err := smInst.HandleEvent("EVENT_CONNECT", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if smInst.CurrState() != "STATE_WAIT_ACK" {
		t.Fatalf("state[%v] error", smInst.CurrState())
	}

	evtCtxt := &fsm.FsmEventContext{"EVENT_CONNECT", nil}
	err = smInst.HandleEventContext(evtCtxt)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if smInst.CurrState() != "STATE_IDLE" {
		t.Fatalf("state[%v] error", smInst.CurrState())
	}

	err = smInst.HandleEvent("EVENT_CONNECT", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if smInst.CurrState() != "STATE_WAIT_ACK" {
		t.Fatalf("state[%v] error", smInst.CurrState())
	}

	err = smInst.HandleEvent("EVENT_RCV_ACK", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if smInst.CurrState() != "STATE_WAIT_CLOSE" {
		t.Fatalf("state[%v] error", smInst.CurrState())
	}

	err = smInst.HandleEvent("EVENT_CLOSE", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if smInst.CurrState() != "STATE_IDLE" {
		t.Fatalf("state[%v] error", smInst.CurrState())
	}

	err = smInst.HandleEvent("EVENT_CONNECT", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if smInst.CurrState() != "STATE_WAIT_ACK" {
		t.Fatalf("state[%v] error", smInst.CurrState())
	}

	err = smInst.HandleEvent("EVENT_RCV_ACK", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if smInst.CurrState() != "STATE_IDLE" {
		t.Fatalf("state[%v] error", smInst.CurrState())
	}
}
