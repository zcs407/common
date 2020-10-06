package conn

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DBClient struct {
	Engine *xorm.Engine
}

func NewDBClient(user, host, pwd, dbName string) (db *DBClient, err error) {
	db = &DBClient{}
	dataSourceName := fmt.Sprintf("%s:%s@%s/%s?charset=utf8", user, pwd, host, dbName)
	db.Engine, err = xorm.NewEngine("mysql", dataSourceName)
	err = db.Engine.Ping()
	return
}

func (d *DBClient) Close() {
	_ = d.Engine.Close()
}

func (d *DBClient) NewSession() xorm.Session {
	return d.NewSession()
}
