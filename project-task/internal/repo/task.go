package repo

import (
	"context"
	"ms_project/project-task/internal/model"
)

type TaskRepo interface {
	FindTaskByStageCode(ctx context.Context, stageCode int64) ([]*model.Task, error)
	FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberId int64) (*model.TaskMember, error)
	FindTaskMaxIdNum(ctx context.Context, projectCode int64) (int, error)
	FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (int, error)
	FindTaskById(ctx context.Context, taskCode int64) (*model.Task, error)
	FindTaskByAssignTo(ctx context.Context, memberId int64, t int32, page int64, pageSize int64) ([]*model.Task, int64, error)
	FindTaskByMemberId(ctx context.Context, memberId int64, t int32, page int64, pageSize int64) ([]*model.Task, int64, error)
	FindTaskByCreateBy(ctx context.Context, memberId int64, t int32, page int64, pageSize int64) ([]*model.Task, int64, error)
	FindTaskMemberPage(ctx context.Context, taskCode int64, page int64, pageSize int64) ([]*model.TaskMember, int64, error)
	SaveTaskWorkTime(ctx context.Context, twt *model.TaskWorkTime) error
	FindWorkTimeList(ctx context.Context, taskCode int64) ([]*model.TaskWorkTime, error)
	FindTaskByIds(ctx context.Context, taskIdList []int64) ([]*model.Task, error)
}
