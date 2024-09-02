package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-user/internal/database/gorms"
	"ms_project/project-user/internal/model"
)

type OrganizationDao struct {
	*gorm.DB
}

func (dao *OrganizationDao) FindOrganizationByMemberId(ctx context.Context, id int64) (orgs []*model.Organization, err error) {
	err = dao.DB.Model(&model.Organization{}).Where("member_id = ?", id).Find(&orgs).Error
	return
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{gorms.NewDBClient()}
}
func (dao *OrganizationDao) SaveOrganization(ctx context.Context, organization *model.Organization) error {
	return dao.DB.Model(&model.Organization{}).Create(organization).Error
}
