// Package core Create a dispatcher, add jobs to it, and then start it with a number of workers
package core

import (
	"fmt"
)

type JobExecutable func() error

// IProcessor is an interface that has a Run method that returns an error.
// @property {error} Run - This is the main function of the processor. It will be called by the processor manager.
type IProcessor interface {
	Run() error
}

// A Job is a struct that has an ID and a Processor.
// @property {int} ID - The ID of the job.
// @property {IProcessor} Processor - This is the interface that the job will use to process itself.
type Job struct {
	ID        int
	Processor IProcessor
}

// A Worker is a struct that has an ID, a jobs channel, a dispatchStatus channel, and a Quit channel.
// @property {int} ID - The ID of the worker.
// @property jobs - This is the channel that the worker will use to receive jobs from the dispatcher.
// @property dispatchStatus - This is a channel that will be used to send a message to the worker when it's done processing
// a job.
// @property Quit - This is a channel that will be used to tell the worker to stop working.
type Worker struct {
	ID             int
	jobs           chan *Job
	dispatchStatus chan *DispatchStatus
	Quit           chan bool
}

// CreateNewWorker Create a new worker, add it to the worker queue, and return it
func CreateNewWorker(id int, workerQueue chan *Worker, jobQueue chan *Job, dStatus chan *DispatchStatus) *Worker {
	w := &Worker{
		ID:             id,
		jobs:           jobQueue,
		dispatchStatus: dStatus,
	}

	go func() { workerQueue <- w }()
	return w
}

// Start It's a goroutine that is waiting for a job to be added to the worker's job channel.
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

// DispatchStatus is a struct with three fields, Type, ID, and Status.
//
// The first field, Type, is a string. The second field, ID, is an int. The third field, Status, is a string.
//
// The DispatchStatus type is a struct type. A struct type is a collection of fields.
//
// The DispatchStatus type has three fields. Each field has a name and a type.
//
// The DispatchStatus type has three fields: Type, ID, and Status.
//
// The Type field has a name and a type. The name is Type
// @property {string} Type - The type of the dispatch. This can be either "order" or "quote".
// @property {int} ID - The ID of the dispatch
// @property {string} Status - The status of the dispatch. This can be one of the following:
type DispatchStatus struct {
	Type   string
	ID     int
	Status string
}

// A Dispatcher is a struct that contains a jobCounter, jobQueue, dispatchStatus, workQueue, and workerQueue.
// @property {int} jobCounter - an internal counter for the number of jobs submitted
// @property jobQueue - This is the channel that jobs are submitted to.
// @property dispatchStatus - This is a channel that will be used to report the status of jobs and workers.
// @property workQueue - This is the channel that the workers will use to send back the results of their work.
// @property workerQueue - a channel of workers
type Dispatcher struct {
	jobCounter     int                  // internal counter for number of jobs
	jobQueue       chan *Job            // channel of jobs submitted
	dispatchStatus chan *DispatchStatus // channel for job/worker status reports
	workQueue      chan *Job            // channel of work dispatched
	workerQueue    chan *Worker         // channel of workers
}

// CreateNewDispatcher Create a new dispatcher, and return a pointer to it
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

// Start Creating a number of workers and then waiting for jobs to be added to the jobQueue.
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

// AddJob Adding a job to the jobQueue.
// It's adding a job to the jobQueue.
func (d *Dispatcher) AddJob(je IProcessor) {
	j := &Job{ID: d.jobCounter, Processor: je}
	go func() { d.jobQueue <- j }()
	d.jobCounter++
	fmt.Printf("Number Of Jobs: %d\n", d.jobCounter)
}

// AddJobAndDispatch Adding a job to the jobQueue, and then starting the dispatcher with a number of workers.
// It's adding a job to the jobQueue, and then starting the dispatcher with a number of workers.
func (d *Dispatcher) AddJobAndDispatch(je IProcessor, numWorkers int) {
	j := &Job{ID: d.jobCounter, Processor: je}
	go func() { d.jobQueue <- j }()
	d.jobCounter++
	fmt.Printf("Number Of Jobs: %d\n", d.jobCounter)
	d.Start(numWorkers)
}

// Finished It's a method that returns true if the jobCounter is less than 1.
// It's a method that returns true if the jobCounter is less than 1.
func (d *Dispatcher) Finished() bool {
	if d.jobCounter < 1 {
		return true
	} else {
		return false
	}
}
