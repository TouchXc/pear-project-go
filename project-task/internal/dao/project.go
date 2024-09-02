package dao

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type ProjectDao struct {
	*gorm.DB
}

func (dao *ProjectDao) FindProjectNodeAll(ctx context.Context) ([]*model.ProjectNode, error) {
	var pns []*model.ProjectNode
	err := dao.DB.Model(&model.ProjectNode{}).Find(&pns).Error
	return pns, err
}

func (dao *ProjectDao) FindProjectByPids(ctx context.Context, pids []int64) (projectList []*model.Project, err error) {
	err = dao.DB.Model(&model.Project{}).Where("id in (?)", pids).Find(&projectList).Error
	return
}

func (dao *ProjectDao) FindProjectById(ctx context.Context, projectCode int64) (pj *model.Project, err error) {
	pj = &model.Project{}
	err = dao.DB.Model(&model.Project{}).Where("id = ?", projectCode).Take(pj).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return
}

func (dao *ProjectDao) FindProjectByPid(ctx context.Context, projectCode int64) (mps []*model.MemberProject, total int64, err error) {
	err = dao.DB.Model(&model.MemberProject{}).Where("project_code = ?", projectCode).Find(&mps).Error
	err = dao.DB.Model(&model.MemberProject{}).Where("project_code = ?", projectCode).Count(&total).Error
	return
}

func (dao *ProjectDao) UpdateProject(ctx context.Context, project *model.Project) error {
	return dao.DB.Model(&model.Project{}).Update(&project).Error
}

func (dao *ProjectDao) DeleteProjectCollect(ctx context.Context, memberId int64, projectCode int64) error {
	return dao.DB.Where("member_code = ? and project_code = ?", memberId, projectCode).Delete(&model.CollectionProject{}).Error
}

func (dao *ProjectDao) SaveProjectCollect(ctx context.Context, pc *model.CollectionProject) error {
	return dao.DB.Model(&model.CollectionProject{}).Save(&pc).Error
}

func (dao *ProjectDao) UpdateDeleteProject(ctx context.Context, projectCode int64, deleted bool) error {
	code := 0
	if deleted {
		code = 1
	}
	return dao.DB.Model(&model.Project{}).Where("id = ?", projectCode).Update("deleted", code).Error
}

func (dao *ProjectDao) FindProjectByPidAndMemberId(ctx context.Context, projectCode int64, memId int64) (*model.ProjectAndMember, error) {
	pms := &model.ProjectAndMember{}
	sql := fmt.Sprintf("select a.*,b.project_code,b.member_code,b.join_time,b.is_owner,b.authorize from ms_project a,ms_project_member b where a.id = b.project_code and b.member_code = ? and b.project_code = ? limit 1")
	err := dao.DB.Raw(sql, memId, projectCode).Scan(&pms).Error
	return pms, err
}

func (dao *ProjectDao) FindCollectProjectByPidAndMemberId(ctx context.Context, projectCode int64, memberId int64) (bool, error) {
	var count int64
	err := dao.DB.Model(&model.CollectionProject{}).Where("member_code = ? AND project_code = ?", memberId, projectCode).Count(&count).Error
	return count > 0, err
}

func (dao *ProjectDao) FindCollectProjectByMemberId(ctx context.Context, id int64, page int64, size int64) (mp []*model.ProjectAndMember, total int64, err error) {
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project where id in (select project_code from ms_project_collection where member_code=? ) order by sort limit ?,?")
	err = dao.DB.Raw(sql, id, index, size).Scan(&mp).Error
	err = dao.DB.Model(&model.CollectionProject{}).Where("member_code = ?", id).Count(&total).Error
	return
}

func (dao *ProjectDao) FindProjectByMemberId(ctx context.Context, id int64, condition string, page int64, size int64) (projectAndMembers []*model.ProjectAndMember, total int64, err error) {
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project a, ms_project_member b where a.id=b.project_code and b.member_code=? %s order by sort limit ?,?", condition)
	err = dao.DB.Raw(sql, id, index, size).Scan(&projectAndMembers).Error
	//query := fmt.Sprintf("member_code = ?", condition)
	dao.DB.Model(&model.MemberProject{}).Where("member_code = ?", id).Count(&total)
	return
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{gorms.NewDBClient()}
}
