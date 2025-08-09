package patcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vldcreation/go-patcher/selector"
)

func TestNewSelector(t *testing.T) {
	s := NewSelector(selector.WithTable("users"))
	assert.NotNil(t, s)
}
