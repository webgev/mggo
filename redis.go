package mggo

import (
    "github.com/go-redis/redis"
    "fmt"
)

var redicClient *redis.Client

func init() {
    InitCallback(func () {
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
        redicClient = redis.NewClient(&redis.Options{
            Addr:     v.String(),
            Password: password, 
            DB:       0,  
        })
    
        pong, err := redicClient.Ping().Result()
        fmt.Println(pong, err)
    })
}

//Redis get redis client
func Redis() *redis.Client{
    return redicClient
}