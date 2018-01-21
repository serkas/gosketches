package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)



type Ctl struct {
	wg     *sync.WaitGroup
	stopCh chan struct{}
}

func main() {
	fmt.Println("Start.")
	var n sync.WaitGroup
	var stopChannel = make(chan struct{})

	ctl := Ctl{wg: &n, stopCh: stopChannel}

	for i := 1; i <= 3; i++ {
		ctl.wg.Add(1)
		go work(fmt.Sprintf("Task %d", i), ctl)
	}


	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-signalChan
		fmt.Printf("\n\nCaptured %v. Exiting...\n", sig)
		close(stopChannel)
		ctl.wg.Wait()
		os.Exit(0)
	}
}

func work(taskName string, ctl Ctl) {
	defer ctl.wg.Done()
	for {
		select {
		case <-ctl.stopCh:
			fmt.Println(taskName + " ->  Worker stopped.")
			return
		default:
			fmt.Println(taskName + " -> I'm doing ... ")
			time.Sleep(5 * time.Second)
			fmt.Println(taskName + " -> Done")
		}

	}
}
