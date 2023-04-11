package core

import (
	"errors"
	"testing"
)

type TestProcessor struct {
	ID   int
	Name string
}

func (tp *TestProcessor) Run() error {
	if tp.ID == 0 {
		return errors.New("Invalid Test Processor ID")
	}
	return nil
}

func TestDispatcher(t *testing.T) {
	// create a new dispatcher
	dispatcher := CreateNewDispatcher()

	// add a job to the job queue
	testProcessor := &TestProcessor{ID: 1, Name: "Test Processor 1"}
	dispatcher.AddJob(testProcessor)

	// start the dispatcher with one worker
	dispatcher.Start(1)

	// wait until the dispatcher has finished processing all jobs
	for !dispatcher.Finished() {
	}

	// check if the job was processed successfully
	if testProcessor.ID != 1 {
		t.Errorf("Test Processor did not run successfully. Expected ID: 1, Actual ID: %d", testProcessor.ID)
	}
}
