package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type MenuDao struct {
	*gorm.DB
}

func NewMenuDao() *MenuDao {
	return &MenuDao{gorms.NewDBClient()}
}

func (dao *MenuDao) FindMenus(ctx context.Context) (promenus []*model.ProjectMenu, err error) {
	err = dao.DB.Model(&model.ProjectMenu{}).Order("pid ,sort asc,id asc").Find(&promenus).Error
	return
}
