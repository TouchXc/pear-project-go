package model

import (
	"github.com/jinzhu/copier"
	"ms_project/project-common/encrypts"
	"ms_project/project-common/tms"
)

type TaskWorkTime struct {
	Id         int64
	TaskCode   int64
	MemberCode int64
	CreateTime int64
	Content    string
	BeginTime  int64
	Num        int
}

func (*TaskWorkTime) TableName() string {
	return "ms_task_work_time"
}

type TaskWorkTimeDisplay struct {
	Id         int64
	TaskCode   string
	MemberCode string
	CreateTime string
	Content    string
	BeginTime  string
	Num        int
	Member     Member
}

func (t *TaskWorkTime) ToDisplay() *TaskWorkTimeDisplay {
	td := &TaskWorkTimeDisplay{}
	copier.Copy(td, t)
	td.MemberCode, _ = encrypts.EncryptInt64(t.MemberCode, AESKey)
	td.TaskCode, _ = encrypts.EncryptInt64(t.TaskCode, AESKey)
	td.CreateTime = tms.FormatByMill(t.CreateTime)
	td.BeginTime = tms.FormatByMill(t.BeginTime)
	return td
}
