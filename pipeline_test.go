package kocto_test

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/matryer/is"
	"kamaesoft.visualstudio.com/kocto/_git/kocto"
)

type doubler struct {
	t     *testing.T
	l     sync.Mutex
	count int
}

func (s *doubler) Process(msg kocto.Message) ([]kocto.Message, error) {
	s.l.Lock()
	defer s.l.Unlock()

	s.count += 1

	m := msg.(int)
	r := m * 2

	s.t.Logf("doubler: processing %d => %d", m, r)
	return []kocto.Message{r}, nil
}

type power struct {
	t     *testing.T
	l     sync.Mutex
	count int
}

func (s *power) Process(msg kocto.Message) ([]kocto.Message, error) {
	s.l.Lock()
	defer s.l.Unlock()

	s.count += 1

	m := msg.(int)
	r := int(math.Pow(float64(m), 2))
	s.t.Logf("power: processing %d => %d", m, r)

	return []kocto.Message{r}, nil
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

	d := &doubler{t, sync.Mutex{}, 0}
	pw := &power{t, sync.Mutex{}, 0}

	p.AddStage(d, &kocto.StageOptions{Concurrency: 2})
	p.AddStage(pw, &kocto.StageOptions{Concurrency: 2})

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
			msg := msg
			//t.Log("out:", msg)
			msgs.Add(msg)
		}
	}()

	numMsgs := 100
	expected := numMsgs

	for i := 1; i <= numMsgs; i++ {
		p.Input() <- i
	}

	time.Sleep(time.Millisecond * 100)
	p.Stop()

	is := is.NewRelaxed(t)

	is.Equal(d.count, numMsgs)  // doubler should have seen 100 msgs
	is.Equal(pw.count, numMsgs) // power should have seen 100 msgs

	msgs.l.Lock()
	if len(msgs.messages) != expected {
		is.Equal(len(msgs.messages), expected)
	}

	sum := 0
	for _, m := range msgs.messages {
		sum += m.(int)
	}

	is.Equal(sum, 1353400) // the sum of (x * 2) ^2 , where x E [1, 100] should be 1353400

	msgs.l.Unlock()
}
