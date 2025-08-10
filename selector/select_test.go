package selector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vldcreation/go-patcher/selector"
)

type User struct {
	ID           int    `db:"id"`
	IgnoredField string `db:"-"`
	EmptyTag     string `db:""`
	CommaTag     string `db:"comma,omitempty"`
}

type UserWhere struct {
	ID *int `db:"id"`
}

func NewUserWhere(id int) *UserWhere {
	return &UserWhere{
		ID: &id,
	}
}

func (u *UserWhere) Where() (sqlStr string, sqlArgs []any) {
	if u.ID == nil {
		return "", nil
	}
	return "id = ?", []any{*u.ID}
}

func TestNewSelector(t *testing.T) {
	s := selector.New(selector.WithTable("users"))
	assert.NotNil(t, s)
}

func TestSelector_From(t *testing.T) {
	type User struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	s := selector.New(selector.WithTable("users"))
	s.From(User{})

	sql, args, err := s.GenerateSQL()
	assert.NoError(t, err)
	expectedSQL := "SELECT id, name FROM users"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, args)
}

func TestSelector_FromWithEmbedded(t *testing.T) {
	type Embedded struct {
		Name string `db:"name"`
	}

	type User struct {
		ID int `db:"id"`
		Embedded
	}

	s := selector.New(selector.WithTable("users"))
	s.From(User{})

	sql, args, err := s.GenerateSQL()
	assert.NoError(t, err)
	expectedSQL := "SELECT id, name FROM users"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, args)
}

func TestSelector_FromWithIgnoredField(t *testing.T) {
	s := selector.New(selector.WithTable("users"))
	s.From(User{})

	sql, args, err := s.GenerateSQL()
	assert.NoError(t, err)
	expectedSQL := "SELECT id FROM users"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, args)
}

func TestSelector_FromWithEmptyTag(t *testing.T) {
	s := selector.New(selector.WithTable("users"))
	s.From(User{})

	sql, args, err := s.GenerateSQL()
	assert.NoError(t, err)
	expectedSQL := "SELECT id FROM users"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, args)
}

func TestSelector_FromWithCommaTag(t *testing.T) {
	s := selector.New(selector.WithTable("users"))
	s.From(User{})

	sql, args, err := s.GenerateSQL()
	assert.NoError(t, err)
	expectedSQL := "SELECT id, name FROM users"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, args)
}

func TestSelector_WithWhere(t *testing.T) {
	wherer := NewUserWhere(1)
	s := selector.New(selector.WithTable("users"), selector.WithWhere(wherer))
	s.From(User{})

	sql, args, err := s.GenerateSQL()
	assert.NoError(t, err)
	expectedSQL := "SELECT id, comma FROM users\nWHERE (1=1)\nAND id = ?"
	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, []any{1}, args)
}

func TestSelector_WithWhereNoValue(t *testing.T) {
	s := selector.New(selector.WithTable("users"), selector.WithWhere(&UserWhere{}))
	s.From(User{})

	sql, args, err := s.GenerateSQL()
	assert.NoError(t, err)
	expectedSQL := "SELECT id, comma FROM users"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, args)
}
