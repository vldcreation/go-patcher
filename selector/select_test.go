package selector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vldcreation/go-patcher/selector"
)

func TestNewSelector(t *testing.T) {
	s := selector.New(selector.WithTable("users"))
	assert.NotNil(t, s)
}
