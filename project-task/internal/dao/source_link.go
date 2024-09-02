package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type SourceLinkDao struct {
	*gorm.DB
}

func (dao *SourceLinkDao) Save(ctx context.Context, link *model.SourceLink) (err error) {
	err = dao.DB.Model(&model.SourceLink{}).Save(&link).Error
	return
}

func (dao *SourceLinkDao) FindByTaskCode(ctx context.Context, taskCode int64) (sourceLinkList []*model.SourceLink, err error) {
	err = dao.DB.Model(&model.SourceLink{}).Where("link_type = ? and link_code = ?", "task", taskCode).Find(&sourceLinkList).Error
	return
}

func NewSourceLinkDao() *SourceLinkDao {
	return &SourceLinkDao{gorms.NewDBClient()}
}
