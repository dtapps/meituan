package dorm

import (
	"github.com/beego/beego/v2/client/orm"
)

type BeegoClient struct {
	Db orm.Ormer // 驱动
}
