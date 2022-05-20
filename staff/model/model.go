package model

import (
	"github.com/google/uuid"
)

type Staff struct {
	ID     uuid.UUID   `json:"id,omitempty" bson:"id, omitempty"`
	Name   string      `json:"name,omitempty" bson:"name, omitempty"`
	Age    int         `json:"age,omitempty" bson:"age, omitempty"`
	Salary float32     `json:"salary,omitempty" bson:"salary, omitempty"`
	Team   []uuid.UUID `json:"team,omitempty" bson:"team, omitempty"`
}

type TeamInter struct {
	ID      uuid.UUID
	Name    string
	Members []Staff
}

type StaffInter struct {
	ID     uuid.UUID
	Name   string
	Age    int
	Salary float32
	Team   []string
}

type DataInter struct {
	Option string
	Data interface{}
}
