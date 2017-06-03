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
}

func (q *Queue) Push(i interface{}) {
	q.Lock()
	defer q.Unlock()

	e := Entry{
		Seq:   q.nextSeq,
		Value: i,
	}
	q.nextSeq++
	q.entries = append(q.entries, &e)
}

func (q *Queue) Pop(seq uint64) (*Entry, bool) {
	q.Lock()
	defer q.Unlock()

	q.ack(seq)

	if len(q.entries) > 0 {
		entry := q.entries[0]
		if entry.Seq+1 > q.unacknowledgedSeq {
			q.unacknowledgedSeq = entry.Seq + 1
		}
		return entry, true
	}
	return nil, false
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

func seqdiff(a, b uint64) uint64 {
	if a <= b {
		return b - a
	} else {
		return math.MaxUint64 - (a - b + 1)
	}
}

func inseq(a, b, seq uint64) bool {
	if a <= b {
		return a <= seq && seq < b
	} else {
		return a <= seq || seq < b
	}
}

// TODO: chan part...
