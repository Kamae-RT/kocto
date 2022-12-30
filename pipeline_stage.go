package kocto

import (
	"sync"
)

type Stage interface {
	Process(message Message) ([]Message, error)
}

type StageOptions struct {
	Concurrency int
}

type stageWorker struct {
	wg                *sync.WaitGroup
	concurrentWorkers int
	in                chan Message
	out               chan Message
	stage             Stage
}

func newStageWorker(concurrentWorkers int, in chan Message, out chan Message, stage Stage) stageWorker {
	return stageWorker{
		wg:                &sync.WaitGroup{},
		concurrentWorkers: concurrentWorkers,
		in:                in,
		out:               out,
		stage:             stage,
	}
}

func (w *stageWorker) Start() error {
	for i := 0; i < w.concurrentWorkers; i++ {
		w.wg.Add(1)

		// this function is required to avoid weird shadowing bugs
		go work(w.wg, w.Input(), w.Output(), w.stage)
	}

	return nil
}

func (w *stageWorker) Stop() error {
	close(w.Input())
	w.wg.Wait()

	return nil
}

func (w *stageWorker) Input() chan Message {
	return w.in
}

func (w *stageWorker) Output() chan Message {
	return w.out
}

func work(wg *sync.WaitGroup, in, out chan Message, stage Stage) {
	defer wg.Done()

	for msg := range in {
		result, err := stage.Process(msg)
		if err != nil {
			continue
		}

		for _, r := range result {
			out <- r
		}
	}
}
