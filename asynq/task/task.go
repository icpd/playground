package task

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeWelcomeEmail  = "email:welcome"
	TypeReminderEmail = "email:reminder"
	TypeDefaultEmail  = "email:"
	TypeCreateOrder   = "order:create"
	TypeDefaultOrder  = "order:"
)

// Task payload for any email related tasks.
type emailTaskPayload struct {
	// ID for the email recipient.
	UserID int
}

func NewWelcomeEmailTask(id int) (*asynq.Task, error) {
	payload, err := json.Marshal(emailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWelcomeEmail, payload), nil
}

func NewReminderEmailTask(id int) (*asynq.Task, error) {
	payload, err := json.Marshal(emailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeReminderEmail, payload), nil
}

func NewOrderCreateTask(id string) (*asynq.Task, error) {
	return asynq.NewTask(TypeCreateOrder, []byte(id)), nil
}

var attempt int

func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	defer func() {
		attempt++
	}()

	time.Sleep(5 * time.Second)
	if attempt < 1 {
		return errors.New("not implemented")
	}

	var p emailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Send Welcome Email to User %d", p.UserID)
	return nil
}

func HandleReminderEmailTask(ctx context.Context, t *asynq.Task) error {
	var p emailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Send Reminder Email to User %d", p.UserID)
	return nil
}

func HandleOrderCreateTask(ctx context.Context, t *asynq.Task) error {
	time.Sleep(30 * time.Second)
	log.Printf(" [*] Handle Order Create Task, order: %s", t.Payload())
	return nil
}
