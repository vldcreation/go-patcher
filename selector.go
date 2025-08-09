package patcher

import "github.com/jacobbrewer1/patcher/selector"

// NewSelector creates a new selector.
//
// This is a convenience function that wraps the selector package.
func NewSelector(opts ...selector.SelectOpt) *selector.SQLSelect {
	return selector.New(opts...)
}
