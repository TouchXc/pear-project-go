package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type MenuRepo interface {
	FindMenus(ctx context.Context) ([]*model.ProjectMenu, error)
}
