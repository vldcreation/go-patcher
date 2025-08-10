package placeholder

import (
	"bytes"
	"strconv"
	"strings"
)

type Type int

const (
	Question Type = iota
	Dollar
)

func Rebind(sql string, placeholderType Type) string {
	if placeholderType == Question {
		return sql
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(sql)))
	i := 1
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // A ?? escaped question mark
			buf.WriteString(sql[:p+1])
			sql = sql[p+2:]
		} else {
			buf.WriteString(sql[:p])
			buf.WriteString("$")
			buf.WriteString(strconv.Itoa(i))
			sql = sql[p+1:]
			i++
		}
	}

	buf.WriteString(sql)
	return buf.String()
}
