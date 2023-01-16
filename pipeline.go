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
	//return &ConcurrentPipeline{workers: make([]stageWorker, 0)}
	return &ConcPipeline{workers: make([]concWorker, 0)}
}

type ConcPipeline struct {
	workers []concWorker
}

func (p *ConcPipeline) AddStage(stage Stage, opt *StageOptions) {
	if opt == nil {
		opt = &StageOptions{Concurrency: 10} // should be > 1
	}

	in := make(chan Message, 10)
	out := make(chan Message, 10)

	// if it's not the first stage the input channel must be the output of the previous stage
	for _, w := range p.workers {
		in = w.Output()
	}

	worker := newConcWorker(opt.Concurrency, in, out, stage)
	p.workers = append(p.workers, worker)
}

func (p *ConcPipeline) Start() error {
	for _, w := range p.workers {
		w.Start()
	}

	return nil
}

func (p *ConcPipeline) Stop() error {
	for _, w := range p.workers {
		w.Stop()
	}

	lastIdx := len(p.workers) - 1
	close(p.workers[lastIdx].Output())
	return nil
}

func (p *ConcPipeline) Output() <-chan Message {
	sz := len(p.workers)
	return p.workers[sz-1].Output()
}

func (p *ConcPipeline) Input() chan<- Message {
	return p.workers[0].Input()
}
