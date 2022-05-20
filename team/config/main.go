package config

import (
	"Basic/Trainning4/redis/team/model"
	"github.com/google/uuid"
)

type Database interface {
	NewDB()
	GetName() string
	FindOneTeam(uuid.UUID) model.TeamInter
	FindNameTeam(uuid.UUID) model.TeamInter
	FindTeam() []model.Team
	InsertOneTeam(model.Team)
	UpdateOneTeam(uuid.UUID, []uuid.UUID)
	DeleteOneTeam(uuid.UUID)
}
