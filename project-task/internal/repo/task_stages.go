package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type TaskStagesRepo interface {
	FindStagesByProjectId(ctx context.Context, projectCode int64, page int64, pageSize int64) (list []*model.TaskStages, total int64, err error)
	FindById(ctx context.Context, stageCode int64) (ts *model.TaskStages, err error)
}
