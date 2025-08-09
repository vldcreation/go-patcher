package patcher

import (
	"testing"

	"github.com/jacobbrewer1/patcher/selector"
	"github.com/stretchr/testify/assert"
)

func TestNewSelector(t *testing.T) {
	s := NewSelector(selector.WithTable("users"))
	assert.NotNil(t, s)
}
