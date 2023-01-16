package kocto

import (
	"github.com/sourcegraph/conc/pool"
)

type Stage interface {
	Process(message Message) ([]Message, error)
}

type StageOptions struct {
	Concurrency int
}

type pipeWorker struct {
	In  chan Message
	Out chan Message

	logger      Logger
	pool        *pool.Pool
	concurrency int

	stage Stage
}

func newPipeWorker(l Logger, c int, in chan Message, out chan Message, stage Stage) *pipeWorker {
	return &pipeWorker{
		logger:      l,
		pool:        pool.New().WithMaxGoroutines(c),
		concurrency: c,
		In:          in,
		Out:         out,
		stage:       stage,
	}
}

func (w *pipeWorker) Start() error {
	for i := 0; i < w.concurrency; i++ {
		w.pool.Go(func() {
			for msg := range w.In {
				msg := msg

				res, err := w.stage.Process(msg)
				if err != nil {
					w.logger.Errorw("unable to process message", "error", err)
					continue
				}

				for _, m := range res {
					w.Out <- m
				}
			}
		})
	}
	return nil
}

func (w *pipeWorker) Stop() error {
	close(w.In)
	w.pool.Wait()

	return nil
}
