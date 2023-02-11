package jobs

import "delta/core"

type LogEventProcessor struct {
	LightNode *core.DeltaNode
	LogEvent  core.LogEvent
}

func NewLogEvent(ln *core.DeltaNode, logEvent core.LogEvent) IProcessor {
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
