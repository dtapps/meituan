package dorm

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq"
)

func NewBeegoPostgresClient(dns string) (*BeegoClient, error) {

	c := &BeegoClient{}

	err := orm.RegisterDataBase("default", "postgres", dns)
	if err != nil {
		return nil, err
	}
	c.Db = orm.NewOrm()

	return c, err
}
