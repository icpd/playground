package main

import (
	"fmt"
	"log"
	"time"

	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisOptions := &redis.Options{
		Addr: "172.16.0.224:6379",
	}
	redisClient := redis.NewClient(redisOptions)
	locker, err := redislock.NewRedisLocker(redisClient, redislock.WithTries(1))
	if err != nil {
		panic(err)
	}

	s, err := gocron.NewScheduler(gocron.WithDistributedLocker(locker))
	if err != nil {
		// handle the error
	}
	j, err := s.NewJob(gocron.DurationJob(3*time.Second), gocron.NewTask(func() {
		time.Sleep(3 * time.Second)
		log.Println("what are you doing?")
	}))
	if err != nil {
		panic(err)
	}

	fmt.Println(j.ID())

	s.Start()

	select {}
}
