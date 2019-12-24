package ticket

import (
	"time"

	"github.com/google/uuid"
)

// Service ...
type Service interface {
	CreateTicket(ticket *Model) error
	FindTicketByID(id string) (*Model, error)
	FindAllTickets() ([]*Model, error)
}

type ticketService struct {
	repo Repository
}

// NewTicketService ...
func NewTicketService(repo Repository) Service {
	return &ticketService{
		repo,
	}
}

func (s *ticketService) CreateTicket(ticket *Model) error {
	ticket.ID = uuid.New().String()
	ticket.Created = time.Now()
	ticket.Updated = time.Now()
	ticket.Status = "open"
	return s.repo.Create(ticket)
}

func (s *ticketService) FindTicketByID(id string) (*Model, error) {
	return s.repo.FindByID(id)
}

func (s *ticketService) FindAllTickets() ([]*Model, error) {
	return s.repo.FindAll()
}
