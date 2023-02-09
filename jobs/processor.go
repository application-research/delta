package jobs

import (
	"context"
	"delta/core"
)

type JobExecutable func() error
type IProcessor interface {
	Run() error
}

type Processor struct {
	Context   context.Context
	LightNode *core.LightNode
}
