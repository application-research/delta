package jobs

import "fc-deal-making-service/core"

type Processor struct {
	ProcessorInterface
	LightNode *core.LightNode
}

type ProcessorInterface interface {
	PreProcess()
	PostProcess()
	Run()
	Verify()
}
