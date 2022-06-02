package postgresql

import (
	"Basic/Trainning4/redis/team/model"
	"Basic/Trainning4/redis/team/redis"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)
var client = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).GetClient()

func (DB *Postgresql) FindTeam() []model.Team {
	var teams []model.Team
	data, err := DB.postgresql.Query("select * from teams")
	defer data.Close()
	if err != nil {
		log.Fatal(err)
	}
	for data.Next() {
		var team model.Team
		e := data.Scan(&team.ID, &team.Name)
		if e != nil {
			log.Fatal(e)
		}
		team.Members = FindStaffInTeam(team.ID, DB)
		teams = append(teams, team)
	}
	return teams
}

func (DB *Postgresql) FindOneTeam(id uuid.UUID) model.TeamInter {
	var team model.Team
	data, err := DB.postgresql.Query("select * from teams where id=$1", id)
	defer data.Close()
	if err != nil {
		log.Fatal(err)
	}
	for data.Next() {
		e := data.Scan(&team.ID, &team.Name)
		if e != nil {
			log.Fatal(e)
		}
		team.Members = FindStaffInTeam(team.ID, DB)
	}
	var client = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).GetClient()
	teamInter := model.TeamInter{
		ID: team.ID,
		Name: team.Name,
	}
	if len(team.Members) > 0 {
		//	//param := "?"
		//	//for _, m := range team.Members {
		//	//	param += "ids=" + m.String() + "&"
		//	//}
		//	//path := fmt.Sprintf(os.Getenv("PortStaff")+"staff/Many%s", param)
		//	//resp, err := http.Get(path)
		//	//if err != nil {
		//	//	log.Println(err)
		//	//}
		//	//defer resp.Body.Close()
		//	//json.NewDecoder(resp.Body).Decode(&members)
		var data model.DataInter
		data.Option = "GetManyStaff"
		data.Data = team.Members
		jsonData, _ := json.Marshal(data)
		//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//defer cancel()
		//client.RPush(ctx,"Team_Staff",jsonData)
		client.RPush("Team_Staff",jsonData)

		for range time.Tick(500 * time.Microsecond) {
			//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			//defer cancel()
			//jsonData,_ :=client.BLPop(ctx, 1*time.Second,"ReturnManyStaff").Result()
			jsonData,_ :=client.BLPop(1*time.Second,"ReturnManyStaff").Result()
			if len(jsonData)>0{
				var data []model.Staff
				json.Unmarshal([]byte(jsonData[1]),&data)
				teamInter.Members = data
				return teamInter
			}
		}
	}
	return teamInter
}

func (DB *Postgresql) FindNameTeam(id uuid.UUID) model.TeamInter {
	var team model.Team
	data, err := DB.postgresql.Query("select * from teams where id=$1", id)
	defer data.Close()
	if err != nil {
		log.Fatal(err)
	}
	for data.Next() {
		e := data.Scan(&team.ID, &team.Name)
		if e != nil {
			log.Fatal(e)
		}
		team.Members = FindStaffInTeam(team.ID, DB)
	}
	teamInter := model.TeamInter{
		ID: team.ID,
		Name: team.Name,
	}
	return teamInter
}

func (DB *Postgresql) InsertOneTeam(team model.Team) {
	stmt := "Insert into teams(id, name) values ($1,$2)"
	_, err := DB.postgresql.Exec(stmt, team.ID, team.Name)
	if err != nil {
		log.Fatal(err)
	}
	if len(team.Members) > 0 {
		AddManyStaffInTeam(team.Members, team.ID, DB)
	}
	log.Println("Create team successfully postgresql")
}

func (DB *Postgresql) UpdateOneTeam(idTeam uuid.UUID, idStaffs []uuid.UUID) {
	DeleteAllStaffInTeam(idTeam, DB)
	AddManyStaffInTeam(idStaffs, idTeam, DB)
	log.Println("Update team successfully postgresql")
}

func (DB *Postgresql) DeleteOneTeam(id uuid.UUID) {
	DeleteAllStaffInTeam(id, DB)
	stmt := "Delete from teams where id = $1"
	_, err := DB.postgresql.Exec(stmt, id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Delete team successfully postgresql")
}
