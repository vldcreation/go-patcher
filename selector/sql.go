package selector

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vldcreation/go-patcher/common"
)

func New(opts ...SelectOpt) *SQLSelect {
	return newSelectDefaults(opts...)
}

func (s *SQLSelect) GenerateSQL() (sqlStr string, args []any, err error) {
	sqlBuilder := new(strings.Builder)
	sqlBuilder.WriteString("SELECT ")

	if len(s.fields) == 0 {
		sqlBuilder.WriteString("*")
	} else {
		sqlBuilder.WriteString(strings.Join(s.fields, ", "))
	}

	sqlBuilder.WriteString(" FROM ")
	sqlBuilder.WriteString(s.table)
	sqlBuilder.WriteString("\n")

	if s.joinSql.String() != "" {
		sqlBuilder.WriteString(s.joinSql.String())
	}

	if s.whereSql.String() != "" {
		sqlBuilder.WriteString("WHERE (1=1)\n")
		sqlBuilder.WriteString(s.whereSql.String())
	}

	if s.limit > 0 {
		sqlBuilder.WriteString(fmt.Sprintf("LIMIT %d\n", s.limit))
	}

	if s.offset > 0 {
		sqlBuilder.WriteString(fmt.Sprintf("OFFSET %d\n", s.offset))
	}

	sqlArgs := s.joinArgs
	sqlArgs = append(sqlArgs, s.whereArgs...)

	return strings.TrimSpace(sqlBuilder.String()), sqlArgs, nil
}

func (s *SQLSelect) Perform(dest any) error {
	if err := s.validate(); err != nil {
		return err
	}

	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("destination must be a pointer to a slice")
	}

	sliceType := destVal.Elem().Type().Elem()

	sqlStr, args, err := s.GenerateSQL()
	if err != nil {
		return fmt.Errorf("generating SQL: %w", err)
	}

	rows, err := s.db.Query(sqlStr, args...)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	slice := reflect.MakeSlice(destVal.Elem().Type(), 0, 0)

	for rows.Next() {
		newVal := reflect.New(sliceType).Interface()
		if err := s.scan(rows, newVal); err != nil {
			return fmt.Errorf("scanning row: %w", err)
		}
		slice = reflect.Append(slice, reflect.ValueOf(newVal).Elem())
	}

	destVal.Elem().Set(slice)

	return nil
}

func (s *SQLSelect) validate() error {
	if s.db == nil {
		return fmt.Errorf("database connection is not set")
	}
	if s.table == "" {
		return fmt.Errorf("table name is not set")
	}
	return nil
}

func (s *SQLSelect) scan(rows common.Rows, dest any) error {
	destVal := reflect.ValueOf(dest).Elem()
	destType := destVal.Type()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := make([]any, len(columns))
	fieldMap := make(map[string]reflect.Value)

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		tag := field.Tag.Get(s.tagName)
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		fieldMap[tag] = destVal.Field(i)
	}

	for i, col := range columns {
		if field, ok := fieldMap[col]; ok {
			values[i] = field.Addr().Interface()
		} else {
			var discard any
			values[i] = &discard
		}
	}

	return rows.Scan(values...)
}
