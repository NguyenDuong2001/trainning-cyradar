package main

import (
	"Basic/Trainning4/redis/team/config"
	"Basic/Trainning4/redis/team/database/mongodb"
	"Basic/Trainning4/redis/team/database/mysql"
	"Basic/Trainning4/redis/team/database/postgresql"
	"Basic/Trainning4/redis/team/model"
	"Basic/Trainning4/redis/team/redis"
	"Basic/Trainning4/redis/team/service"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	//"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	// redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).CreateClient()
	routeTeam := chi.NewRouter()
	routeTeam.Use(middleware.Logger)
	routeTeam.Get("/", getDB)
	routeTeam.Post("/", choiceDB)
	routeTeam.Route("/team", func(r chi.Router) {
		r.Use(checkDB)
		r.Get("/", GetAllTeam)
		r.Post("/", CreateTeam)
		r.Route("/detail", func(r chi.Router) {
			r.Get("/", GetATeam)
			r.Put("/", UpdateTeam)
			r.Delete("/", DeleteTeam)
		})
	})
	c := cors.New(cors.Options{
		AllowedMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:  []string{"*"},
		AllowOriginFunc: func(origin string) bool { return true },
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	handler := c.Handler(routeTeam)
	log.Printf("Starting up team on http://localhost:%s", os.Getenv("Port"))
	go run()
	http.ListenAndServe(":"+os.Getenv("Port"), handler)
}

var client = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).GetClient()

var DB config.Database = &mongodb.MongoDB{}

func run() {
	DB.NewDB()
	go func() {
		for {
			//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			//defer cancel()
			//jsonData, _ := client.BLPop(ctx, 5*time.Second, "Staff_Team").Result()
			jsonData, _ := client.BLPop(1*time.Second, "Staff_Team").Result()
			if len(jsonData) > 0 {
				var data model.DataInter
				json.Unmarshal([]byte(jsonData[1]), &data)
				switch data.Option {
				case "ReqGetTeam":
					var name []string
					log.Print(data)
					for _, idInter := range data.Data.([]interface{}) {
						idString := idInter.(string)
						id, e := uuid.Parse(idString)
						if e != nil {
							log.Fatal(e)
						}
						team := service.GetNameTeamDB(id, DB)
						log.Print(team.Name)
						name = append(name, team.Name)
					}
					log.Print("out for")
					var dataName model.DataInter
					dataName.Option = "NameTeam"
					dataName.Data = name
					jsonD, _ := json.Marshal(dataName)
					log.Print(dataName)
					//ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
					//defer cancel()
					//e := client.RPush(ctx, "ReturnNameTeam", jsonD).Err()
					e := client.RPush("ReturnNameTeam", jsonD).Err()

					if e != nil {
						log.Print(1)
						log.Fatal(e)
					}
				case "ChoiceDB":
					choice := data.Data.(string)
					if choice == "MongoDB" {
						DB = &mongodb.MongoDB{}
						log.Println("MongoDB start")
					} else if choice == "MySQL" {
						DB = &mysql.Mysql{}
						log.Println("MySql start")
					} else {
						DB = &postgresql.Postgresql{}
						log.Println("Postgresql start")
					}
					DB.NewDB()
				case "GetDB":
					var data string
					if DB == nil {
						data = ""
					} else {
						data = DB.GetName()
					}
					jsonData, e := json.Marshal(&data)
					if e != nil {
						log.Fatal(e)
					}
					client.RPush("Database", jsonData)
				default:
					break
				}
				//if data.Option == "ReqGetTeam" {
				//	var name []string
				//	for _, idInter := range data.Data.([]interface{}){
				//		idString := idInter.(string)
				//		id, e := uuid.Parse(idString)
				//		if e != nil {
				//			log.Fatal(e)
				//		}
				//		team := service.GetATeamDB(id, DB)
				//		name = append(name, team.Name)
				//	}
				//	var dataName model.DataInter
				//	dataName.Option = "NameTeam"
				//	dataName.Data = name
				//	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
				//	defer cancel()
				//	jsonD, _ := json.Marshal(dataName)
				//	e := client.RPush(ctx, "ReturnNameTeam", jsonD).Err()
				//	if e != nil {
				//		log.Println(e)
				//	}
				//}
			}
		}
	}()
}

func getDB(w http.ResponseWriter, r *http.Request) {
	database := []string{
		"MongoDB", "MySql", "Postgresql",
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(database)
}

func choiceDB(w http.ResponseWriter, r *http.Request) {
	Choice := r.URL.Query()["DB"][0]
	if Choice == "MongoDB" {
		DB = &mongodb.MongoDB{}
		log.Println("MongoDB start")
	} else if Choice == "MySQL" {
		DB = &mysql.Mysql{}
		log.Println("MySql start")
	} else if Choice == "Postgresql" {
		DB = &postgresql.Postgresql{}
		log.Println("Postgresql start")
	} else {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte("Error choice database"))
		return
	}
	DB.NewDB()
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("Choice database " + Choice))
}

