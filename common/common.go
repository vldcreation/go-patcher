package common

import (
	"reflect"
	"strings"
)

const (
	DefaultDbTagName = "db"
	DBTagPrimaryKey  = "pk"
	TagOptsName      = "patcher"
	TagOptSeparator  = ","
	TagOptSkip       = "-"
	TagOptOmitempty  = "omitempty"
)

type IgnoreFieldsFunc func(field *reflect.StructField) bool

type Rows interface {
	Close() error
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
}

type Wherer interface {
	Where() (string, []any)
}

type WhereTyper interface {
	Wherer
	WhereType() WhereType
}

type WhereType string

const (
	WhereTypeAnd WhereType = "AND"
	WhereTypeOr  WhereType = "OR"
)

func (w WhereType) IsValid() bool {
	switch w {
	case WhereTypeAnd, WhereTypeOr:
		return true
	}
	return false
}

type Joiner interface {
	Join() (string, []any)
}

// IsValidType checks if the given value is of a type that can be stored as a database field.
func IsValidType(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Struct, reflect.Ptr:
		return true
	default:
		return false
	}
}

func AppendWhere(where Wherer, builder *strings.Builder, args *[]any) {
	if where == nil {
		return
	}
	wSQL, fwArgs := where.Where()
	if fwArgs == nil {
		fwArgs = make([]any, 0)
	}
	wtStr := WhereTypeAnd // default to AND
	wt, ok := where.(WhereTyper)
	if ok && wt.WhereType().IsValid() {
		wtStr = wt.WhereType()
	}
	builder.WriteString(string(wtStr) + " ")
	builder.WriteString(strings.TrimSpace(wSQL))
	builder.WriteString("\n")
	*args = append(*args, fwArgs...)
}

func AppendJoin(join Joiner, builder *strings.Builder, args *[]any) {
	if join == nil {
		return
	}
	jSQL, jArgs := join.Join()
	if jArgs == nil {
		jArgs = make([]any, 0)
	}
	builder.WriteString(strings.TrimSpace(jSQL))
	builder.WriteString("\n")
	*args = append(*args, jArgs...)
}
