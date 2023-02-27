// Create a dispatcher, add jobs to it, and then start it with a number of workers
package core

import (
	"fmt"
)

type JobExecutable func() error
type IProcessor interface {
	Run() error
}

type Job struct {
	ID        int
	Processor IProcessor
}

type Worker struct {
	ID             int
	jobs           chan *Job
	dispatchStatus chan *DispatchStatus
	Quit           chan bool
}

// Creating a new worker and adding it to the workerQueue.
func CreateNewWorker(id int, workerQueue chan *Worker, jobQueue chan *Job, dStatus chan *DispatchStatus) *Worker {
	w := &Worker{
		ID:             id,
		jobs:           jobQueue,
		dispatchStatus: dStatus,
	}

	go func() { workerQueue <- w }()
	return w
}

// It's a goroutine that is waiting for a job to be added to the worker's job channel.
// When a job is added, it is executed and then the worker sends a quit message to the dispatcher.
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.jobs:
				fmt.Printf("Worker[%d] executing job[%d].\n", w.ID, job.ID)
				job.Processor.Run()
				w.dispatchStatus <- &DispatchStatus{Type: "worker", ID: w.ID, Status: "quit"}
				w.Quit <- true
			case <-w.Quit:
				return
			}
		}
	}()
}

type DispatchStatus struct {
	Type   string
	ID     int
	Status string
}

type Dispatcher struct {
	jobCounter     int                  // internal counter for number of jobs
	jobQueue       chan *Job            // channel of jobs submitted
	dispatchStatus chan *DispatchStatus // channel for job/worker status reports
	workQueue      chan *Job            // channel of work dispatched
	workerQueue    chan *Worker         // channel of workers
}

// Create a new dispatcher, and return a pointer to it
func CreateNewDispatcher() *Dispatcher {
	d := &Dispatcher{
		jobCounter:     0,
		jobQueue:       make(chan *Job),
		dispatchStatus: make(chan *DispatchStatus),
		workQueue:      make(chan *Job),
		workerQueue:    make(chan *Worker),
	}
	return d
}

// Creating a number of workers and then waiting for jobs to be added to the jobQueue.
// When a job is added, it is sent to the workQueue, which is the channel that the workers are listening to.
func (d *Dispatcher) Start(numWorkers int) {
	// Create numWorkers:
	for i := 0; i < numWorkers; i++ {
		worker := CreateNewWorker(i, d.workerQueue, d.workQueue, d.dispatchStatus)
		worker.Start()
	}

	// wait for work to be added then pass it off.
	go func() {
		for {
			select {
			case job := <-d.jobQueue:
				fmt.Printf("Got a job in the queue to dispatch: %d\n", job.ID)
				// Sending it off;
				d.workQueue <- job
			case ds := <-d.dispatchStatus:
				fmt.Printf("Got a dispatch status: Type[%s] - ID[%d] - Status[%s]", ds.Type, ds.ID, ds.Status)
				if ds.Type == "worker" {
					if ds.Status == "quit" {
						d.jobCounter--
					}
				}
			}
		}
	}()
}

// Adding a job to the jobQueue.
func (d *Dispatcher) AddJob(je IProcessor) {
	j := &Job{ID: d.jobCounter, Processor: je}
	go func() { d.jobQueue <- j }()
	d.jobCounter++
	fmt.Printf("Number Of Jobs: %d\n", d.jobCounter)
}

// Adding a job to the jobQueue, and then starting the dispatcher with a number of workers.
func (d *Dispatcher) AddJobAndDispatch(je IProcessor, numWorkers int) {
	j := &Job{ID: d.jobCounter, Processor: je}
	go func() { d.jobQueue <- j }()
	d.jobCounter++
	fmt.Printf("Number Of Jobs: %d\n", d.jobCounter)
	d.Start(numWorkers)
}

// It's a method that returns true if the jobCounter is less than 1.
func (d *Dispatcher) Finished() bool {
	if d.jobCounter < 1 {
		return true
	} else {
		return false
	}
}
