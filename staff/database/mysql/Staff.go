package mysql

import (
	"Basic/Trainning4/redis/staff/model"
	"Basic/Trainning4/redis/staff/redis"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

var client = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).GetClient()

func (DB *Mysql) FindStaff() []model.StaffInter {
	var staffs []model.StaffInter
	data, e := DB.mysql.Query("select * from staffs")
	defer data.Close()
	if e != nil {
		log.Fatal(e)
	}
	for data.Next() {
		var staff model.Staff
		err := data.Scan(&staff.ID, &staff.Name, &staff.Age, &staff.Salary)
		if err != nil {
			log.Fatal(err)
		}
		data, err := DB.mysql.Query("select  team from team_staff where staff = ?", staff.ID)
		defer data.Close()
		if err != nil {
			log.Fatal(err)
		}
		var teams []string
		var staffInter model.StaffInter
		staffInter.ID = staff.ID
		staffInter.Name = staff.Name
		staffInter.Age = staff.Age
		staffInter.Salary = staff.Salary
		var IDTeams []uuid.UUID
		if !data.Next() {
			staffInter.Team = teams
			if err != nil {
				log.Fatal(err)
			}
			break
		} else {
			var idTeam uuid.UUID
			err := data.Scan(&idTeam)
			if err != nil {
				log.Fatal(err)
			}
			IDTeams = append(IDTeams, idTeam)
		}
		for data.Next() {
			var idTeam uuid.UUID
			err := data.Scan(&idTeam)
			if err != nil {
				log.Fatal(err)
			}
			IDTeams = append(IDTeams, idTeam)
		}
		if len(IDTeams) > 0 {
			teams = GetNameTeams(IDTeams)
		}
		staffInter.Team = teams
		staffs = append(staffs, staffInter)
	}
	return staffs
}
func (DB *Mysql) FindOneStaff(id uuid.UUID) model.Staff {
	var staff model.Staff
	data, err := DB.mysql.Query("select  * from staffs where id = ?", id)
	defer data.Close()
	if err != nil {
		log.Fatal(err)
	}
	for data.Next() {
		e := data.Scan(&staff.ID, &staff.Name, &staff.Age, &staff.Salary)
		if e != nil {
			log.Fatal(e)
		}
	}
	return staff
}

func (DB *Mysql) FindManyStaff(ids []uuid.UUID) []model.Staff {
	var staffs []model.Staff
	if len(ids) > 0 {
		query := "select * from staffs where id in ("
		for i, id := range ids {
			query += "'"
			query += id.String()
			query += "'"
			if i < len(ids)-1 {
				query += ","
			}
		}
		query += ")"
		data, err := DB.mysql.Query(query)
		defer data.Close()
		if err != nil {
			log.Fatal(err)
		}
		for data.Next() {
			var staff model.Staff
			e := data.Scan(&staff.ID, &staff.Name, &staff.Age, &staff.Salary)
			if e != nil {
				log.Fatal(e)
			}
			staffs = append(staffs, staff)
		}
	}
	return staffs
}

func (DB *Mysql) InsertOneStaff(staff model.Staff) {
	stmt, err := DB.mysql.Prepare("Insert into staffs(id, name, age, salary) value (?,?,?,?)")
	_, err = stmt.Exec(staff.ID, staff.Name, staff.Age, staff.Salary)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Create staff successfully mysql !")
}

func (DB *Mysql) UpdateOneStaff(staff model.Staff, id uuid.UUID) {
	stmt, err := DB.mysql.Prepare("update staffs set name=?, age = ?, salary =? where id=?")
	_, err = stmt.Exec(staff.Name, staff.Age, staff.Salary, id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Update staff successfully mysql!")
}

func (DB *Mysql) DeleteOneStaff(id uuid.UUID) []string {
	var teams []string
	data, err := DB.mysql.Query("select  team from team_staff where staff = ?", id)
	defer data.Close()
	if err != nil {
		log.Fatal(err)
	}
	var IDTeams []uuid.UUID
	if !data.Next() {
		stmt, err := DB.mysql.Prepare("DELETE from staffs where id= ?")
		_, err = stmt.Exec(id)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Delete staff successfully postgresql")
		return teams
	} else {
		var idTeam uuid.UUID
		err := data.Scan(&idTeam)
		if err != nil {
			log.Fatal(err)
		}
		IDTeams = append(IDTeams, idTeam)
	}
	for data.Next() {
		var idTeam uuid.UUID
		err := data.Scan(&idTeam)
		if err != nil {
			log.Fatal(err)
		}
		IDTeams = append(IDTeams, idTeam)
	}
	teams = GetNameTeams(IDTeams)
	return teams
}

func GetNameTeams(teams []uuid.UUID) []string {
	var team []string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var data model.DataInter
	data.Option = "ReqGetTeam"
	data.Data = teams
	jsonData, _ := json.Marshal(data)
	client.RPush(ctx, "Staff_Team", jsonData)
	for range time.Tick(500 * time.Microsecond) {
		jsonD, _ := client.BLPop(ctx, 5*time.Second, "ReturnNameTeam").Result()
		if len(jsonD) > 0 {
			var data model.DataInter
			json.Unmarshal([]byte(jsonD[1]), &data)
			array := data.Data.([]interface{})
			for _, t := range array {
				team = append(team, t.(string))
			}
			return team
		}
	}
	return team
}

func (M *Mysql) PullTeam(idTeam uuid.UUID, idStaff uuid.UUID) {}

func (M *Mysql) PushTeam(idTeam uuid.UUID, idStaff uuid.UUID) {}
