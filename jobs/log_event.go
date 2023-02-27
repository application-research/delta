package jobs

import (
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
)

// `LogEventProcessor` is a struct that contains a `LightNode` and a `LogEvent`.
// @property LightNode - The node that the event is being processed for.
// @property LogEvent - This is the event that we want to process.
type LogEventProcessor struct {
	LightNode *core.DeltaNode
	LogEvent  model.LogEvent
}

// > This function creates a new LogEventProcessor object and returns it
func NewLogEvent(ln *core.DeltaNode, logEvent model.LogEvent) IProcessor {
	return &LogEventProcessor{
		LightNode: ln,
		LogEvent:  logEvent,
	}
}

// Saving the log event to the database.
func (l LogEventProcessor) Run() error {
	// save log event
	l.LightNode.DB.Create(&l.LogEvent)
	return nil
}
