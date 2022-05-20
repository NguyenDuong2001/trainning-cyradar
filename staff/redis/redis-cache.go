package redis

import (
	"Basic/Trainning4/redis/staff/model"
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

var Client *redis.Client = redis.NewClient(&redis.Options{
	Addr:        os.Getenv("Redis_Host"),
	Password:    "",
	ReadTimeout: 7 * time.Second,
	DialTimeout: 3 * time.Minute,
	PoolSize:    1000,
	PoolTimeout: 100 * time.Minute,
})
type redisCache struct {
	host   string
	db     int
	expire time.Duration
}
type PostCache interface {
	SetStaff(key uuid.UUID, value model.Staff)
	GetStaff(key uuid.UUID) model.Staff
	DelStaff(key uuid.UUID)
	DelTeam(key uuid.UUID)
	GetClient() *redis.Client
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
	pong, err := Client.Ping(ctx).Result()
	if err != nil {

		log.Fatal(err)

	}

	// return pong if server is online

	log.Println(pong)

	return Client
}
func CreateClient()  {
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
func (cache *redisCache) SetStaff(key uuid.UUID, value model.Staff) {
	client := cache.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keyString := key.String()
	v, err := json.Marshal(value)
	err = client.Set(ctx, keyString, v, cache.expire*time.Second).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func (cache *redisCache) GetStaff(key uuid.UUID) model.Staff {
	client := cache.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keyString := key.String()
	fmt.Println(keyString)
	val, _ := client.Get(ctx, keyString).Result()
	var result model.Staff
	json.Unmarshal([]byte(val), &result)
	return result
}

func (cache *redisCache) DelStaff(key uuid.UUID) {
	client := cache.GetClient()
	keyString := key.String()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Del(ctx, keyString)
}

func (cache *redisCache) DelTeam(key uuid.UUID) {
	client := cache.GetClient()
	keyString := key.String()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Del(ctx, keyString)
}
