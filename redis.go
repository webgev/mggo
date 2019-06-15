package mggo

import (
	"fmt"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func init() {
	InitCallback(func() {
		section, err := config.GetSection("redis")
		if err != nil {
			return
		}
		v, err := section.GetKey("address")
		if err != nil {
			return
		}
		var password string
		pass, err := section.GetKey("password")
		if err == nil {
			password = pass.String()
		}
		redisClient = redis.NewClient(&redis.Options{
			Addr:     v.String(),
			Password: password,
			DB:       0,
		})

		pong, err := redisClient.Ping().Result()
		fmt.Println(pong, err)
	})
}

//Redis get redis client
func Redis() *redis.Client {
	return redisClient
}
