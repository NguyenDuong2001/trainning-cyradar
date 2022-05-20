package postgresql

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

func FindStaffInTeam(id uuid.UUID, DB *Postgresql) []uuid.UUID {
	var ids []uuid.UUID
	fmt.Println(id)
	data, err := DB.postgresql.Query("select  staff from team_staff where team = $1", id)
	defer data.Close()
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

func AddManyStaffInTeam(staffs []uuid.UUID, team uuid.UUID, DB *Postgresql) {
	query := "Insert into team_staff(staff, team) values"
	log.Println(staffs)

	for i, staff := range staffs {
		query += "('" + staff.String() + "','" + team.String() + "')"
		if i == len(staffs)-1 {
			query += ";"
		} else {
			query += ","
		}
	}
	_, e := DB.postgresql.Exec(query)
	if e != nil {
		log.Fatal(e)
	}
}

func DeleteAllStaffInTeam(team uuid.UUID, DB *Postgresql) {
	stmt := "DELETE from team_staff where team=$1"
	_, err := DB.postgresql.Exec(stmt, team)
	if err != nil {
		log.Fatal(err)
	}
}
