package main

import "github.com/nsf/termbox-go"

func runInputParser(envRcv <-chan *environment, envRqst chan<- bool, envMdfy chan<- bool, stRcv <-chan uint8, stRqst chan<- bool, stMdfy chan<- uint8) {
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Ch == 'e' {
				stMdfy <- stateExit
			}else if event.Ch == 'n' {
				stRqst <- true
				newState := <- stRcv
				stMdfy <- ((newState + 1) % stateExit)
			}
		}
	}
}