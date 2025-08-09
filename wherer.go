package patcher

import "github.com/jacobbrewer1/patcher/common"

type Wherer = common.Wherer
type WhereTyper = common.WhereTyper
type WhereType = common.WhereType

const (
	WhereTypeAnd = common.WhereTypeAnd
	WhereTypeOr  = common.WhereTypeOr
)

type whereStringOption struct {
	where string
	args  []any
}

func (w *whereStringOption) Where() (sqlStr string, args []any) {
	return w.where, w.args
}
