package main

import "fmt"

////states////

const (
	stateMainMenu uint = iota
	stateNewGame
	stateRunningGame
	statePausedGame
	stateExit
	totalStates
)


////subStates////

//subStates: stateMainMenu
const (
	stateMainMenuSelectorIndex uint = iota
	stateMainMenuSubStates
)
const (
	stateMainMenuSelectorInit int = 0
)

//subStates: stateNewGame
const (
	stateNewGameSizeSelectorIndex uint = iota
	stateNewGameSubStates
)
const (
	stateNewGameSizeSelectorInit int = 0
)

//subStates: stateRunningGame
const (
	stateRunningGameSubStates uint = iota
)

//subStates: statePausedGame
const (
	statePausedGameSelectorIndex uint = iota
	statePausedGameSubStates
)
const (
	statePausedGameSelectorInit int = 0
)

//subStates: stateExit
const (
	stateExitSubStates uint = iota
)

//A shorthand way of accessing the number of substates per state.
var (
	totalSubStates [totalStates]uint = [totalStates]uint{
		stateMainMenuSubStates,
		stateNewGameSubStates,
		stateRunningGameSubStates,
		statePausedGameSubStates,
		stateExitSubStates,
	}
)

//A shorthand was of accessing the initial values of substates for each state.
var (
	subStateInitialValues [totalStates][]int = [totalStates][]int{
		[]int{stateMainMenuSelectorInit},
		[]int{stateNewGameSizeSelectorInit},
		[]int{},
		[]int{statePausedGameSelectorInit},
		[]int{},
	}
)


////Data structures////

const (
	stateType uint8 = iota
	subStateType
	errorType
)

type stateRequest struct {
	reqType uint8
	subStateIndex uint
}

func initStateReq() stateRequest {
	sr := stateRequest{
		reqType: stateType,
		subStateIndex: 0,
	}
	return sr
}

func initSubStateReq(i uint) stateRequest {
	ssr := stateRequest{
		reqType: subStateType,
		subStateIndex: i,
	}
	return ssr
}

type stateDescriptor struct {
	descType uint8
	state uint
	subState int
	subStateIndex uint
}

func initStateDesc(s uint) stateDescriptor {
	sd := stateDescriptor{
		descType: stateType,
		state: s,
		subState: 0,
		subStateIndex: 0,
	}
	return sd
}

func initSubStateDesc(s uint, ss int, i uint) stateDescriptor {
	ssd := stateDescriptor{
		descType: subStateType,
		state: s,
		subState: ss,
		subStateIndex: i,
	}
	return ssd
}

func initErrorDesc() stateDescriptor {
	err := stateDescriptor{
		descType: errorType,
		state: 0,
		subState: 0,
		subStateIndex: 0,
	}
	return err
}


////Full Request Functions////

func getState(rqst chan<- stateRequest, rcv <-chan stateDescriptor) uint {
	rqst <- initStateReq()
	return (<- rcv).state
}

func getSubState(rqst chan<- stateRequest, rcv <-chan stateDescriptor, state uint, subStateIndex uint) int {
	if subStateIndex >= totalSubStates[state] {
		panic(fmt.Sprintf("State %d does not have a subState %d.", state, subStateIndex))
	}
	rqst <- initSubStateReq(subStateIndex)
	desc := <- rcv
	if desc.descType == errorType {
		panic(fmt.Sprintf("State has changed from %d to %d.", state, getState(rqst, rcv)))
	}else if desc.descType == subStateType && desc.state != state {
		panic(fmt.Sprintf("State has changed from %d to %d.", state, desc.state))
	}
	return desc.subState
}