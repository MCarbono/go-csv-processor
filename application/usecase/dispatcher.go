package usecase

type Dispathcer interface {
	Stop()
	LaunchWorker(w WorkerLauncher)
	Launch(input []string)
}

type dispatcher struct {
	inCh chan []string
}

func (d *dispatcher) Stop() {
	close(d.inCh)
}

func (d *dispatcher) LaunchWorker(w WorkerLauncher) {
	w.LaunchWorker(d.inCh)
}

func (d *dispatcher) Launch(input []string) {
	d.inCh <- input
}

func NewDispatcher(bufferSize int) Dispathcer {
	return &dispatcher{
		inCh: make(chan []string, bufferSize),
	}
}
