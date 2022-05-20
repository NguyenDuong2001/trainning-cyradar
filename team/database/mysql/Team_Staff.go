package mysql

import (
	"github.com/google/uuid"
	"log"
)

func FindStaffInTeam(id uuid.UUID, DB *Mysql) []uuid.UUID {
	var ids []uuid.UUID
	data, err := DB.mysql.Query("select staff from team_staff where team =?", id)
	if err != nil {
		log.Fatal(err)
	}
	for data.Next() {
		var v uuid.UUID
		data.Scan(&v)
		ids = append(ids, v)
	}
	return ids
}

func AddManyStaffInTeam(staffs []uuid.UUID, team uuid.UUID, DB *Mysql) {
	query := "Insert into team_staff(staff, team) value"

	for i, staff := range staffs {
		query += "('" + staff.String() + "','" + team.String() + "')"
		if i == len(staffs)-1 {
			query += ";"
		} else {
			query += ","
		}
	}
	_, e := DB.mysql.Exec(query)
	if e != nil {
		log.Fatal(e)
	}
}
func DeleteAllStaffInTeam(team uuid.UUID, DB *Mysql) {
	stmt, err := DB.mysql.Prepare("DELETE from team_staff where team=?")
	_, err = stmt.Exec(team)
	if err != nil {
		log.Fatal(err)
	}
}
