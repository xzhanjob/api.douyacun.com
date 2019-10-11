package subscribe

import (
	"context"
	"dyc/internal/db"
	"time"
)

const SubscriberIndex = "subscriber"

type Subscriber struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

func NewSubscriber(email string) *Subscriber {
	return &Subscriber{
		Email: email,
		Date:  time.Now(),
	}
}
func (s *Subscriber) Store() error {
	_, err := db.ES.Index().Index(SubscriberIndex).BodyJson(s).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
