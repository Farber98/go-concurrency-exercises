//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var sigintCallings int

func main() {
	// Create a process
	proc := MockProcess{}

	// Run the process (non blocking)
	go proc.Run()

	// Create channel that listen SIGINT (^C) signals.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	defer close(c)

	for {
		select {
		case <-c: // First time it receives a signal calls proc.Stop(). Second time calls os.Exit(1).
			sigintCallings++
			if sigintCallings > 1 {
				fmt.Println("\nNon gracefully shutdown.")
				os.Exit(1)
			}
			go proc.Stop()
		}
	}
}
