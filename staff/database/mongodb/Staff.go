package mongodb

import (
	"Basic/Trainning4/redis/staff/model"
	"Basic/Trainning4/redis/staff/redis"
	"context"
	"encoding/json"
	"os"
	//"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

var client = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).GetClient()

func (DB *MongoDB) FindStaff() []model.StaffInter {
	var staffs []model.StaffInter
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := DB.staffCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var staff model.Staff
		if e := cursor.Decode(&staff); err != nil {
			log.Fatal(e)
		}
		var staffInter model.StaffInter
		var teams []string
		if len(staff.Team) > 0 {
			teams = GetNameTeams(staff.Team)
		}
		staffInter.ID = staff.ID
		staffInter.Name = staff.Name
		staffInter.Age = staff.Age
		staffInter.Salary = staff.Salary
		staffInter.Team = teams
		staffs = append(staffs, staffInter)
	}
	return staffs
}

func (DB *MongoDB) FindOneStaff(id uuid.UUID) model.Staff {
	var staff model.Staff
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	DB.staffCollection.FindOne(ctx, bson.M{"id": id}).Decode(&staff)
	return staff
}

func (DB *MongoDB) FindManyStaff(ids []uuid.UUID) []model.Staff {
	var staffs []model.Staff
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := DB.staffCollection.Find(ctx, bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var staff model.Staff
		if err := cursor.Decode(&staff); err != nil {
			log.Fatal(err)
		}
		staffs = append(staffs, staff)
	}
	return staffs
}

func (DB *MongoDB) InsertOneStaff(staff model.Staff) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	newStaff := bson.M{
		"id":     staff.ID,
		"name":   staff.Name,
		"age":    staff.Age,
		"salary": staff.Salary,
	}
	_, err := DB.staffCollection.InsertOne(ctx, newStaff)
	if err != nil {
		log.Fatal("Error add staff")
	} else {
		log.Println("Create staff successfully mongodb!")
	}
}

func (M *MongoDB) UpdateOneStaff(staff model.Staff, id uuid.UUID) {
	updateStaff := bson.M{
		"name":   staff.Name,
		"age":    staff.Age,
		"salary": staff.Salary,
		"team":   staff.Team,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := M.staffCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": updateStaff})
	if err != nil {
		log.Fatal("Error update staff")
	} else {
		log.Println("Update staff successfully mongoDB!")
	}
}

func (M *MongoDB) DeleteOneStaff(id uuid.UUID) []string {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	staff := M.FindOneStaff(id)
	team := []string{}
	if len(staff.Team) == 0 {
		_, err := M.staffCollection.DeleteOne(ctx, bson.M{"id": id})
		if err != nil {
			log.Fatal("Error delete staff")
		} else {
			log.Println("Delete staff successfully mongoDB!")
		}
	} else {
		team = GetNameTeams(staff.Team)
	}
	return team
}
func (M *MongoDB) PullTeam(idTeam uuid.UUID, idStaff uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := M.staffCollection.UpdateOne(ctx, bson.M{"id": idStaff}, bson.M{"$pull": bson.M{"team": idTeam}})
	if err != nil {
		log.Fatal("Error update staff")
	} else {
		log.Println("Pull team successfully mongoDB!")
	}
}

func (M *MongoDB) PushTeam(idTeam uuid.UUID, idStaff uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := M.staffCollection.UpdateOne(ctx, bson.M{"id": idStaff}, bson.M{"$push": bson.M{"team": idTeam}})
	if err != nil {
		log.Fatal("Error update staff")
	} else {
		log.Println("Push team successfully mongoDB!")
	}
}

func GetNameTeams(teams []uuid.UUID) []string {
	var team []string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var data model.DataInter
	data.Option = "ReqGetTeam"
	data.Data = teams
	jsonData, _ := json.Marshal(data)
	err := client.RPush(ctx, "Staff_Team", jsonData).Err()
	if err != nil {
		log.Fatal(err)
	}
	for range time.Tick(500 * time.Microsecond) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// log.Print(client)
		jsonD, _ := client.BLPop(ctx, 5*time.Second, "ReturnNameTeam").Result()
		// if err != nil {
		// 	log.Fatal("err",err)
		// }
		if len(jsonD) > 0 {
			log.Print(jsonD)
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
