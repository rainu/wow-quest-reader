package processor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdIterator_Next(t *testing.T) {
	toTest := newQuestIter([]int64{1, 2})

	assert.Equal(t, int64(0), toTest.Next())
	assert.Equal(t, int64(3), toTest.Next())
	assert.Equal(t, int64(4), toTest.Next())
	assert.Equal(t, int64(5), toTest.Next())
}
