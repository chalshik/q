package main

import "time"

type Worker struct {
	status string
}

func (w *Worker) sendHeartbeat(heartbeatCh chan bool) {
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			heartbeatCh <- true
		}
		// stops here — monitor should detect death
	}()
}
func main() {
	worker := &Worker{}
	heartbeatCh := make(chan bool)
	worker.sendHeartbeat(heartbeatCh)

	for {
		select {
		case <-time.After(5 * time.Second):
			println("No heartbeat received from worker for 5 seconds")
		case <-heartbeatCh:
			println("Heartbeat received from worker")
		}
	}

}
