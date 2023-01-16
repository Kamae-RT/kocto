package kocto_test

import (
	"sync"
	"testing"
	"time"

	"kamaesoft.visualstudio.com/kocto/_git/kocto"
)

type doubler struct{}

func (s doubler) Process(msg kocto.Message) ([]kocto.Message, error) {
	return []kocto.Message{msg.(int) * 2}, nil
}

type power struct{}

func (s power) Process(msg kocto.Message) ([]kocto.Message, error) {
	return []kocto.Message{msg.(int) ^ 2}, nil
}

type splitter struct{}

func (s splitter) Process(msg kocto.Message) ([]kocto.Message, error) {
	n := msg.(int) / 2

	return []kocto.Message{n, n}, nil
}

type msgsS struct {
	messages []kocto.Message
	l        sync.Mutex
}

func (m *msgsS) Add(msg kocto.Message) {
	m.l.Lock()
	defer m.l.Unlock()

	m.messages = append(m.messages, msg)
}

func TestPipeline(t *testing.T) {
	p := kocto.NewPipeline()

	p.AddStage(doubler{}, nil)
	p.AddStage(power{}, nil)
	p.AddStage(splitter{}, nil)

	if err := p.Start(); err != nil {
		t.Log(err)
		t.FailNow()
	}

	msgs := msgsS{
		messages: make([]kocto.Message, 0),
		l:        sync.Mutex{},
	}

	go func() {
		for msg := range p.Output() {
			t.Log("out: ", msg)
			msgs.Add(msg)
		}
	}()

	numMsgs := 100
	expected := numMsgs * 2

	for i := 1; i <= numMsgs; i++ {
		p.Input() <- i
		t.Log("In: ", i)
	}

	p.Stop()
	time.Sleep(time.Millisecond * 10)

	msgs.l.Lock()
	if len(msgs.messages) != expected {
		t.Log("wrong number of messages: ", len(msgs.messages), " expected ", expected)
		t.Fail()
	}
	msgs.l.Unlock()
}
