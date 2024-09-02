package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type DepartmentRepo interface {
	FindDepartmentById(departmentCode int64) (*model.Department, error)
	FindDepartment(ctx context.Context, organizationCode int64, PCode int64, name string) (*model.Department, error)
	SaveDepartment(ctx context.Context, dp *model.Department) error
	ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, pageSize int64) ([]*model.Department, int64, error)
}
