package redis

import (
	"encoding/json"

	"github.com/yishanzhilu/api-template/examples/hex/internal/ticket"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const ticketTable = "tickets"

type ticketRepository struct {
	client       *redis.Client
	DBRepository ticket.Repository
}

// NewRedisModelRepository ..
func NewRedisModelRepository(client *redis.Client, DBRepository ticket.Repository) ticket.Repository {
	return &ticketRepository{
		client,
		DBRepository,
	}
}

func (r *ticketRepository) Create(ticket *ticket.Model) error {
	encoded, err := json.Marshal(ticket)

	if err != nil {
		logrus.Error("Unable to marshal ticket")
		return err
	}

	r.client.HSet(ticketTable, ticket.ID, encoded) //Don't expire
	return nil
}

func (r *ticketRepository) FindByID(id string) (*ticket.Model, error) {
	b, err := r.client.HGet(ticketTable, id).Bytes()

	if err == nil {
		t := new(ticket.Model)
		err = json.Unmarshal(b, t)

		if err != nil {
			logrus.WithField("id", id).Error("Unable to unmarshal ticket")
			return nil, err
		}

		return t, nil

	}
	logrus.WithField("id", id).Info("Cache contains no Model with id")
	t, err := r.DBRepository.FindByID(id)
	if err != nil {
		logrus.WithField("id", id).Error("Unable to find ticket by id from DBRepository")
		return nil, err
	}
	err = r.Create(t)
	if err != nil {
		logrus.WithField("id", id).Error("Unable to create ticket in cache")
		return nil, err
	}
	return t, nil

}

func (r *ticketRepository) FindAll() (tickets []*ticket.Model, err error) {
	ts := r.client.HGetAll(ticketTable).Val()
	for key, value := range ts {
		t := new(ticket.Model)
		err = json.Unmarshal([]byte(value), t)

		if err != nil {
			logrus.WithField("id", key).Error("Unable to unmarshal ticket")
			return nil, err
		}

		t.ID = key
		tickets = append(tickets, t)
	}
	return tickets, nil
}
