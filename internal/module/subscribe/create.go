package subscribe

import (
	"context"
	"dyc/internal/db"
	"time"
)

const SubscriberIndex = "subscriber"

var Email _subscriber

type _subscriber struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

func (s *_subscriber) Store(email string) error {
	s.Email = email
	_, err := db.ES.Index().Index(SubscriberIndex).BodyJson(s).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