func checkDB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if DB == nil {
			var data model.DataInter
			data.Option = "GetDB"
			jsonData, e := json.Marshal(&data)
			if e != nil {
				log.Fatal(e)
			}
			//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			//defer cancel()
			//client.RPush(ctx, "Team_Staff", jsonData)
			client.RPush("Team_Staff", jsonData)

			for range time.Tick(500 * time.Microsecond) {
				//jsonDB, _ := client.BLPop(ctx, 1*time.Second, "Database").Result()
				jsonDB, _ := client.BLPop(1*time.Second, "Database").Result()
				if len(jsonDB) > 0 {
					var data string
					json.Unmarshal([]byte(jsonDB[1]), &data)
					if data == "MongoDB" {
						DB = &mongodb.MongoDB{}
						log.Println("MongoDB start")
						break
					} else if data == "MySQL" {
						DB = &mysql.Mysql{}
						log.Println("MySql start")
						break
					} else if data == "Postgresql" {
						DB = &postgresql.Postgresql{}
						log.Println("Postgresql start")
						break
					} else {
						http.Error(w, "Error connect DB", 404)
						return
					}
				}
			}
		}
		DB.NewDB()
		next.ServeHTTP(w, r)
	})
}

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	memberStrings := r.Form["members"]
	if len(name) < 2 {
		w.WriteHeader(409)
		w.Write([]byte("Create failed"))
		return
	}
	newTeam := model.Team{
		ID:   uuid.New(),
		Name: name,
	}
	if len(memberStrings) >= 1 && len(memberStrings[0]) > 0 {
		for _, m := range memberStrings {
			member, err := uuid.Parse(m)
			if err != nil {
				log.Println(err)
				w.WriteHeader(409)
				w.Write([]byte("Create failed"))
				return
			}
			newTeam.Members = append(newTeam.Members, member)
		}
	}
	service.CreateTeamDB(newTeam, DB)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&newTeam)
}

func GetAllTeam(w http.ResponseWriter, r *http.Request) {
	teams := service.GetTeamDB(DB)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&teams)
}

func GetATeam(w http.ResponseWriter, r *http.Request) {
	idQuery := r.URL.Query().Get("id")
	if len(idQuery) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not found team"))
		return
	}
	id, err := uuid.Parse(idQuery)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Not found team"))
		return
	}
	team := service.GetATeamDB(id, DB)
	if team.Name == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not found team"))
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(team)
}

func UpdateTeam(w http.ResponseWriter, r *http.Request) {
	idQuery := r.URL.Query().Get("id")
	r.FormValue("members")
	memberStrings := r.Form["members"]
	if len(memberStrings) == 0 && len(memberStrings[0]) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Team not empty"))
		return
	}
	if len(idQuery) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not found team"))
		return
	}
	id, err := uuid.Parse(idQuery)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Not found team"))
		return
	}
	service.UpdateTeamDB(id, memberStrings, DB)
	w.WriteHeader(200)
	w.Write([]byte("Update successful"))
}

func DeleteTeam(w http.ResponseWriter, r *http.Request) {
	idQuery := r.URL.Query().Get("id")
	if len(idQuery) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not found team"))
		return
	}
	id, err := uuid.Parse(idQuery)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Not found team"))
		return
	}
	log.Print(id)
	service.DeleteTeamDB(id, DB)
	w.WriteHeader(200)
	w.Write([]byte("Delete team successful"))
}
