package piecechains

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptySequence(t *testing.T) {
	s := NewSequence()
	assert.Equal(t, 0, s.Len())
	assert.Equal(t, "", s.String())

	span, idx := s.spanForIndex(0)
	assert.Nil(t, span)
	assert.Equal(t, 0, idx)
}

func TestInsert(t *testing.T) {
	s := NewSequence()
	err := s.Insert(1, []byte("Hello"))
	assert.NotNil(t, err)

	err = s.Insert(0, []byte("Hello"))
	assert.Nil(t, err)
	assert.Equal(t, "Hello", s.String())

	err = s.Insert(6, []byte("Hello"))
	assert.NotNil(t, err)
	// insert did not happen
	assert.Equal(t, "Hello", s.String())

	// insert at end
	err = s.Insert(5, []byte("Bye"))
	assert.Nil(t, err)
	assert.Equal(t, "HelloBye", s.String())

	err = s.Insert(5, []byte("Middle"))
	assert.Nil(t, err)
	assert.Equal(t, "HelloMiddleBye", s.String())
}
