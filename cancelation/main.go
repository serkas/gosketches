package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Control
type Ctl struct {
	wg     sync.WaitGroup
	stopCh chan struct{}
}

func main() {
	nTasks := 3
	fmt.Printf("Start %d tasks\n", nTasks)

	var stopChannel = make(chan struct{}) // Channel to notify all workers
	ctl := &Ctl{stopCh: stopChannel}

	// Spawn workers
	for i := 1; i <= nTasks; i++ {
		ctl.wg.Add(1)
		taskName := fmt.Sprintf("Task %d", i)
		go work(taskName, ctl)
	}

	// Stopping
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	fmt.Printf("\nCaptured %v. Exiting...\n", sig)
	close(stopChannel) // notify workers

	ctl.wg.Wait()
	os.Exit(0)
}

func work(taskName string, ctl *Ctl) {
	defer ctl.wg.Done()
	for {
		select {
		case <-ctl.stopCh:
			fmt.Println(taskName + " ->  Worker stopped.")
			return
		default:
			fmt.Println(taskName + " -> I'm doing ... ")
			time.Sleep(5 * time.Second) // Task logic to do
			fmt.Println(taskName + " -> Done")
		}
	}
}
