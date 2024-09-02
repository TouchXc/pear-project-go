package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type SourceLinkRepo interface {
	Save(ctx context.Context, link *model.SourceLink) error
	FindByTaskCode(ctx context.Context, taskCode int64) ([]*model.SourceLink, error)
}
