package seqqueue

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func description(q *Queue) string {
	return fmt.Sprintf("%+v %+v %+v", q.nextSeq, q.unacknowledgedSeq, q.entries)
}

func isEqual(a, b *Queue) bool {
	return a.nextSeq == b.nextSeq && a.unacknowledgedSeq == b.unacknowledgedSeq && reflect.DeepEqual(a.entries, b.entries)
}

func TestPush0(t *testing.T) {
	found := Queue{}
	expected := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 0,
		entries: []*Entry{
			&Entry{
				Seq:   0,
				Value: 127,
			},
		},
	}
	found.Push(127)

	if !isEqual(&expected, &found) {
		t.Errorf("Expected: %+v,\nfound: %+v", description(&expected), description(&found))
	}
}

func TestPop0(t *testing.T) {
	expectedEntry := &Entry{
		Seq:   0,
		Value: 127,
	}
	expected := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 1,
		entries: []*Entry{
			expectedEntry,
		},
	}
	found := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 0,
		entries: []*Entry{
			expectedEntry,
		},
	}
	entry, ok := found.Pop(0)

	if !ok {
		t.Errorf("Bad pop %+v", description(&found))
	}

	if !reflect.DeepEqual(entry, expectedEntry) {
		t.Errorf("Pop returns bad Entry %+v", entry)
	}

	if !isEqual(&expected, &found) {
		t.Errorf("Expected: %+v,\nfound: %+v", description(&expected), description(&found))
	}
}

func TestPop1(t *testing.T) {
	expected := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 0,
		entries:           []*Entry{},
	}
	found := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 0,
		entries:           []*Entry{},
	}
	entry, ok := found.Pop(0)

	if ok {
		t.Errorf("Bad pop %+v %+v", description(&found), entry)
	}

	if !isEqual(&expected, &found) {
		t.Errorf("Expected: %+v,\nfound: %+v", description(&expected), description(&found))
	}
}

func TestAck0(t *testing.T) {
	expected := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 1,
		entries:           []*Entry{},
	}
	found := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 1,
		entries: []*Entry{
			&Entry{
				Seq:   0,
				Value: 127,
			},
		},
	}

	found.ack(0)

	if !isEqual(&expected, &found) {
		t.Errorf("Expected: %+v,\nfound: %+v", description(&expected), description(&found))
	}
}

func TestAck1(t *testing.T) {
	expected := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 1,
		entries: []*Entry{
			&Entry{
				Seq:   0,
				Value: 127,
			},
		},
	}
	found := Queue{
		nextSeq:           1,
		unacknowledgedSeq: 1,
		entries: []*Entry{
			&Entry{
				Seq:   0,
				Value: 127,
			},
		},
	}

	found.ack(1)

	if !isEqual(&expected, &found) {
		t.Errorf("Expected: %+v,\nfound: %+v", description(&expected), description(&found))
	}
}

func TestInseq0(t *testing.T) {
	if inseq(0, 0, 0) {
		t.Errorf("Bad inseq")
	}
	if !inseq(0, 1, 0) {
		t.Errorf("Bad inseq")
	}
	if !inseq(0, 10, 5) {
		t.Errorf("Bad inseq")
	}
	if inseq(10, 0, 5) {
		t.Errorf("Bad inseq")
	}
	if !inseq(10, 9, 5) {
		t.Errorf("Bad inseq")
	}
	if inseq(0, 10, 15) {
		t.Errorf("Bad inseq")
	}
	if !inseq(10, 0, 15) {
		t.Errorf("Bad inseq")
	}
}

func TestSeqdiff0(t *testing.T) {
	if seqdiff(0, 0) != 0 {
		t.Error("Bad seqdiff")
	}
	if seqdiff(0, 1) != 1 {
		t.Error("Bad seqdiff")
	}
	if seqdiff(0, 10) != 10 {
		t.Error("Bad seqdiff")
	}
	if seqdiff(10, 0) != math.MaxUint64-11 {
		t.Error("Bad seqdiff")
	}
}