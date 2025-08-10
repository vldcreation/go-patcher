package selector

import (
	"database/sql"

	"github.com/vldcreation/go-patcher/common"
	"github.com/vldcreation/go-patcher/placeholder"
)

const (
	DefaultPage  = 1
	DefaultLimit = 50
)

type SelectOpt func(*SQLSelect)

// WithTable sets the table name to use in the SQL statement
func WithTable(table string) SelectOpt {
	return func(s *SQLSelect) {
		s.table = table
	}
}

// WithDB sets the database connection to use
func WithDB(db *sql.DB) SelectOpt {
	return func(s *SQLSelect) {
		s.db = db
	}
}

// WithFields sets the fields to select in the SQL statement. If no fields are provided, it will default to "*"
func WithFields(fields ...string) SelectOpt {
	return func(s *SQLSelect) {
		s.fields = fields
	}
}

// WithWhere sets the where clause to use in the SQL statement
func WithWhere(where common.Wherer) SelectOpt {
	return func(s *SQLSelect) {
		common.AppendWhere(where, s.whereSql, &s.whereArgs)
	}
}

// WithJoin sets the join clause to use in the SQL statement
func WithJoin(join common.Joiner) SelectOpt {
	return func(s *SQLSelect) {
		common.AppendJoin(join, s.joinSql, &s.joinArgs)
	}
}

// WithLimit sets the limit for the SQL query.
func WithLimit(limit int) SelectOpt {
	return func(s *SQLSelect) {
		s.limit = limit
	}
}

// WithOffset sets the offset for the SQL query.
func WithOffset(offset int) SelectOpt {
	return func(s *SQLSelect) {
		s.offset = offset
	}
}

// WithPlaceholderFormat sets the placeholder format for the SQL query.
func WithPlaceholderFormat(placeholderType placeholder.Type) SelectOpt {
	return func(s *SQLSelect) {
		s.placeholderType = placeholderType
	}
}

func WithPagination(page, limit int) SelectOpt {
	return func(s *SQLSelect) {
		if page < 1 {
			page = DefaultPage
		}
		if limit < 1 {
			limit = DefaultLimit
		}
		s.limit = limit
		s.offset = (page - 1) * limit
	}
}
