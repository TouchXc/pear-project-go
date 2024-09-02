package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type FileDao struct {
	*gorm.DB
}

func (dao *FileDao) Save(ctx context.Context, file *model.File) (err error) {
	err = dao.DB.Model(&model.File{}).Save(&file).Error
	return
}

func (dao *FileDao) FindByIds(ctx context.Context, ids []int64) (fileList []*model.File, err error) {
	err = dao.DB.Model(&model.File{}).Where("id in (?)", ids).Find(&fileList).Error
	return
}

func NewFileDao() *FileDao {
	return &FileDao{gorms.NewDBClient()}
}
