package dao

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"ms_project/project-user/internal/database/gorms"
	"ms_project/project-user/internal/model"
)

// dao层 对数据库进行操作

type MemberDao struct {
	*gorm.DB
}

func (dao *MemberDao) FindMemberByIds(ctx context.Context, ids []int64) (mList []*model.Member, err error) {
	if len(ids) <= 0 {
		return nil, nil
	}
	err = dao.DB.Model(&model.Member{}).Where("id in (?)", ids).Find(&mList).Error
	return
}

func (dao *MemberDao) FindMemberById(ctx context.Context, id int64) (*model.Member, error) {
	mem := &model.Member{}
	err := dao.DB.Model(&model.Member{}).Where("id = ?", id).First(mem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return mem, err
}

func NewMemberDao() *MemberDao {
	return &MemberDao{gorms.NewDBClient()}
}
func (dao *MemberDao) FindMemberByAccount(ctx context.Context, account string, pwd string) (*model.Member, error) {
	member := &model.Member{}
	err := dao.DB.Model(&model.Member{}).Where("account = ? AND password = ?", account, pwd).First(member).Error
	return member, err
}

func (dao *MemberDao) GetMemberByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := dao.DB.Model(&model.Member{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (dao *MemberDao) GetMemberByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := dao.DB.Model(&model.Member{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (dao *MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := dao.DB.Model(&model.Member{}).Where("mobile = ?", mobile).Count(&count).Error
	return count > 0, err
}
func (dao *MemberDao) SaveMember(ctx context.Context, member *model.Member) error {
	return dao.DB.Model(&model.Member{}).Create(member).Error
}
