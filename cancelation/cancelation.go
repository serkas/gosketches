package main

import (
	"log"
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
	var stopChannel = make(chan struct{}) // Channel to notify all workers
	ctl := &Ctl{stopCh: stopChannel}

	nTasks := 3
	log.Printf("Start %d tasks\n", nTasks)
	// Spawn workers
	for i := 1; i <= nTasks; i++ {
		ctl.wg.Add(1)
		go work(i, ctl)
	}

	// Stopping. Capture OS signals and route them to our channel
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	log.Printf("Captured %v. Exiting...\n", sig)
	close(stopChannel) // notify workers

	ctl.wg.Wait()
	os.Exit(0)
}

func work(taskId int, ctl *Ctl) {
	counter := 0
	defer ctl.wg.Done()
	for {
		select {
		case <-ctl.stopCh:
			log.Printf("Task %d ->  Worker stopped.", taskId)
			return
		default:
			log.Printf( "Task %d.%d -> Started...\n", taskId, counter)
			time.Sleep(5 * time.Second) // Task logic to do
			log.Printf("Task %d.%d -> Done\n", taskId, counter)
			counter++
		}
	}
}
