package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type AccountRepo interface {
	FindAccountList(ctx context.Context, condition string, organizationCode int64, departmentCode int64, page int64, pageSize int64) ([]*model.MemberAccount, int64, error)
	FindAuthIdByMemberId(ctx context.Context, memberId int64) (int64, error)
}
