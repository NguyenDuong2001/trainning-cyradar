package mysql

import (
	"Basic/Trainning4/redis/team/model"
	"Basic/Trainning4/redis/team/redis"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)
var client = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).GetClient()

func (DB *Mysql) FindTeam() []model.Team {
	var teams []model.Team
	data, err := DB.mysql.Query("select * from teams")
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

func (DB *Mysql) FindOneTeam(id uuid.UUID) model.TeamInter {
	var team model.Team
	data, err := DB.mysql.Query("select * from teams where id=?", id)
	defer data.Close()
	log.Println(data)
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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var ids []uuid.UUID
		for _, m := range team.Members{
			ids = append(ids,m)
		}
		var data model.DataInter
		data.Option = "GetManyStaff"
		data.Data = ids
		jsonData, _ := json.Marshal(data)
		client.RPush(ctx,"Team_Staff",jsonData)
		for range time.Tick(500 * time.Microsecond) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			jsonData,_ :=client.BLPop(ctx, 1*time.Second,"ReturnManyStaff").Result()
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

func (DB *Mysql) FindNameTeam(id uuid.UUID) model.TeamInter {
	var team model.Team
	data, err := DB.mysql.Query("select * from teams where id=?", id)
	defer data.Close()
	log.Println(data)
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

func (DB *Mysql) InsertOneTeam(team model.Team) {
	stmt, err := DB.mysql.Prepare("Insert into teams(id, name) value (?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, e := stmt.Exec(team.ID, team.Name)
	if e != nil {
		log.Fatal(e)
	}
	if len(team.Members) > 0 {
		AddManyStaffInTeam(team.Members, team.ID, DB)
	}
	log.Println("Create team successfully mysql")
}

func (DB *Mysql) UpdateOneTeam(idTeam uuid.UUID, idStaffs []uuid.UUID) {
	DeleteAllStaffInTeam(idTeam, DB)
	AddManyStaffInTeam(idStaffs, idTeam, DB)
	log.Println("Update team successfully mysql")
}

func (DB *Mysql) DeleteOneTeam(id uuid.UUID) {
	DeleteAllStaffInTeam(id, DB)
	stmt, _ := DB.mysql.Prepare("DELETE from teams where id= ?")
	_, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Delete team successfully mysql")
}
