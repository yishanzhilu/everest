package workspace

import (
	"github.com/yishanzhilu/everest/pkg/common"
)

// WorkprofileModel define workprofile properties
type WorkprofileModel struct {
	common.BaseModel
	Minutes int `json:"minutes" db:"minutes"`
}
