package service

import (
	"Basic/Trainning4/redis/staff/config"
	"Basic/Trainning4/redis/staff/model"
	"Basic/Trainning4/redis/staff/redis"
	"github.com/google/uuid"
	"os"
)

var Cache = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10)

func GetStaffDB(database config.Database) []model.StaffInter {
	return database.FindStaff()
}

func GetAStaffDB(id uuid.UUID, database config.Database) model.Staff {
	var staff model.Staff
	staff = Cache.GetStaff(id)
	if staff.Age == 0 {
		staff = database.FindOneStaff(id)
		Cache.SetStaff(id, staff)
	}
	return staff
}

func CreateStaffDB(staff model.Staff, database config.Database) {
	database.InsertOneStaff(staff)
	Cache.SetStaff(staff.ID, staff)
}

func UpdateStaffDB(staff model.Staff, id uuid.UUID, database config.Database) {
	database.UpdateOneStaff(staff, id)
	Cache.DelStaff(id)
	for _, idTeam := range staff.Team {
		Cache.DelTeam(idTeam)
	}
}

func DeleteStaffDB(id uuid.UUID, database config.Database) []string {
	team := database.DeleteOneStaff(id)
	if len(team) == 0 {
		Cache.DelStaff(id)
	}
	return team
}

func FindManyStaffDB(ids []uuid.UUID, database config.Database) []model.Staff {
	staffs := database.FindManyStaff(ids)
	return staffs
}

func PushTeam(idTeam uuid.UUID, idStaff uuid.UUID, database config.Database) {
	database.PushTeam(idTeam, idStaff)
	Cache.DelStaff(idStaff)
}

func PullTeam(idTeam uuid.UUID, idStaff uuid.UUID, database config.Database) {
	database.PullTeam(idTeam, idStaff)
	Cache.DelStaff(idStaff)
}
