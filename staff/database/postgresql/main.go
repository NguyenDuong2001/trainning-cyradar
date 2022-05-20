package postgresql

import (
	"Basic/Trainning4/redis/staff/config"
	"database/sql"
)

type Postgresql struct {
	postgresql *sql.DB
}

func (DB *Postgresql) NewDB() {
	DB.postgresql = config.GetDBPostgresql()
}

func (DB *Postgresql) GetName() string {
	return "Postgresql"
}
