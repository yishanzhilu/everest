package workspace

import (
	"github.com/rs/xid"
)

// WorkprofileService abstract workspace business logic
type WorkprofileService interface {
	CreateWorkprofile(workprofile *WorkprofileModel) error
	FindByID(id string) (*WorkprofileModel, error)
	FindList() ([]WorkprofileModel, error)
}

type workProfileService struct {
	repo WorkprofileRepository
}

// NewWorkprofileService 需要一个 WorkprofileRepository
func NewWorkprofileService(repo WorkprofileRepository) WorkprofileService {
	return &workProfileService{
		repo,
	}
}

func (s *workProfileService) CreateWorkprofile(workprofile *WorkprofileModel) error {
	workprofile.ID = xid.New().String()
	return s.repo.Create(workprofile)
}

func (s *workProfileService) FindByID(id string) (*WorkprofileModel, error) {
	return s.repo.FindByID(id)
}

func (s *workProfileService) FindList() ([]WorkprofileModel, error) {
	return s.repo.FindList()
}
