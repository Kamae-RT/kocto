package kocto

type Message any

type Pipeline interface {
	Start() error
	Stop() error

	Input() chan<- Message
	Output() <-chan Message

	AddStage(Stage, *StageOptions)
}

func NewPipeline() Pipeline {
	return &ConcurrentPipeline{workers: make([]stageWorker, 0)}
}

type ConcurrentPipeline struct {
	workers []stageWorker
}

func (p *ConcurrentPipeline) AddStage(stage Stage, opt *StageOptions) {
	if opt == nil {
		opt = &StageOptions{Concurrency: 10} // should be > 1
	}

	in := make(chan Message, 10)
	out := make(chan Message, 10)

	// if it's not the first stage the input channel must be the output of the previous stage
	for _, w := range p.workers {
		in = w.Output()
	}

	worker := newStageWorker(opt.Concurrency, in, out, stage)
	p.workers = append(p.workers, worker)
}

func (p *ConcurrentPipeline) Start() error {
	for _, w := range p.workers {
		if err := w.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (p *ConcurrentPipeline) Stop() error {
	for _, w := range p.workers {
		w.Stop()
	}

	lastIdx := len(p.workers) - 1
	close(p.workers[lastIdx].Output())
	return nil
}

func (p *ConcurrentPipeline) Output() <-chan Message {
	sz := len(p.workers)
	return p.workers[sz-1].Output()
}

func (p *ConcurrentPipeline) Input() chan<- Message {
	return p.workers[0].Input()
}
