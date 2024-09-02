package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type ProjectAuthRepo interface {
	FindAuthList(ctx context.Context, organizationCode int64) ([]*model.ProjectAuth, error)
	FindAuthListPage(ctx context.Context, organizationCode int64, page int64, pageSize int64) ([]*model.ProjectAuth, int64, error)
	FindProjectAuthNodeByAuthId(ctx context.Context, authId int64) ([]string, error)
}
