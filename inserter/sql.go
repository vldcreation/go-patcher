package inserter

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/vldcreation/go-patcher/common"
	"github.com/vldcreation/go-patcher/placeholder"
)

func NewBatch(resources []any, opts ...BatchOpt) *SQLBatch {
	b := newBatchDefaults(opts...)
	for _, opt := range opts {
		opt(b)
	}
	b.genBatch(resources)
	return b
}

func (b *SQLBatch) genBatch(resources []any) {
	uniqueFields := make(map[string]struct{})

	for _, r := range resources {
		t := reflect.TypeOf(r)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			continue
		}

		v := reflect.ValueOf(r)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		for i := range t.NumField() {
			f := t.Field(i)
			fVal := v.Field(i)

			if !common.IsValidType(fVal) || !f.IsExported() || b.checkSkipField(&f) {
				continue
			}

			tag := f.Tag.Get(b.tagName)
			if tag == common.TagOptSkip {
				continue
			}

			if tag == "" {
				tag = f.Name
			} else {
				tag = strings.Split(tag, common.TagOptSeparator)[0]
			}

			b.args = append(b.args, b.getFieldValue(fVal, &f))

			if _, ok := uniqueFields[tag]; ok {
				continue
			}

			b.fields = append(b.fields, tag)
			uniqueFields[tag] = struct{}{}
		}
	}
}

func (b *SQLBatch) getFieldValue(v reflect.Value, f *reflect.StructField) any {
	if f.Type.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	} else if f.Type.Kind() == reflect.Ptr {
		return v.Elem().Interface()
	}

	return v.Interface()
}

func (b *SQLBatch) GenerateSQL() (sqlStr string, args []any, err error) {
	if err := b.validateSQLGen(); err != nil {
		return "", nil, err
	}

	sqlBuilder := new(strings.Builder)
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(b.table)
	sqlBuilder.WriteString(" (")
	sqlBuilder.WriteString(strings.Join(b.fields, ", "))
	sqlBuilder.WriteString(") VALUES ")

	sqlBuilder.WriteString(b.buildPlaceholders())

	return sqlBuilder.String(), b.args, nil
}

func (b *SQLBatch) buildPlaceholders() string {
	if b.placeholderType == placeholder.Dollar {
		return b.buildDollarPlaceholders()
	}

	return b.buildQuestionMarkPlaceholders()
}

func (b *SQLBatch) buildQuestionMarkPlaceholders() string {
	placeholder := "(" + strings.Repeat("?, ", len(b.fields)-1) + "?)"
	placeholders := strings.Repeat(placeholder+", ", len(b.args)/len(b.fields))
	return placeholders[:len(placeholders)-2]
}

func (b *SQLBatch) buildDollarPlaceholders() string {
	placeholders := new(strings.Builder)
	rowCount := len(b.args) / len(b.fields)
	for i := 0; i < rowCount; i++ {
		placeholders.WriteString("(")
		for j := 0; j < len(b.fields); j++ {
			placeholders.WriteString("$")
			placeholders.WriteString(fmt.Sprintf("%d", i*len(b.fields)+j+1))
			if j < len(b.fields)-1 {
				placeholders.WriteString(", ")
			}
		}
		placeholders.WriteString(")")
		if i < rowCount-1 {
			placeholders.WriteString(", ")
		}
	}
	return placeholders.String()
}

func (b *SQLBatch) Perform() (sql.Result, error) {
	if err := b.validateSQLInsert(); err != nil {
		return nil, fmt.Errorf("validate SQL generation: %w", err)
	}

	sqlStr, args, err := b.GenerateSQL()
	if err != nil {
		return nil, fmt.Errorf("generate SQL: %w", err)
	}

	return b.db.Exec(sqlStr, args...)
}
