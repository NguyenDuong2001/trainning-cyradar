package main

import (
	"Basic/Trainning4/redis/staff/config"
	"Basic/Trainning4/redis/staff/database/mongodb"
	"Basic/Trainning4/redis/staff/database/mysql"
	"Basic/Trainning4/redis/staff/database/postgresql"
	"Basic/Trainning4/redis/staff/model"
	"Basic/Trainning4/redis/staff/redis"
	"Basic/Trainning4/redis/staff/service"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	//"github.com/go-redis/redis/v8"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var DB config.Database = &mongodb.MongoDB{}

func main() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	//redis.CreateClient()
	routeStaff := chi.NewRouter()
	routeStaff.Use(middleware.Logger)
	routeStaff.Get("/", getDB)
	routeStaff.Post("/", choiceDB)
	routeStaff.Route("/staff", func(r chi.Router) {
		//r.Use(CorsMiddleware)
		r.Use(checkDB)
		r.Get("/", GetAllStaff)
		r.Get("/editTeam", GetAllStaffToEdit)
		r.Get("/Many", GetManyStaff)
		r.Post("/", CreateStaff)
		r.Put("/Pull", PullStaff)
		r.Put("/Push", PushStaff)
		r.Route("/detail", func(r chi.Router) {
			r.Get("/", GetAStaff)
			r.Put("/", UpdateStaff)
			r.Delete("/", DeleteStaff)
		})
	})
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Enable Debugging for testing, consider disabling in production
		AllowOriginFunc: func(origin string) bool { return true },
		Debug:           true,
		AllowedHeaders:  []string{"accept", "authorization", "content-type"},
	})
	handler := c.Handler(routeStaff)

	// Insert the middleware
	handler = c.Handler(handler)
	log.Printf("Starting up staff on http://localhost:%s", os.Getenv("Port"))
	go run()
	http.ListenAndServe(":"+os.Getenv("Port"), handler)
}

var client = redis.NewRedisCache(os.Getenv("Redis_Host"), 0, 10).GetClient()

func run() {
	DB.NewDB()
	go func() {
		for {
			//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			//defer cancel()
			//jsonData, _ := client.BLPop(ctx, 1*time.Second, "Team_Staff").Result()
			jsonData, _ := client.BLPop(1*time.Second, "Team_Staff").Result()
			if len(jsonData) > 0 {
				var data model.DataInter
				json.Unmarshal([]byte(jsonData[1]), &data)
				switch data.Option {
				case "GetManyStaff":
					var dataID []uuid.UUID
					dataInter := data.Data.([]interface{})
					for _, idIner := range dataInter {
						idString := idIner.(string)
						id, e := uuid.Parse(idString)
						if e != nil {
							log.Fatal(e)
						}
						dataID = append(dataID, id)
					}
					staffs := service.FindManyStaffDB(dataID, DB)
					log.Print(staffs)
					d, _ := json.Marshal(staffs)
					//e := client.RPush(ctx, "ReturnManyStaff", d).Err()
					e := client.RPush("ReturnManyStaff", d).Err()
					if e != nil {
						log.Fatal(e)
					}
				case "PushStaff":
					data := data.Data.(map[string]interface{})
					idTeam, e := uuid.Parse(data["idTeam"].(string))
					idStaff, e := uuid.Parse(data["idStaff"].(string))
					if e != nil {
						log.Fatal(e)
					}
					service.PushTeam(idTeam, idStaff, DB)
				case "PullStaff":
					data := data.Data.(map[string]interface{})
					idTeam, e := uuid.Parse(data["idTeam"].(string))
					idStaff, e := uuid.Parse(data["idStaff"].(string))
					if e != nil {
						log.Fatal(e)
					}
					service.PullTeam(idTeam, idStaff, DB)
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
					//client.RPush(ctx, "Database", jsonData)
					client.RPush("Database", jsonData)
				default:
					break
				}
				//if data.Option == "GetManyStaff" {
				//	staffs := service.FindManyStaffDB(data.Data, DB)
				//	d,_ := json.Marshal(staffs)
				//	client.RPush(ctx,"ReturnManyStaff", d)
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
	choice := r.URL.Query()["DB"][0]
	if choice == "MongoDB" {
		DB = &mongodb.MongoDB{}
		log.Println("MongoDB start")
	} else if choice == "MySQL" {
		DB = &mysql.Mysql{}
		log.Println("MySql start")
	} else if choice == "Postgresql" {
		DB = &postgresql.Postgresql{}
		log.Println("Postgresql start")
	} else {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte("Error choice database"))
		return
	}
	//params := "?DB=" + url.QueryEscape(choice)
	//path := fmt.Sprintf(os.Getenv("PortTeam")+"%s", params)
	//_, err := http.Post(path, "application/json", bytes.NewBuffer([]byte("1")))
	//if err != nil {
	//	log.Fatal(err)
	//}
	var data model.DataInter
	data.Option = "ChoiceDB"
	data.Data = choice
	jsonData, e := json.Marshal(&data)
	if e != nil {
		log.Fatal(e)
	}
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//client.RPush(ctx, "Staff_Team", jsonData)
	client.RPush("Staff_Team", jsonData)
	DB.NewDB()
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("Choice database " + choice))
}

func checkDB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if DB == nil {
			var data model.DataInter
			data.Option = "GetDB"
			jsonData, e := json.Marshal(&data)
			if e != nil {
				log.Fatal(e)
			}
			//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			//defer cancel()
			//client.RPush(ctx, "Staff_Team", jsonData)
			client.RPush("Staff_Team", jsonData)
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

func GetAllStaff(w http.ResponseWriter, r *http.Request) {
	staffs := service.GetStaffDB(DB, false)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&staffs)
}

func GetAllStaffToEdit(w http.ResponseWriter, r *http.Request) {
	staffs := service.GetStaffDB(DB, true)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&staffs)
}

