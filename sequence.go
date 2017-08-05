package piecechains

import (
	"bytes"
	"container/list"
	"fmt"
)

type Sequence struct {
	spans *list.List
	// modification buffer
	editsBuffer *bytes.Buffer
}

func NewSequence() *Sequence {
	return &Sequence{
		spans:       list.New(),
		editsBuffer: new(bytes.Buffer),
	}
}

func (s *Sequence) Len() int {
	len := 0
	for e := s.spans.Front(); e != nil; e = e.Next() {
		span := e.Value.(*Span)
		len = len + span.len
	}
	return len
}

// string representation
// this can be expensive
func (s *Sequence) String() string {
	buf := new(bytes.Buffer)
	for e := s.spans.Front(); e != nil; e = e.Next() {
		span := e.Value.(*Span)
		buf.Write(span.buffer.Bytes()[span.offset : span.offset+span.len])
	}
	return buf.String()
}

// for a given index into the sequence, returns the
// span element and the offset.
//
// a nil span and offset > 0 means index is out of bonds
func (s *Sequence) spanElementForIndex(idx int) (*list.Element, int) {
	currLen := 0
	for e := s.spans.Front(); e != nil; e = e.Next() {
		span := e.Value.(*Span)
		if idx >= currLen && idx < currLen+span.len {
			return e, idx - currLen
		}
		currLen = currLen + span.len
	}
	return nil, idx - currLen
}

func (s *Sequence) spanForIndex(idx int) (*Span, int) {
	e, offset := s.spanElementForIndex(idx)
	if e != nil {
		span := e.Value.(*Span)
		return span, offset
	}
	return nil, offset
}

func (s *Sequence) NewEditSpan(content []byte) (*Span, error) {
	origLen := s.editsBuffer.Len()
	n, err := s.editsBuffer.Write(content)
	if err != nil {
		return nil, err
	}
	span := Span{offset: origLen,
		len: n, buffer: s.editsBuffer}
	return &span, nil
}

func (span *Span) Split(idx int) (*Span, *Span) {
	span1 := Span{offset: span.offset,
		len: idx, buffer: span.buffer}
	span2 := Span{offset: span1.len,
		len: span.len - idx, buffer: span.buffer}
	return &span1, &span2
}

// inserts string at index.
// errors if index is out of bonds
func (s *Sequence) Insert(idx int, content []byte) error {
	e, offset := s.spanElementForIndex(idx)
	// insert at a boundary
	if offset == 0 {
		span, err := s.NewEditSpan(content)
		if err != nil {
			return err
		}
		// insert in between
		if e != nil {
			s.spans.InsertBefore(span, e)
		} else { // first or last span
			s.spans.PushBack(span)
		}
	} else {
		if e == nil {
			return fmt.Errorf("Out of bounds")
		}
		prev := e
		oldSpan := prev.Value.(*Span)
		span1, span2 := oldSpan.Split(offset)
		e := s.spans.InsertAfter(span1, prev)
		span, err := s.NewEditSpan(content)
		if err != nil {
			return err
		}
		e = s.spans.InsertAfter(span, e)
		e = s.spans.InsertAfter(span2, e)
		s.spans.Remove(prev)
	}
	return nil
}

// a span
type Span struct {
	offset int
	len    int
	buffer *bytes.Buffer
}
