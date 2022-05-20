package mongodb

import (
	"Basic/Trainning4/redis/team/model"
	"Basic/Trainning4/redis/team/redis"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"time"
)

var client = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).GetClient()
func (DB *MongoDB) FindTeam() []model.Team {
	var teams []model.Team
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := DB.teamCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Error get team")
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var team model.Team
		cursor.Decode(&team)
		teams = append(teams, team)
	}
	return teams
}

func (DB *MongoDB) FindOneTeam(id uuid.UUID) model.TeamInter {
	var team model.Team
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	DB.teamCollection.FindOne(ctx, bson.M{"id": id}).Decode(&team)
	if team.Name == "" {
		return model.TeamInter{}
	}
	teamInter := model.TeamInter{
		ID:team.ID,
		Name: team.Name,
	}
	log.Println("Find", team)
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
		e := client.RPush(ctx,"Team_Staff",jsonData).Err()
		if e != nil {
			log.Fatal(e)
		}
		for range time.Tick(500 * time.Microsecond) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			jsonData,_ :=client.BLPop(ctx, 1*time.Second,"ReturnManyStaff").Result()
			if len(jsonData)>0{
				log.Println("Return Staff",jsonData)
				var data []model.Staff
				json.Unmarshal([]byte(jsonData[1]),&data)
				teamInter.Members = data
				return teamInter
			}
		}
	}
	return teamInter
}

func (DB *MongoDB) FindNameTeam(id uuid.UUID) model.TeamInter {
	var team model.Team
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	DB.teamCollection.FindOne(ctx, bson.M{"id": id}).Decode(&team)
	if team.Name == "" {
		return model.TeamInter{}
	}
	teamInter := model.TeamInter{
		ID:team.ID,
		Name: team.Name,
	}
	return teamInter
}

func (DB *MongoDB) InsertOneTeam(team model.Team) {
	log.Println(team)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	newTeam := bson.M{
		"id":      team.ID,
		"name":    team.Name,
		"members": team.Members,
	}
	_, err := DB.teamCollection.InsertOne(ctx, newTeam)

	for _, m := range team.Members {
		//path := fmt.Sprintf(os.Getenv("PortStaff") + "staff/Push")
		ids := map[string]uuid.UUID{
			"idTeam":  team.ID,
			"idStaff": m,
		}
		var data model.DataInter
		data.Option = "PushStaff"
		data.Data = ids
		jsonData, e := json.Marshal(data)
		if e != nil {
			log.Fatal(e)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client.RPush(ctx, "Team_Staff", jsonData)

		//req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(body))
		//req.Header.Set("Content-Type", "application/json; charset=utf-8")
		//client := &http.Client{}
		//client.Do(req)
	}
	if err != nil {
		log.Fatal("Error add staff")
	} else {
		log.Println("Create team successfully mongoDB")
	}
}

func (DB *MongoDB) UpdateOneTeam(id uuid.UUID, newMembers []uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var team model.Team
	DB.teamCollection.FindOne(ctx, bson.M{"id": id}).Decode(&team)
	updateTeam := bson.M{
		"name":    team.Name,
		"members": newMembers,
	}
	_, err := DB.teamCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": updateTeam})
	if err != nil {
		log.Fatal("Error update team")
	}
	if len(team.Members) > 0 {
		for _, Tm := range team.Members {
		CheckPull:
			for _, m := range newMembers {
				if Tm == m {
					continue CheckPull
				}
			}
			ids := map[string]uuid.UUID{
				"idTeam":  id,
				"idStaff": Tm,
			}
			var data model.DataInter
			data.Option ="PullStaff"
			data.Data = ids
			jsonData, e := json.Marshal(data)
			if e != nil {
				log.Fatal(e)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			client.RPush(ctx, "Team_Staff", jsonData)
			//path := fmt.Sprintf(os.Getenv("PortStaff") + "staff/Pull")
			//body, e := json.Marshal(ids)
			//if e != nil {
			//	log.Fatal(e)
			//}
			//req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(body))
			//req.Header.Set("Content-Type", "application/json; charset=utf-8")
			//client := &http.Client{}
			//client.Do(req)
		}
		for _, m := range newMembers {
		CheckPush:
			for _, Tm := range team.Members {
				if m == Tm {
					continue CheckPush
				}
			}
			ids := map[string]uuid.UUID{
				"idTeam":  team.ID,
				"idStaff": m,
			}
			var data model.DataInter
			data.Option ="PushStaff"
			data.Data = ids
			jsonData, e := json.Marshal(data)
			if e != nil {
				log.Fatal(e)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			client.RPush(ctx, "Team_Staff", jsonData)
			//path := fmt.Sprintf(os.Getenv("PortStaff") + "staff/Push")
			//body, e := json.Marshal(ids)
			//if e != nil {
			//	log.Fatal(e)
			//}
			//req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(body))
			//req.Header.Set("Content-Type", "application/json; charset=utf-8")
			//client := &http.Client{}
			//client.Do(req)
		}
	} else {
		for _, m := range newMembers {
			ids := map[string]uuid.UUID{
				"idTeam":  team.ID,
				"idStaff": m,
			}
			var data model.DataInter
			data.Option ="PushStaff"
			data.Data = ids
			jsonData, e := json.Marshal(data)
			if e != nil {
				log.Fatal(e)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			client.RPush(ctx, "Team_Staff", jsonData)
			//path := fmt.Sprintf(os.Getenv("PortStaff") + "staff/Push")
			//body, e := json.Marshal(ids)
			//if e != nil {
			//	log.Fatal(e)
			//}
			//req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(body))
			//req.Header.Set("Content-Type", "application/json; charset=utf-8")
			//client := &http.Client{}
			//client.Do(req)
		}
	}
	log.Println("Update team successfully mongoDB")
}

func (DB *MongoDB) DeleteOneTeam(id uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var team model.Team
	DB.teamCollection.FindOne(ctx, bson.M{"id": id}).Decode(&team)
	if len(team.Members) > 0 {
		for _, m := range team.Members {
			ids := map[string]uuid.UUID{
				"idTeam":  team.ID,
				"idStaff": m,
			}
			var data model.DataInter
			data.Option ="PullStaff"
			data.Data = ids
			jsonData, e := json.Marshal(data)
			if e != nil {
				log.Fatal(e)
			}
			client.RPush(ctx, "Team_Staff", jsonData)
			//path := fmt.Sprintf(os.Getenv("PortStaff") + "staff/Pull")
			//body, e := json.Marshal(ids)
			//if e != nil {
			//	log.Fatal(e)
			//}
			//req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(body))
			//req.Header.Set("Content-Type", "application/json; charset=utf-8")
			//client := &http.client{}
			//client.Do(req)
		}
	}
	_, err := DB.teamCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Delete team successfully mongoDB")
		return
	}
}
