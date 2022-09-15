package dorm

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

func NewBeegoMysqlClient(dns string) (*BeegoClient, error) {

	c := &BeegoClient{}

	err := orm.RegisterDataBase("default", "mysql", dns)
	if err != nil {
		return nil, err
	}
	c.Db = orm.NewOrm()

	return c, err
}
