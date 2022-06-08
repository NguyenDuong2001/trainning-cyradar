package redis

import (
	"Basic/Trainning4/redis/staff/model"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	// "net"
	"time"
)

var Client *redis.Client = redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_URL"),
	// Addr:        net.JoinHostPort(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	// Password:    os.Getenv("REDIS_PASSWORD"),
	Password:    os.Getenv("REDIS_PASSWORD"),
	DB:          0,
	PoolSize:    1000,
	PoolTimeout: 10 * time.Second,
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
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//pong, err := Client.Ping(ctx).Result()
	pong, err := Client.Ping().Result()
	if err != nil {

		log.Fatal(err)

	}

	// return pong if server is online

	log.Println(pong)

	return Client
}
func CreateClient() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	Client = redis.NewClient(&redis.Options{
		Addr:        os.Getenv("Redis_Host"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		ReadTimeout: 7 * time.Second,
		DialTimeout: 3 * time.Minute,
		PoolSize:    1000,
		PoolTimeout: 100 * time.Minute,
	})
}
func (cache *redisCache) SetStaff(key uuid.UUID, value model.Staff) {
	client := cache.GetClient()
	keyString := key.String()
	v, err := json.Marshal(value)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//err = client.Set(ctx, keyString, v, cache.expire*time.Second).Err()
	err = client.Set(keyString, v, 10*time.Minute).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func (cache *redisCache) GetStaff(key uuid.UUID) model.Staff {
	client := cache.GetClient()
	keyString := key.String()
	fmt.Println(keyString)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//val, _ := client.Get(ctx, keyString).Result()
	val, _ := client.Get(keyString).Result()
	var result model.Staff
	json.Unmarshal([]byte(val), &result)
	return result
}

func (cache *redisCache) DelStaff(key uuid.UUID) {
	client := cache.GetClient()
	keyString := key.String()
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//client.Del(ctx, keyString)
	client.Del(keyString)
}

func (cache *redisCache) DelTeam(key uuid.UUID) {
	client := cache.GetClient()
	keyString := key.String()
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//client.Del(ctx, keyString)
	client.Del(keyString)
}
