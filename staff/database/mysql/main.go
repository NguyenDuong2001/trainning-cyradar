package mysql

import (
	"Basic/Trainning4/redis/staff/config"
	"database/sql"
)

type Mysql struct {
	mysql *sql.DB
}

func (DB *Mysql) NewDB() {
	DB.mysql = config.GetDBMySql()
}

func (DB *Mysql) GetName() string {
	return "MySQL"
}
