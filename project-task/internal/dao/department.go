package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type DepartmentDao struct {
	*gorm.DB
}

func (dao *DepartmentDao) ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, pageSize int64) ([]*model.Department, int64, error) {
	var dpList []*model.Department
	var total int64
	var err error
	err = dao.DB.Model(&model.Department{}).Where("organization_code = ? and pcode = ?", organizationCode, parentDepartmentCode).Limit(pageSize).Offset((page - 1) * pageSize).Find(&dpList).Error
	err = dao.DB.Model(&model.Department{}).Where("organization_code = ? and pcode = ?", organizationCode, parentDepartmentCode).Count(&total).Error
	return dpList, total, err
}

func (dao *DepartmentDao) SaveDepartment(ctx context.Context, dp *model.Department) (err error) {
	err = dao.DB.Model(&model.Department{}).Save(&dp).Error
	return
}

func (dao *DepartmentDao) FindDepartment(ctx context.Context, organizationCode int64, PCode int64, name string) (*model.Department, error) {
	dp := &model.Department{}
	var err error
	if PCode > 0 {
		err = dao.DB.Model(&model.Department{}).
			Where("organization_code = ? AND name = ? AND pcode = ?", organizationCode, name, PCode).
			Take(&dp).Error
	} else {
		err = dao.DB.Model(&model.Department{}).
			Where("organization_code = ? AND name = ?", organizationCode, name).
			Take(&dp).Error
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return dp, err
}

func (dao *DepartmentDao) FindDepartmentById(departmentCode int64) (*model.Department, error) {
	dt := &model.Department{}
	err := dao.DB.Model(&model.Department{}).Where("id = ?", departmentCode).First(&dt).Error
	return dt, err
}

func NewDepartmentDao() *DepartmentDao {
	return &DepartmentDao{gorms.NewDBClient()}
}
