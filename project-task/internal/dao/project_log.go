package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type ProjectLogDao struct {
	*gorm.DB
}

func (dao *ProjectLogDao) FindLogByMemberCode(ctx context.Context, memberId int64, page int64, size int64) (prList []*model.ProjectLog, total int64, err error) {
	err = dao.DB.Model(&model.ProjectLog{}).Where("member_code = ?", memberId).Find(&prList).Limit(size).Offset((page - 1) * size).Order("create_time desc").Error
	err = dao.DB.Model(&model.ProjectLog{}).Where("member_code = ?", memberId).Count(&total).Error
	return
}

func (dao *ProjectLogDao) SaveProjectLog(pl *model.ProjectLog) {
	dao.DB.Model(&model.ProjectLog{}).Save(&pl)
}

//不分页

func (dao *ProjectLogDao) FindLogByTaskCode(ctx context.Context, taskCode int64, comment int) (list []*model.ProjectLog, total int64, err error) {
	if comment == 1 {
		err = dao.DB.Model(&model.ProjectLog{}).Where("source_code=? and is_comment=?", taskCode, comment).Find(&list).Error
		err = dao.DB.Model(&model.ProjectLog{}).Where("source_code=? and is_comment=?", taskCode, comment).Count(&total).Error
	} else {
		err = dao.DB.Model(&model.ProjectLog{}).Where("source_code=?", taskCode).Find(&list).Error
		err = dao.DB.Model(&model.ProjectLog{}).Where("source_code=?", taskCode).Count(&total).Error
	}
	return
}

//分页

func (dao *ProjectLogDao) FindLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) ([]*model.ProjectLog, int64, error) {
	offset := (page - 1) * pageSize
	var list []*model.ProjectLog
	var total int64
	var err error
	if comment == 1 {
		err = dao.DB.Model(&model.ProjectLog{}).Where("source_code=? and is_comment=?", taskCode, comment).Limit(pageSize).Offset(offset).Find(&list).Error
		err = dao.DB.Model(&model.ProjectLog{}).Where("source_code=? and is_comment=?", taskCode, comment).Count(&total).Error
	} else {
		err = dao.DB.Model(&model.ProjectLog{}).Where("source_code=?", taskCode).Limit(pageSize).Offset(offset).Find(&list).Error
		err = dao.DB.Model(&model.ProjectLog{}).Where("source_code=?", taskCode).Count(&total).Error
	}
	return list, total, err
}
func NewProjectLogDao() *ProjectLogDao {
	return &ProjectLogDao{gorms.NewDBClient()}
}
