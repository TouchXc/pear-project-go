package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
)

type TaskDao struct {
	*gorm.DB
}

func (dao *TaskDao) FindTaskByIds(ctx context.Context, taskIdList []int64) (taskList []*model.Task, err error) {
	err = dao.DB.Model(&model.Task{}).Where("id in (?)", taskIdList).Find(&taskList).Error
	return
}

func (dao *TaskDao) SaveTaskWorkTime(ctx context.Context, twt *model.TaskWorkTime) (err error) {
	err = dao.DB.Model(&model.TaskWorkTime{}).Save(&twt).Error
	return
}

func (dao *TaskDao) FindWorkTimeList(ctx context.Context, taskCode int64) (twtList []*model.TaskWorkTime, err error) {
	err = dao.DB.Model(&model.TaskWorkTime{}).Where("task_code = ?", taskCode).Find(&twtList).Error
	return
}

func (dao *TaskDao) FindTaskMemberPage(ctx context.Context, taskCode int64, page int64, pageSize int64) ([]*model.TaskMember, int64, error) {
	var list []*model.TaskMember
	var total int64
	err := dao.DB.Model(&model.TaskMember{}).Where("task_code = ?", taskCode).Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	err = dao.DB.Model(&model.TaskMember{}).Where("task_code = ?", taskCode).Count(&total).Error
	return list, total, err
}

func (dao *TaskDao) FindTaskByAssignTo(ctx context.Context, memberId int64, t int32, page int64, pageSize int64) (taskList []*model.Task, total int64, err error) {
	err = dao.DB.Model(&model.Task{}).Where("assign_to = ? and deleted = 0 and done = ?", memberId, t).Limit(pageSize).Offset((page - 1) * pageSize).Find(&taskList).Error
	err = dao.DB.Model(&model.Task{}).Where("assign_to = ? and deleted = 0 and done = ?", memberId, t).Count(&total).Error
	return
}

type Total struct {
	total int64
}

func (dao *TaskDao) FindTaskByMemberId(ctx context.Context, memberId int64, t int32, page int64, pageSize int64) (taskList []*model.Task, total int64, err error) {
	var T Total
	sql := "select a.* from ms_task a,ms_task_member b where a.id=b.task_code and member_code=? and a.deleted=0 and a.done=?"
	err = dao.DB.Raw(sql, memberId, t).Limit(pageSize).Offset((page - 1) * pageSize).Scan(&taskList).Error
	if err != nil {
		return nil, 0, err
	}
	sqlCount := "select count(*) from ms_task a,ms_task_member b where a.id=b.task_code and member_code=? and a.deleted=0 and a.done=?"
	err = dao.DB.Raw(sqlCount, memberId, t).Scan(&T).Error
	return taskList, T.total, err
}

func (dao *TaskDao) FindTaskByCreateBy(ctx context.Context, memberId int64, t int32, page int64, pageSize int64) (taskList []*model.Task, total int64, err error) {
	err = dao.DB.Model(&model.Task{}).Where("create_by = ? and deleted = 0 and done = ?", memberId, t).Limit(pageSize).Offset((page - 1) * pageSize).Find(&taskList).Error
	err = dao.DB.Model(&model.Task{}).Where("create_by = ? and deleted = 0 and done = ?", memberId, t).Count(&total).Error
	return
}

func (dao *TaskDao) FindTaskById(ctx context.Context, taskCode int64) (*model.Task, error) {
	ts := &model.Task{}
	err := dao.DB.Model(&model.Task{}).Where("id = ?", taskCode).First(&ts).Error
	return ts, err
}

func (dao *TaskDao) FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (maxSort int, err error) {
	task := &model.Task{}
	err = dao.DB.Model(&model.Task{}).
		Where("project_code = ? AND stage_code = ?", projectCode, stageCode).
		Order("sort desc").First(task).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	}
	return task.Sort, err
}

func (dao *TaskDao) FindTaskMaxIdNum(ctx context.Context, projectCode int64) (int, error) {
	task := &model.Task{}
	err := dao.DB.Model(&model.Task{}).Where("project_code = ?", projectCode).Order("id_num desc").First(task).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	}
	return task.IdNum, err
}

func (dao *TaskDao) FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberId int64) (task *model.TaskMember, err error) {
	err = dao.DB.Model(&model.TaskMember{}).Where("task_code = ? and member_code = ?", taskCode, memberId).Take(&task).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return
}

func (dao *TaskDao) FindTaskByStageCode(ctx context.Context, stageCode int64) (list []*model.Task, err error) {
	err = dao.DB.Model(&model.Task{}).Where("stage_code = ? and deleted = 0", stageCode).Order("sort asc").Find(&list).Error
	return
}

func NewTaskDao() *TaskDao {
	return &TaskDao{gorms.NewDBClient()}
}
