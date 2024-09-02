package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type ProjectTemplateDao struct {
	*gorm.DB
}

func (dao *ProjectTemplateDao) FindProjectTemplateSystem(ctx context.Context, page int64, size int64) (pts []model.ProjectTemplate, total int64, err error) {
	err = dao.DB.Model(&model.ProjectTemplate{}).
		Where("is_system = ?", 1).
		Limit(size).Offset((page - 1) * size).
		Find(&pts).Error
	if err != nil {
		return pts, total, err
	}
	err = dao.DB.Model(&model.ProjectTemplate{}).Where("is_system = ?", 1).Count(&total).Error
	return pts, total, err
}

func (dao *ProjectTemplateDao) FindProjectTemplateCustom(ctx context.Context, memId int64, organizationCode int64, page int64, size int64) (pts []model.ProjectTemplate, total int64, err error) {
	err = dao.DB.Model(&model.ProjectTemplate{}).
		Where("is_system = ? AND member_code = ?  AND organization_code = ? ", 0, memId, organizationCode).
		Limit(size).Offset((page - 1) * size).
		Find(&pts).Error
	dao.DB.Model(&model.ProjectTemplate{}).
		Where("is_system = ? AND member_code = ?  AND organization_code = ? ", 0, memId, organizationCode).
		Count(&total)
	return pts, total, err
}

func (dao *ProjectTemplateDao) FindProjectTemplateAll(ctx context.Context, organizationCode int64, page int64, size int64) (pts []model.ProjectTemplate, total int64, err error) {
	err = dao.DB.Model(&model.ProjectTemplate{}).Where("organization_code = ?", organizationCode).
		Limit(size).Offset((page - 1) * size).
		Find(&pts).Error
	dao.DB.Model(&model.ProjectTemplate{}).Where("organization_code = ?", organizationCode).Count(&total)
	return pts, total, err
}

func NewProjectTemplateDao() *ProjectTemplateDao {
	return &ProjectTemplateDao{gorms.NewDBClient()}
}
