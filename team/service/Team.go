package service

import (
	"Basic/Trainning4/redis/team/config"
	"Basic/Trainning4/redis/team/model"
	"Basic/Trainning4/redis/team/redis"
	"github.com/google/uuid"
	"log"
	"os"
)

var Cache = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10)

func GetTeamDB(database config.Database) []model.Team {
	return database.FindTeam()
}

func GetATeamDB(id uuid.UUID, database config.Database) model.TeamInter {
	var team model.TeamInter
	team = Cache.GetTeam(id)
	log.Print("Get A Team",id, team);
	if team.Name == "" {
		log.Print("In if")
		team = database.FindOneTeam(id)
		Cache.SetTeam(id, team)
	}
	return team
}

func GetNameTeamDB(id uuid.UUID, database config.Database) model.TeamInter {
	var team model.TeamInter
	team = Cache.GetTeam(id)
	log.Print("Get A Team ",id, team);
	if team.Name == "" {
		log.Print("In if")
		team = database.FindNameTeam(id)
	}
	return team
}

func CreateTeamDB(team model.Team, database config.Database) {
	database.InsertOneTeam(team)
}

func UpdateTeamDB(id uuid.UUID, idStaffs []string, database config.Database) {
	var members []uuid.UUID
	for _, idStaff := range idStaffs {
		id, err := uuid.Parse(idStaff)
		if err != nil {
			log.Fatal(err)
			return
		}
		members = append(members, id)
	}
	database.UpdateOneTeam(id, members)
	Cache.DelTeam(id)
}

func DeleteTeamDB(id uuid.UUID, database config.Database) {
	database.DeleteOneTeam(id)
	Cache.DelTeam(id)
}
