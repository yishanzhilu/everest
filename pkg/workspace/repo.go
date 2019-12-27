package workspace

// WorkprofileRepository abstract workprofile repo layer logic
type WorkprofileRepository interface {
	Create(workprofile *WorkprofileModel) error
	FindByID(id string) (*WorkprofileModel, error)
	FindList() ([]WorkprofileModel, error)
}
