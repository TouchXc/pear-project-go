package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type FileRepo interface {
	Save(ctx context.Context, file *model.File) error
	FindByIds(background context.Context, ids []int64) ([]*model.File, error)
}
