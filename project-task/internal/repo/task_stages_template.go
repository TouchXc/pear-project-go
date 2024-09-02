package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type TaskStagesTemplateRepo interface {
	FindInProTemIds(ctx context.Context, id []int) ([]model.MsTaskStagesTemplate, error)
	FindStagesByProjectTemplateCode(ctx context.Context, ptCode int) (list []*model.MsTaskStagesTemplate, err error)
}
