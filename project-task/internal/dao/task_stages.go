package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type TaskStagesDao struct {
	*gorm.DB
}

func (dao *TaskStagesDao) FindById(ctx context.Context, stageCode int64) (ts *model.TaskStages, err error) {
	ts = &model.TaskStages{}
	err = dao.DB.Model(&model.TaskStages{}).Where("id = ?", stageCode).Take(ts).Error
	return
}

func (dao *TaskStagesDao) FindStagesByProjectId(ctx context.Context, projectCode int64, page int64, pageSize int64) (list []*model.TaskStages, total int64, err error) {
	err = dao.DB.Model(&model.TaskStages{}).
		Where("project_code = ?", projectCode).
		Order("sort asc").
		Find(&list).Limit(pageSize).Offset((page - 1) * pageSize).Error
	err = dao.DB.Model(&model.TaskStages{}).Where("project_code = ?", projectCode).Count(&total).Error
	return
}

func NewTaskStagesDao() *TaskStagesDao {
	return &TaskStagesDao{gorms.NewDBClient()}
}
