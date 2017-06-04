package seqqueue

import (
	"math"
	"sync"
)

type Queue struct {
	sync.Mutex
	nextSeq           uint64
	unacknowledgedSeq uint64
	entries           []*Entry

	in  chan interface{}
	out chan *Entry
}

func NewQueue() *Queue {
	q := Queue{
		entries: []*Entry{},
		in:      make(chan interface{}),
		out:     make(chan *Entry),
	}
	go q.loop()
	return &q
}

func (q *Queue) Dispose() {
	close(q.in)
}

func (q *Queue) In() chan interface{} {
	return q.in
}

func (q *Queue) Out(seq uint64) <-chan *Entry {
	q.ack(seq)
	q.fakeRead()
	return q.out
}

func (q *Queue) ack(seq uint64) {
	if len(q.entries) == 0 {
		return
	}

	firstEntity := q.entries[0]

	if !inseq(firstEntity.Seq, q.unacknowledgedSeq, seq) {
		return
	}

	acked := seqdiff(firstEntity.Seq, seq) + 1
	q.entries = q.entries[acked:]
}

func (q *Queue) fakeRead() {
	select {
	case <-q.out:
	default:
	}
}

func (q *Queue) loop() {
	var in chan interface{}
	var out chan *Entry

	in = q.in

	for {
		var entry *Entry
		if len(q.entries) > 0 {
			out = q.out
			entry = q.entries[0]
		} else {
			out = nil
		}

		if in == nil && out == nil {
			close(q.out)
			return
		}

		select {
		case i, ok := <-in:
			if ok {
				q.push(i)
			} else {
				in = nil
			}
		case out <- entry:
			q.unacknowledgedSeq = entry.Seq + 1
			continue
		}
	}
}

func (q *Queue) push(i interface{}) {
	e := Entry{
		Seq:   q.nextSeq,
		Value: i,
	}
	q.nextSeq++
	q.entries = append(q.entries, &e)
}

func seqdiff(a, b uint64) uint64 {
	if a <= b {
		return b - a
	}
	return math.MaxUint64 - (a - b + 1)
}

func inseq(a, b, seq uint64) bool {
	if a <= b {
		return a <= seq && seq < b
	}
	return a <= seq || seq < b
}
