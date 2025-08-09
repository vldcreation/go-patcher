package selector

import (
	"github.com/jacobbrewer1/patcher/common"
	"github.com/stretchr/testify/mock"
)

type mockWherer struct {
	mock.Mock
}

func (m *mockWherer) Where() (string, []any) {
	args := m.Called()
	return args.String(0), args.Get(1).([]any)
}

func (m *mockWherer) WhereType() common.WhereType {
	args := m.Called()
	if args.Get(0) == nil {
		return ""
	}
	return args.Get(0).(common.WhereType)
}
