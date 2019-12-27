package workspace

import (
	"github.com/yishanzhilu/everest/pkg/common"
)

// TaskModel is user created scalar job.
type TaskModel struct {
	common.BaseModel
	Content string
	Review  string
}
