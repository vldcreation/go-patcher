package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Embedded struct {
	Name string `db:"name"`
}

func TestFrom(t *testing.T) {
	type testCase struct {
		name     string
		from     any
		expected []string
	}

	testCases := []testCase{
		{
			name: "simple struct",
			from: struct {
				ID   int    `db:"id"`
				Name string `db:"name"`
			}{},
			expected: []string{"id", "name"},
		},
		{
			name: "struct with anonymous field",
			from: struct {
				ID int `db:"id"`
				Embedded
			}{},
			expected: []string{"id", "name"},
		},
		{
			name: "struct with ignored field",
			from: struct {
				ID   int    `db:"id"`
				Name string `db:"-"`
			}{},
			expected: []string{"id"},
		},
		{
			name: "struct with empty tag",
			from: struct {
				ID   int `db:"id"`
				Name string
			}{},
			expected: []string{"id"},
		},
		{
			name: "struct with comma",
			from: struct {
				ID   int    `db:"id,omitempty"`
				Name string `db:"name,omitempty"`
			}{},
			expected: []string{"id", "name"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			selector := New()
			selector.From(tc.from)

			assert.Equal(t, tc.expected, selector.fields)
		})
	}
}
