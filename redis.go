package mggo

import (
    "github.com/go-redis/redis"
    "fmt"
)

var redicClient *redis.Client

func init() {
    InitCallback(func () {
        redicClient = redis.NewClient(&redis.Options{
            Addr:     "localhost:6379",
            Password: "", 
            DB:       0,  
        })
    
        pong, err := redicClient.Ping().Result()
        fmt.Println(pong, err)
    })
}

func Redis() *redis.StatusCmd{
    redicClient = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", 
        DB:       0,  
    })
    return redicClient.Ping()
}