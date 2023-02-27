package jobs

import (
	"context"
	"delta/core"
)

type JobExecutable func() error

// IProcessor is an interface that has a Run method that returns an error.
// @property {error} Run - This is the main function of the processor. It will be called by the processor manager.
type IProcessor interface {
	Run() error
}

// `Processor` is a struct that contains a `context.Context` and a `*core.DeltaNode`.
// @property Context - The context of the processor.
// @property LightNode - This is the node that is being processed.
type Processor struct {
	Context   context.Context
	LightNode *core.DeltaNode
}
