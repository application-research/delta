package jobs

import (
	"delta/core"
)

type RepairProcessor struct {
	LightNode *core.DeltaNode
}

// `NewRepairProcessor` creates a new `RepairProcessor` struct and returns it
func NewRepairProcessor(ln *core.DeltaNode) IProcessor {
	return &RepairProcessor{
		LightNode: ln,
	}
}

// Run DB heavy process. We need to check the status of the content and requeue the job if needed.
func (i RepairProcessor) Run() error {
	return nil
}
