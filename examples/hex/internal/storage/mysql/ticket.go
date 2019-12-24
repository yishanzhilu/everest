package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/api-template/examples/hex/internal/ticket"
)

type ticketRepository struct {
	db *gorm.DB
}

// NewMysqlTicketRepository ..
func NewMysqlTicketRepository(db *gorm.DB) ticket.Repository {
	return &ticketRepository{
		db,
	}
}

func (r *ticketRepository) Create(ticket *ticket.Model) error {
	r.db.Model(ticket).FirstOrInit(ticket)
	return nil
}

func (r *ticketRepository) FindByID(id string) (*ticket.Model, error) {
	ticket := new(ticket.Model)
	err := r.db.First(&ticket, id).Error
	return ticket, err
}

func (r *ticketRepository) FindAll() (tickets []*ticket.Model, err error) {
	err = r.db.Find(&tickets).Error
	return
}
