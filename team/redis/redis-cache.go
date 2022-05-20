package redis

import (
	"Basic/Trainning4/redis/team/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type redisCache struct {
	host   string
	db     int
	expire time.Duration
}

var Client *redis.Client = redis.NewClient(&redis.Options{
	Addr:        os.Getenv("Redis_Host"),
	Password:    "",
	ReadTimeout: 7 * time.Second,
	DialTimeout: 3 * time.Minute,
	PoolSize:    1000,
	PoolTimeout: 100 * time.Minute,
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Print(Client)
	pong, err := Client.Ping(ctx).Result()
	if err != nil {

		log.Fatal(err)

	}

	// return pong if server is online

	log.Println(pong)
	return Client

}

func (cache *redisCache) CreateClient()  {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	Client = redis.NewClient(&redis.Options{
		Addr:        os.Getenv("Redis_Host"),
		Password:    "",
		ReadTimeout: 7 * time.Second,
		DialTimeout: 3 * time.Minute,
		PoolSize:    1000,
		PoolTimeout: 100 * time.Minute,
	})
}

func (cache *redisCache) DelStaff(key uuid.UUID) {
	client := cache.GetClient()
	keyString := key.String()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Del(ctx, keyString)
}

func (cache *redisCache) SetTeam(key uuid.UUID, value model.TeamInter) {
	client := cache.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keyString := key.String()
	v, _ := json.Marshal(value)
	client.Set(ctx, keyString, v, cache.expire*time.Second).Err()
}

func (cache *redisCache) GetTeam(key uuid.UUID) model.TeamInter {
	client := cache.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keyString := key.String()
	fmt.Println(keyString)
	val, _ := client.Get(ctx, keyString).Result()
	var result model.TeamInter
	json.Unmarshal([]byte(val), &result)
	return result
}

func (cache *redisCache) DelTeam(key uuid.UUID) {
	client := cache.GetClient()
	keyString := key.String()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Del(ctx, keyString)
}
