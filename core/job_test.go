package core

import (
	"errors"
	"sync"
	"testing"
)

// MockProcessor is a mock implementation of the IProcessor interface
type MockProcessor struct {
	err error
}

func (m *MockProcessor) Run() error {
	return m.err
}

func TestDispatcher(t *testing.T) {
	// Create a new dispatcher
	d := CreateNewDispatcher()

	// Add some jobs to the dispatcher
	numJobs := 5
	var wg sync.WaitGroup
	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			d.AddJob(&MockProcessor{})
		}(i)
	}

	// Start the dispatcher with two workers
	d.Start(2)

	// Wait for all jobs to finish
	wg.Wait()

	// Check that there were no errors processing the jobs
	if !d.Finished() {
		t.Errorf("Dispatcher did not finish all jobs")
	}

	// Check that the dispatchStatus channel received the expected messages
	var workerQuitCount, jobCount int
	for {
		select {
		case ds := <-d.dispatchStatus:
			if ds.Type == "worker" && ds.Status == "quit" {
				workerQuitCount++
			}
			if ds.Type == "job" && ds.Status == "done" {
				jobCount++
			}
			if workerQuitCount == 2 && jobCount == numJobs {
				// Expected number of messages received
				return
			}
		default:
			// No more messages in dispatchStatus channel
			t.Errorf("Did not receive expected dispatch status messages")
			return
		}
	}
}

func TestJobWithError(t *testing.T) {
	// Create a new dispatcher
	d := CreateNewDispatcher()

	// Add a job with an error to the dispatcher
	err := errors.New("job error")
	d.AddJob(&MockProcessor{err})

	// Start the dispatcher with one worker
	d.Start(1)

	// Check that the dispatcher finished with an error
	if !d.Finished() {
		t.Errorf("Dispatcher did not finish all jobs")
	}

	// Check that the dispatchStatus channel received the expected messages
	var workerQuitCount, jobCount int
	for {
		select {
		case ds := <-d.dispatchStatus:
			if ds.Type == "worker" && ds.Status == "quit" {
				workerQuitCount++
			}
			if ds.Type == "job" && ds.Status == "error" {
				jobCount++
			}
			if workerQuitCount == 1 && jobCount == 1 {
				// Expected number of messages received
				return
			}
		default:
			// No more messages in dispatchStatus channel
			t.Errorf("Did not receive expected dispatch status messages")
			return
		}
	}
}
