package kocto

type Message any

type Pipeline struct {
	In  chan<- Message
	Out <-chan Message

	logger  Logger
	workers []*pipeWorker
}

func NewPipeline() *Pipeline {
	//return &ConcurrentPipeline{workers: make([]stageWorker, 0)}
	return &Pipeline{workers: make([]*pipeWorker, 0)}
}

func (p *Pipeline) AddStage(stage Stage, opt *StageOptions) {
	if opt == nil {
		opt = &StageOptions{Concurrency: 10} // should be > 1
	}

	in := make(chan Message, 10)
	out := make(chan Message, 10)

	if len(p.workers) == 0 {
		p.In = in
	} else {
		// if it's not the first stage the input channel must be the output of the previous stage
		in = p.last().Out
	}

	p.Out = out

	worker := newPipeWorker(p.logger, opt.Concurrency, in, out, stage)
	p.workers = append(p.workers, worker)
}

func (p *Pipeline) Start() error {
	for _, w := range p.workers {
		w.Start()
	}

	return nil
}

func (p *Pipeline) last() *pipeWorker {
	lastIdx := len(p.workers) - 1
	return p.workers[lastIdx]
}

func (p *Pipeline) Stop() error {
	for _, w := range p.workers {
		w.Stop()
	}

	close(p.last().Out)
	return nil
}
