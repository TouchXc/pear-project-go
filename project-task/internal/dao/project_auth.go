package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type ProjectAuthDao struct {
	*gorm.DB
}

func (dao *ProjectAuthDao) FindProjectAuthNodeByAuthId(ctx context.Context, authId int64) ([]string, error) {
	var err error
	var nodeList []*model.ProjectAuthNode
	err = dao.DB.Raw("select node from ms_project_auth_node where auth = ?", authId).Scan(&nodeList).Error
	nodes := make([]string, 0)
	for _, v := range nodeList {
		nodes = append(nodes, v.Node)
	}
	return nodes, err
}

func (dao *ProjectAuthDao) FindAuthListPage(ctx context.Context, organizationCode int64, page int64, pageSize int64) ([]*model.ProjectAuth, int64, error) {
	var list []*model.ProjectAuth
	var total int64
	var err error
	err = dao.DB.Model(&model.ProjectAuth{}).
		Where("organization_code = ?", organizationCode).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&list).Error
	err = dao.DB.Model(&model.ProjectAuth{}).
		Where("organization_code = ?", organizationCode).
		Count(&total).Error
	return list, total, err
}

func (dao *ProjectAuthDao) FindAuthList(ctx context.Context, organizationCode int64) (paList []*model.ProjectAuth, err error) {
	err = dao.DB.Model(&model.ProjectAuth{}).Where("organization_code=? and status=1", organizationCode).Find(&paList).Error
	return
}

func NewProjectAuthDao() *ProjectAuthDao {
	return &ProjectAuthDao{gorms.NewDBClient()}
}
