package repo

import (
	"context"
	"ms_project/project-user/internal/model"
)

// Mysql 接口层

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByName(ctx context.Context, name string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(ctx context.Context, member *model.Member) error
	FindMemberByAccount(ctx context.Context, account string, pwd string) (*model.Member, error)
	FindMemberById(ctx context.Context, id int64) (*model.Member, error)
	FindMemberByIds(ctx context.Context, ids []int64) ([]*model.Member, error)
}
