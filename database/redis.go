package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisDB struct {
	Redis *redis.Client
}

var (
	RDB *redisDB
	CTX = context.Background()
)

func InitRedis() (err error) {
	RDB = &redisDB{}
	RDB.Redis = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err = RDB.Redis.Ping(CTX).Result()
	if err != nil {
		log.Println("Error connecting to Redis ", err)
		return err
	}
	log.Println("Connected to redis")
	return nil
}

func (db *redisDB) SetUserAndToken(id uint, token string) (err error) {
	strID := fmt.Sprint(id)
	err = db.Redis.Set(CTX, fmt.Sprintf("%s:%s", strID, token), 0, time.Duration(time.Hour*72)).Err()
	return err
}

func (db *redisDB) GetUserAndToken(id uint, token string) (string, error) {
	strID := fmt.Sprint(id)
	val, err := db.Redis.Get(CTX, fmt.Sprintf("%s:%s", strID, token)).Result()
	return val, err
}

func (db *redisDB) DeleteToken(id uint, token string) (int64, error) {
	strID := fmt.Sprint(id)
	val, err := db.Redis.Del(CTX, fmt.Sprintf("%s:%s", strID, token)).Result()
	return val, err

}
