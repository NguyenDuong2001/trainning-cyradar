package redis

import (
	"Basic/Trainning4/redis/team/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	// "net"
)

type redisCache struct {
	host   string
	db     int
	expire time.Duration
}

// net.JoinHostPort(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
var Client *redis.Client = redis.NewClient(&redis.Options{
	// Addr:        os.Getenv("Redis_Host"),
	Addr: os.Getenv("REDIS_URL"),
	// Addr:        net.JoinHostPort(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	Password:    os.Getenv("REDIS_PASSWORD"),
	DB:          0,
	PoolSize:    1000,
	PoolTimeout: 10 * time.Second,
})

type PostCache interface {
	SetTeam(key uuid.UUID, value model.TeamInter)
	GetTeam(key uuid.UUID) model.TeamInter
	DelTeam(key uuid.UUID)
	GetClient() *redis.Client
	CreateClient()
}

func NewRedisCache(host string, db int, expire time.Duration) PostCache {
	return &redisCache{
		host:   host,
		db:     db,
		expire: expire,
	}
}

func (cache *redisCache) GetClient() *redis.Client {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	pong, err := Client.Ping().Result()

	fmt.Println(pong)
	if err != nil {
		fmt.Println("Error get client team")
		log.Fatal(err)
	}
	log.Println(pong)
	return Client

}

func (cache *redisCache) CreateClient() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	Client = redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    os.Getenv("REDIS_PASSWORD"),
		ReadTimeout: 7 * time.Second,
		DialTimeout: 3 * time.Second,
		PoolSize:    1000,
		PoolTimeout: 10 * time.Second,
	})
}

func (cache *redisCache) DelStaff(key uuid.UUID) {
	fmt.Println("delstaff 1")
	// client := cache.GetClient()
	client := Client
	fmt.Println("delstaff 2", client != nil)

	keyString := key.String()
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//fmt.Println("delstaff 3")
	//
	//err := client.Del(ctx, keyString).Err()
	err := client.Del(keyString).Err()

	fmt.Println("delstaff 4", err)
}

func (cache *redisCache) SetTeam(key uuid.UUID, value model.TeamInter) {
	fmt.Println("setteam 1")

	client := cache.GetClient()
	fmt.Println("setteam 2", client != nil)

	fmt.Println("setteam 3")

	keyString := key.String()
	fmt.Println("KEY", keyString)
	v, _ := json.Marshal(value)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//err := client.Set(ctx, keyString, v, cache.expire*time.Second).Err()
	err := client.Set(keyString, v, 10*time.Minute).Err()
	fmt.Println("setstaff 4", err)

}

func (cache *redisCache) GetTeam(key uuid.UUID) model.TeamInter {
	fmt.Println("getteam 1")

	// client := cache.GetClient()
	client := Client

	fmt.Println("getteam 2", client != nil)

	fmt.Println("getteam 3")

	keyString := key.String()
	fmt.Println(keyString)
	//ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	//defer cancel()
	//// e := client.Set(ctx, "KEY", "Hi", cache.expire*time.Second).Err()
	//// fmt.Println("set in get", e)
	//val, err := client.Get(ctx, keyString).Result()
	val, err := client.Get(keyString).Result()
	fmt.Println("getTeam 4", keyString, err)
	fmt.Println("getTeam val", val)
	var result model.TeamInter
	if err != redis.Nil {
		if err == nil {
			err = json.Unmarshal([]byte(val), &result)
		}
	}
	fmt.Println("getTeam 5", err)
	fmt.Println("get Team 6")
	return result
}

func (cache *redisCache) DelTeam(key uuid.UUID) {
	fmt.Println("DelTeam 1")
	client := cache.GetClient()
	fmt.Println("DelTeam 2", client != nil)

	keyString := key.String()
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//fmt.Println("DelTeam 3")
	//
	//err := client.Del(ctx, keyString).Err()
	err := client.Del(keyString).Err()
	fmt.Println("delstaff 4", err)
}
