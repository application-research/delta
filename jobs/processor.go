package jobs

import (
	"context"
	"fc-deal-making-service/core"
)

type JobExecutable func() error
type IProcessor interface {
	Run() error
}

type Processor struct {
	Context   context.Context
	LightNode *core.LightNode
}

type ContentProcessor struct {
	Context   context.Context
	LightNode *core.LightNode
	Content   core.Content
}
