package config
import (
	"Basic/Trainning4/redis/staff/model"
	"github.com/google/uuid"
)

type Database interface {
	NewDB()
	GetName() string
	FindStaff(bool) []model.StaffInter
	FindOneStaff(uuid.UUID) model.Staff
	InsertOneStaff(model.Staff)
	UpdateOneStaff(model.Staff, uuid.UUID)
	DeleteOneStaff(uuid.UUID) []string
	FindManyStaff([]uuid.UUID) []model.Staff
	PullTeam(uuid.UUID, uuid.UUID)
	PushTeam(uuid.UUID, uuid.UUID)
}
