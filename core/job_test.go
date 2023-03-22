package core

import (
	"reflect"
	"testing"
)

func TestCreateNewDispatcher(t *testing.T) {
	tests := []struct {
		name string
		want *Dispatcher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateNewDispatcher(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateNewDispatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateNewWorker(t *testing.T) {
	type args struct {
		id          int
		workerQueue chan *Worker
		jobQueue    chan *Job
		dStatus     chan *DispatchStatus
	}
	tests := []struct {
		name string
		args args
		want *Worker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateNewWorker(tt.args.id, tt.args.workerQueue, tt.args.jobQueue, tt.args.dStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateNewWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDispatcher_AddJob(t *testing.T) {
	type fields struct {
		jobCounter     int
		jobQueue       chan *Job
		dispatchStatus chan *DispatchStatus
		workQueue      chan *Job
		workerQueue    chan *Worker
	}
	type args struct {
		je IProcessor
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{
				jobCounter:     tt.fields.jobCounter,
				jobQueue:       tt.fields.jobQueue,
				dispatchStatus: tt.fields.dispatchStatus,
				workQueue:      tt.fields.workQueue,
				workerQueue:    tt.fields.workerQueue,
			}
			d.AddJob(tt.args.je)
		})
	}
}

func TestDispatcher_AddJobAndDispatch(t *testing.T) {
	type fields struct {
		jobCounter     int
		jobQueue       chan *Job
		dispatchStatus chan *DispatchStatus
		workQueue      chan *Job
		workerQueue    chan *Worker
	}
	type args struct {
		je         IProcessor
		numWorkers int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{
				jobCounter:     tt.fields.jobCounter,
				jobQueue:       tt.fields.jobQueue,
				dispatchStatus: tt.fields.dispatchStatus,
				workQueue:      tt.fields.workQueue,
				workerQueue:    tt.fields.workerQueue,
			}
			d.AddJobAndDispatch(tt.args.je, tt.args.numWorkers)
		})
	}
}

func TestDispatcher_Finished(t *testing.T) {
	type fields struct {
		jobCounter     int
		jobQueue       chan *Job
		dispatchStatus chan *DispatchStatus
		workQueue      chan *Job
		workerQueue    chan *Worker
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{
				jobCounter:     tt.fields.jobCounter,
				jobQueue:       tt.fields.jobQueue,
				dispatchStatus: tt.fields.dispatchStatus,
				workQueue:      tt.fields.workQueue,
				workerQueue:    tt.fields.workerQueue,
			}
			if got := d.Finished(); got != tt.want {
				t.Errorf("Finished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDispatcher_Start(t *testing.T) {
	type fields struct {
		jobCounter     int
		jobQueue       chan *Job
		dispatchStatus chan *DispatchStatus
		workQueue      chan *Job
		workerQueue    chan *Worker
	}
	type args struct {
		numWorkers int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{
				jobCounter:     tt.fields.jobCounter,
				jobQueue:       tt.fields.jobQueue,
				dispatchStatus: tt.fields.dispatchStatus,
				workQueue:      tt.fields.workQueue,
				workerQueue:    tt.fields.workerQueue,
			}
			d.Start(tt.args.numWorkers)
		})
	}
}

func TestWorker_Start(t *testing.T) {
	type fields struct {
		ID             int
		jobs           chan *Job
		dispatchStatus chan *DispatchStatus
		Quit           chan bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				ID:             tt.fields.ID,
				jobs:           tt.fields.jobs,
				dispatchStatus: tt.fields.dispatchStatus,
				Quit:           tt.fields.Quit,
			}
			w.Start()
		})
	}
}
