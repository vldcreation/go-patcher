package selector

import (
	"database/sql"
	"strings"

	"github.com/vldcreation/go-patcher/common"
	"github.com/vldcreation/go-patcher/placeholder"
)

type SQLSelect struct {
	// fields is the fields to update in the SQL statement
	fields []string

	// args is the arguments to use in the SQL statement
	args []any

	// db is the database connection to use
	db *sql.DB

	// tagName is the tag name to look for in the struct. This is an override from the default tag "db"
	tagName string

	// table is the table name to use in the SQL statement
	table string

	// whereSql is the where clause to use in the SQL statement
	whereSql *strings.Builder

	// whereArgs is the arguments to use in the where clause
	whereArgs []any

	// joinSql is the join clause to use in the SQL statement
	joinSql *strings.Builder

	// joinArgs is the arguments to use in the join clause
	joinArgs []any

	// limit is the limit for the SQL query
	limit int

	// offset is the offset for the SQL query
	offset int

	// placeholderType is the type of placeholder to use in the SQL query
	placeholderType placeholder.Type
}

func newSelectDefaults(opts ...SelectOpt) *SQLSelect {
	s := &SQLSelect{
		fields:    make([]string, 0),
		args:      make([]any, 0),
		db:        nil,
		tagName:   common.DefaultDbTagName,
		table:     "",
		whereSql:  new(strings.Builder),
		whereArgs: make([]any, 0),
		joinSql:   new(strings.Builder),
		joinArgs:  make([]any, 0),
		limit:           0,
		offset:          0,
		placeholderType: placeholder.Question,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// New creates a new selector.
func New(opts ...SelectOpt) *SQLSelect {
	return newSelectDefaults(opts...)
}

// From sets the fields to select from a struct.
func (s *SQLSelect) From(st any) *SQLSelect {
	s.fields = s.parseFields(st)
	return s
}
