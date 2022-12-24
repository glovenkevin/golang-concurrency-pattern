package workers_pool

import "time"

type DispatcherInterface interface {
	LaunchWorker(w WorkerLauncherInterface)
	Stop()
	MakeRequest(Request)
}

type Dispatcher struct {
	In chan Request
}

func NewDispatcher(b int) Dispatcher {
	return Dispatcher{
		In: make(chan Request, b),
	}
}

func (d *Dispatcher) LaunchWorker(w WorkerLauncherInterface) {
	w.LaunchWorker(d.In)
}

func (d *Dispatcher) Stop() {
	close(d.In)
}

func (d *Dispatcher) MakeRequest(r Request) {
	select {
	case d.In <- r:
	case <-time.After(time.Second * 5):
		return
	}
}
