package jobs

import (
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
)

type LogEventProcessor struct {
	LightNode *core.DeltaNode
	LogEvent  model.LogEvent
}

func NewLogEvent(ln *core.DeltaNode, logEvent model.LogEvent) IProcessor {
	return &LogEventProcessor{
		LightNode: ln,
		LogEvent:  logEvent,
	}
}

func (l LogEventProcessor) Run() error {
	// save log event
	l.LightNode.DB.Create(&l.LogEvent)
	return nil
}
