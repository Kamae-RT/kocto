package kocto

import (
	"log"

	"github.com/sourcegraph/conc/pool"
)

type Stage interface {
	Process(message Message) ([]Message, error)
}

type StageOptions struct {
	Concurrency int
}

type concWorker struct {
	pool       *pool.Pool
	concurrent int
	in         chan Message
	out        chan Message
	stage      Stage
}

func newConcWorker(concurrentWorkers int, in chan Message, out chan Message, stage Stage) concWorker {
	return concWorker{
		pool:       pool.New().WithMaxGoroutines(concurrentWorkers),
		concurrent: concurrentWorkers,
		in:         in,
		out:        out,
		stage:      stage,
	}
}

func (w *concWorker) Start() error {
	for i := 0; i < w.concurrent; i++ {
		w.pool.Go(func() {
			log.Println("working")
			for msg := range w.in {
				msg := msg

				log.Println("processing")
				res, err := w.stage.Process(msg)
				if err != nil {
				}

				for _, m := range res {
					w.out <- m
				}
			}
		})
	}
	return nil
}

func (w *concWorker) Stop() error {
	close(w.Input())
	w.pool.Wait()

	return nil
}

func (w *concWorker) Input() chan Message {
	return w.in
}

func (w *concWorker) Output() chan Message {
	return w.out
}
