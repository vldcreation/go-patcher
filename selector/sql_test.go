package selector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vldcreation/go-patcher/common"
	"github.com/vldcreation/go-patcher/placeholder"
)

func TestGenerateSQL(t *testing.T) {
	type testCase struct {
		name     string
		opts     func(wherer *mockWherer) []SelectOpt
		expected string
		args     []any
	}

	testCases := []testCase{
		{
			name: "simple select",
			opts: func(wherer *mockWherer) []SelectOpt {
				return []SelectOpt{
					WithTable("users"),
				}
			},
			expected: "SELECT * FROM users",
			args:     []any{},
		},
		{
			name: "select with fields",
			opts: func(wherer *mockWherer) []SelectOpt {
				return []SelectOpt{
					WithTable("users"),
					WithFields("id", "name"),
				}
			},
			expected: "SELECT id, name FROM users",
			args:     []any{},
		},
		{
			name: "select with where",
			opts: func(wherer *mockWherer) []SelectOpt {
				wherer.On("Where").Return("name = ?", []any{"test"})
				wherer.On("WhereType").Return(common.WhereTypeAnd)
				return []SelectOpt{
					WithTable("users"),
					WithWhere(wherer),
				}
			},
			expected: "SELECT * FROM users\nWHERE (1=1)\nAND name = ?",
			args:     []any{"test"},
		},
		{
			name: "select with limit",
			opts: func(wherer *mockWherer) []SelectOpt {
				return []SelectOpt{
					WithTable("users"),
					WithLimit(10),
				}
			},
			expected: "SELECT * FROM users\nLIMIT 10",
		},
		{
			name: "select with offset",
			opts: func(wherer *mockWherer) []SelectOpt {
				return []SelectOpt{
					WithTable("users"),
					WithOffset(10),
				}
			},
			expected: "SELECT * FROM users\nOFFSET 10",
		},
		{
			name: "select with limit and offset",
			opts: func(wherer *mockWherer) []SelectOpt {
				return []SelectOpt{
					WithTable("users"),
					WithLimit(10),
					WithOffset(5),
				}
			},
			expected: "SELECT * FROM users\nLIMIT 10\nOFFSET 5",
		},
		{
			name: "select with where and dollar placeholder",
			opts: func(wherer *mockWherer) []SelectOpt {
				wherer.On("Where").Return("name = ?", []any{"test"})
				wherer.On("WhereType").Return(common.WhereTypeAnd)
				return []SelectOpt{
					WithTable("users"),
					WithWhere(wherer),
					WithPlaceholderFormat(placeholder.Dollar),
				}
			},
			expected: "SELECT * FROM users\nWHERE (1=1)\nAND name = $1",
			args:     []any{"test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wherer := &mockWherer{}
			wherer.On("WhereType").Maybe().Return(common.WhereTypeAnd)
			opts := tc.opts(wherer)

			selector := New(opts...)
			sql, args, err := selector.GenerateSQL()

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, strings.TrimSpace(sql))
			if tc.args != nil {
				assert.Equal(t, tc.args, args)
			}
			wherer.AssertExpectations(t)
		})
	}
}
