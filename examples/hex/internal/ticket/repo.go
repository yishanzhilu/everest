package ticket

// Repository ..
// TODO
type Repository interface {
	Create(ticket *Model) error
	FindByID(id string) (*Model, error)
	FindAll() ([]*Model, error)
}
