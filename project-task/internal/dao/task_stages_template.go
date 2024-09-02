package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type TaskStagesTemplateDao struct {
	*gorm.DB
}

func (dao *TaskStagesTemplateDao) FindStagesByProjectTemplateCode(ctx context.Context, ptCode int) (list []*model.MsTaskStagesTemplate, err error) {
	err = dao.DB.Model(&model.MsTaskStagesTemplate{}).
		Where("project_template_code = ?", ptCode).
		Order("sort desc,id asc").
		Find(&list).Error
	return
}

func (dao *TaskStagesTemplateDao) FindInProTemIds(ctx context.Context, ids []int) ([]model.MsTaskStagesTemplate, error) {
	var tsts []model.MsTaskStagesTemplate
	err := dao.DB.Model(&model.MsTaskStagesTemplate{}).Where("project_template_code in (?)", ids).Find(&tsts).Error
	return tsts, err
}

func NewTaskStagesTemplateDao() *TaskStagesTemplateDao {
	return &TaskStagesTemplateDao{gorms.NewDBClient()}
}
