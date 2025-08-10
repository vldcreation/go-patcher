package selector

import (
	"reflect"
	"strings"
)

func (s *SQLSelect) parseFields(st any) []string {
	t := reflect.TypeOf(st)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	fields := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			fields = append(fields, s.parseFields(reflect.Zero(field.Type).Interface())...)
			continue
		}

		tag := field.Tag.Get(s.tagName)
		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		if len(parts) > 0 {
			fields = append(fields, parts[0])
		}
	}

	return fields
}