func GetAStaff(w http.ResponseWriter, r *http.Request) {
	idQuery := r.URL.Query().Get("id")
	if len(idQuery) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not found staff"))
		return
	}
	id, err := uuid.Parse(idQuery)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Not found staff"))
		return
	}
	staff := service.GetAStaffDB(id, DB)
	if staff.Name == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not found staff"))
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(staff)
}

func CreateStaff(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if len(name) < 2 {
		w.WriteHeader(409)
		w.Write([]byte("Create failed"))
		return
	}
	ageString := r.FormValue("age")
	age, err := strconv.Atoi(ageString)
	if err != nil {
		log.Println(err)
		w.WriteHeader(409)
		w.Write([]byte("Create failed"))
		return
	}
	salaryString := r.FormValue("salary")
	salary, e := strconv.ParseFloat(salaryString, 32)
	if e != nil {
		log.Println(e)
		w.WriteHeader(409)
		w.Write([]byte("Create failed"))
		return
	}
	newStaff := model.Staff{
		ID:     uuid.New(),
		Name:   name,
		Age:    age,
		Salary: float32(math.Ceil(salary*100) / 100),
	}
	service.CreateStaffDB(newStaff, DB)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&newStaff)
}

func UpdateStaff(w http.ResponseWriter, r *http.Request) {
	idQuery := r.URL.Query().Get("id")
	if len(idQuery) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not found staff"))
		return
	}
	id, errr := uuid.Parse(idQuery)
	if errr != nil {
		w.WriteHeader(404)
		w.Write([]byte("Not found staff"))
		return
	}
	staff := service.GetAStaffDB(id, DB)
	if staff.Name == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not found staff"))
		return
	}
	name := r.FormValue("name")
	if len(name) == 0 {
		name = staff.Name
	} else if len(name) < 2 {
		w.WriteHeader(409)
		w.Write([]byte("Update failed"))
		return
	}
	ageString := r.FormValue("age")
	var age int
	if len(ageString) == 0 {
		age = staff.Age
	} else {
		ageParse, err := strconv.Atoi(ageString)
		if err != nil {
			log.Println(err)
			w.WriteHeader(409)
			w.Write([]byte("Update failed"))
			return
		}
		age = ageParse
	}
	var salary float32
	salaryString := r.FormValue("salary")
	if len(salaryString) == 0 {
		salary = staff.Salary
	} else {
		salaryParse, e := strconv.ParseFloat(salaryString, 32)
		if e != nil {
			log.Println(e)
			w.WriteHeader(409)
			w.Write([]byte("Update failed"))
			return
		}
		salary = float32(math.Round(salaryParse*100) / 100)
	}
	updateStaff := model.Staff{
		uuid.New(),
		name,
		age,
		salary,
		staff.Team,
	}
	service.UpdateStaffDB(updateStaff, staff.ID, DB)
	w.WriteHeader(200)
	w.Write([]byte("Update successful"))
}

func DeleteStaff(w http.ResponseWriter, r *http.Request) {
	idQuery := r.URL.Query().Get("id")
	if len(idQuery) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not found staff"))
		return
	}
	id, err := uuid.Parse(idQuery)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Not found staff"))
		return
	}
	teamCheck := service.DeleteStaffDB(id, DB)
	if len(teamCheck) == 0 {
		w.WriteHeader(200)
		w.Write([]byte("Delete successful"))
		return
	}
	w.WriteHeader(409)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(teamCheck)
	return
}

func GetManyStaff(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var ids []uuid.UUID
	idArray := r.Form["ids"]
	for _, idString := range idArray {
		id, e := uuid.Parse(idString)
		if e != nil {
			log.Fatal(e)
		}
		ids = append(ids, id)
	}
	staffs := service.FindManyStaffDB(ids, DB)
	json.NewEncoder(w).Encode(&staffs)
}

func PullStaff(w http.ResponseWriter, r *http.Request) {
	ids := make(map[string]uuid.UUID)
	json.NewDecoder(r.Body).Decode(&ids)
	service.PullTeam(ids["idTeam"], ids["idStaff"], DB)
}

func PushStaff(w http.ResponseWriter, r *http.Request) {
	ids := make(map[string]uuid.UUID)
	json.NewDecoder(r.Body).Decode(&ids)
	service.PushTeam(ids["idTeam"], ids["idStaff"], DB)
}
