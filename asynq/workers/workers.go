package main

import (
	"context"
	"log"
	"time"

	"playground/asynq/task"

	"github.com/hibiken/asynq"
)

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "172.16.0.224:6379"},
		asynq.Config{Concurrency: 10},
	)

	emailHandler := asynq.NewServeMux()
	emailHandler.HandleFunc(task.TypeWelcomeEmail, task.HandleWelcomeEmailTask)
	emailHandler.HandleFunc(task.TypeReminderEmail, task.HandleReminderEmailTask)

	orderHandler := asynq.NewServeMux()
	orderHandler.HandleFunc(task.TypeCreateOrder, task.HandleOrderCreateTask)

	mux := asynq.NewServeMux()
	mux.Use(func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) (err error) {
			start := time.Now()
			log.Printf("middleware log, start task:%s %s\n", t.Type(), t.Payload())
			defer func() {
				duration := time.Since(start)
				if duration > time.Second {
					duration = duration.Truncate(time.Second)
				}
				if err != nil {
					log.Printf("middleware log, faild task:%s %s duration %s\n", t.Type(), t.Payload(), duration)
					return
				}
				log.Printf("middleware log, success task:%s %s, duration %s\n", t.Type(), t.Payload(), duration)
			}()

			return h.ProcessTask(ctx, t)
		})
	})
	mux.Handle(task.TypeDefaultEmail, emailHandler)
	mux.Handle(task.TypeDefaultOrder, orderHandler)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
