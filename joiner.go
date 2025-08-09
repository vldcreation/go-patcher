package patcher

import "github.com/vldcreation/go-patcher/common"

type Joiner = common.Joiner

type joinStringOption struct {
	join string
	args []any
}

func (j *joinStringOption) Join() (sqlStr string, args []any) {
	return j.join, j.args
}
