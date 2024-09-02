package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type ProjectLogRepo interface {
	FindLogByTaskCode(ctx context.Context, taskCode int64, comment int) ([]*model.ProjectLog, int64, error)
	FindLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) ([]*model.ProjectLog, int64, error)
	SaveProjectLog(pl *model.ProjectLog)
	FindLogByMemberCode(ctx context.Context, memberId int64, page int64, size int64) ([]*model.ProjectLog, int64, error)
}
