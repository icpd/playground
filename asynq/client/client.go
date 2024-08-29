package main

import (
	"log"
	"time"

	"playground/asynq/task"

	"github.com/hibiken/asynq"
)

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "172.16.0.224:6379"})

	t1, err := task.NewWelcomeEmailTask(42)
	if err != nil {
		log.Fatal(err)
	}

	t2, err := task.NewReminderEmailTask(42)
	if err != nil {
		log.Fatal(err)
	}

	t3, err := task.NewOrderCreateTask("order-id-123")
	if err != nil {
		log.Fatal(err)
	}

	info, err := client.Enqueue(t1, asynq.MaxRetry(2))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v\n", info)

	info, err = client.Enqueue(t2, asynq.ProcessIn(time.Minute))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v\n", info)

	info, err = client.Enqueue(t3)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v\n", info)
}
