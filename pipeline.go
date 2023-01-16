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
	return &ConcPipeline{Workers: make([]*concWorker, 0)}
}

type ConcPipeline struct {
	In  chan<- Message
	Out <-chan Message

	Workers []*concWorker
}

func (p *ConcPipeline) AddStage(stage Stage, opt *StageOptions) {
	if opt == nil {
		opt = &StageOptions{Concurrency: 10} // should be > 1
	}

	in := make(chan Message, 10)
	out := make(chan Message, 10)

	if len(p.Workers) == 0 {
		p.In = in
	} else {
		// if it's not the first stage the input channel must be the output of the previous stage
		in = p.last().Output()
	}

	p.Out = out

	worker := newConcWorker(opt.Concurrency, in, out, stage)
	p.Workers = append(p.Workers, worker)
}

func (p *ConcPipeline) Start() error {
	for _, w := range p.Workers {
		w.Start()
	}

	return nil
}

func (p *ConcPipeline) last() *concWorker {
	lastIdx := len(p.Workers) - 1
	return p.Workers[lastIdx]
}

func (p *ConcPipeline) Stop() error {
	for _, w := range p.Workers {
		w.Stop()
	}

	close(p.last().Output())
	return nil
}

func (p *ConcPipeline) Input() chan<- Message {
	return p.In
}

func (p *ConcPipeline) Output() <-chan Message {
	return p.Out
}
