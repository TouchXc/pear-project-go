package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
	"strconv"
)

type AccountDao struct {
	*gorm.DB
}

func (dao *AccountDao) FindAuthIdByMemberId(ctx context.Context, memberId int64) (authId int64, err error) {
	ma := &model.MemberAccount{}
	err = dao.DB.Model(&model.MemberAccount{}).
		Where("member_code = ?", memberId).
		First(&ma).Error
	authId, _ = strconv.ParseInt(ma.Authorize, 10, 64)
	return
}

func (dao *AccountDao) FindAccountList(ctx context.Context, condition string, organizationCode int64, departmentCode int64, page int64, pageSize int64) (list []*model.MemberAccount, total int64, err error) {
	err = dao.DB.Model(&model.MemberAccount{}).
		Where("organization_code = ?", organizationCode).
		Where(condition).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&list).Error
	err = dao.DB.Model(&model.MemberAccount{}).Where("organization_code = ?", organizationCode).Where(condition).Count(&total).Error
	return
}
func NewAccountDao() *AccountDao {
	return &AccountDao{gorms.NewDBClient()}
}
