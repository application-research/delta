package jobs

import (
	"delta/core"
	"delta/core/model"
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
