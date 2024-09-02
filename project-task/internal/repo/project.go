package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type ProjectRepo interface {
	FindProjectByMemberId(ctx context.Context, id int64, condition string, page int64, size int64) ([]*model.ProjectAndMember, int64, error)
	FindCollectProjectByMemberId(ctx context.Context, id int64, page int64, size int64) ([]*model.ProjectAndMember, int64, error)
	FindProjectByPidAndMemberId(ctx context.Context, projectCode int64, memId int64) (*model.ProjectAndMember, error)
	FindCollectProjectByPidAndMemberId(ctx context.Context, projectCode int64, memberId int64) (bool, error)
	UpdateDeleteProject(ctx context.Context, projectCode int64, deleted bool) error
	SaveProjectCollect(ctx context.Context, pc *model.CollectionProject) error
	DeleteProjectCollect(ctx context.Context, memberId int64, projectCode int64) error
	UpdateProject(ctx context.Context, project *model.Project) error
	FindProjectByPid(ctx context.Context, projectCode int64) ([]*model.MemberProject, int64, error)
	FindProjectById(ctx context.Context, projectCode int64) (pj *model.Project, err error)
	FindProjectByPids(ctx context.Context, pids []int64) ([]*model.Project, error)
	FindProjectNodeAll(ctx context.Context) ([]*model.ProjectNode, error)
}
