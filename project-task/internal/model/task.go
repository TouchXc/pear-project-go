package model

import (
	"github.com/jinzhu/copier"
	"ms_project/project-common/encrypts"
	"ms_project/project-common/tms"
)

// 映射数据库表    ms_task

type Task struct {
	Id            int64
	ProjectCode   int64
	Name          string
	Pri           int
	ExecuteStatus int
	Description   string
	CreateBy      int64
	DoneBy        int64
	DoneTime      int64
	CreateTime    int64
	AssignTo      int64
	Deleted       int
	StageCode     int
	TaskTag       string
	Done          int
	BeginTime     int64
	EndTime       int64
	RemindTime    int64
	Pcode         int64
	Sort          int
	Like          int
	Star          int
	DeletedTime   int64
	Private       int
	IdNum         int
	Path          string
	Schedule      int
	VersionCode   int64
	FeaturesCode  int64
	WorkTime      int
	Status        int
}

func (*Task) TableName() string {
	return "ms_task"
}

// 任务执行状态常量
const (
	Wait = iota
	Doing
	Done
	Pause
	Cancel
	Closed
)
const (
	NoStarted = iota
	Started
)
const (
	IsNormal = iota
	IsUrgent
	IsVeryUrgent
)

func (t *Task) GetExecuteStatusStr() string {
	status := t.ExecuteStatus
	if status == Wait {
		return "wait"
	}
	if status == Doing {
		return "doing"
	}
	if status == Done {
		return "done"
	}
	if status == Pause {
		return "pause"
	}
	if status == Cancel {
		return "cancel"
	}
	if status == Closed {
		return "closed"
	}
	return ""
}

//返回接口的模型  对应接口：project/task_stages/tasks

type TaskDisplay struct {
	Id            int64
	ProjectCode   string
	Name          string
	Pri           int
	ExecuteStatus string
	Description   string
	CreateBy      string
	DoneBy        string
	DoneTime      string
	CreateTime    string
	AssignTo      string
	Deleted       int
	StageCode     string
	TaskTag       string
	Done          int
	BeginTime     string
	EndTime       string
	RemindTime    string
	Pcode         string
	Sort          int
	Like          int
	Star          int
	DeletedTime   string
	Private       int
	IdNum         int
	Path          string
	Schedule      int
	VersionCode   string
	FeaturesCode  string
	WorkTime      int
	Status        int
	Code          string
	CanRead       int
	Executor      Executors
	ProjectName   string
	StageName     string
	PriText       string
	StatusText    string
}
type MyTaskDisplay struct {
	Id                 int64
	ProjectCode        string
	Name               string
	Pri                int
	ExecuteStatus      string
	Description        string
	CreateBy           string
	DoneBy             string
	DoneTime           string
	CreateTime         string
	AssignTo           string
	Deleted            int
	StageCode          string
	TaskTag            string
	Done               int
	BeginTime          string
	EndTime            string
	RemindTime         string
	Pcode              string
	Sort               int
	Like               int
	Star               int
	DeletedTime        string
	Private            int
	IdNum              int
	Path               string
	Schedule           int
	VersionCode        string
	FeaturesCode       string
	WorkTime           int
	Status             int
	Code               string
	Cover              string `json:"cover"`
	AccessControlType  string `json:"access_control_type"`
	WhiteList          string `json:"white_list"`
	Order              int    `json:"order"`
	TemplateCode       string `json:"template_code"`
	OrganizationCode   string `json:"organization_code"`
	Prefix             string `json:"prefix"`
	OpenPrefix         int    `json:"open_prefix"`
	Archive            int    `json:"archive"`
	ArchiveTime        string `json:"archive_time"`
	OpenBeginTime      int    `json:"open_begin_time"`
	OpenTaskPrivate    int    `json:"open_task_private"`
	TaskBoardTheme     string `json:"task_board_theme"`
	AutoUpdateSchedule int    `json:"auto_update_schedule"`
	HasUnDone          int    `json:"hasUnDone"`
	ParentDone         int    `json:"parentDone"`
	PriText            string `json:"priText"`
	ProjectName        string
	Executor           *Executory
}
type Executory struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Executors struct {
	Name   string
	Avatar string
	Code   string
}

func (t *Task) GetStatusStr() string {
	status := t.Status
	if status == NoStarted {
		return "未开始"
	}
	if status == Started {
		return "开始"
	}
	return ""
}
func (t *Task) GetPriStr() string {
	status := t.Pri
	if status == IsNormal {
		return "普通"
	}
	if status == IsUrgent {
		return "紧急"
	}
	if status == IsVeryUrgent {
		return "非常紧急"
	}
	return ""
}
func (t *Task) ToTaskDisplay() *TaskDisplay {
	td := &TaskDisplay{}
	copier.Copy(td, t)
	td.CreateTime = tms.FormatByMill(t.CreateTime)
	td.DoneTime = tms.FormatByMill(t.DoneTime)
	td.BeginTime = tms.FormatByMill(t.BeginTime)
	td.EndTime = tms.FormatByMill(t.EndTime)
	td.RemindTime = tms.FormatByMill(t.RemindTime)
	td.DeletedTime = tms.FormatByMill(t.DeletedTime)
	td.CreateBy, _ = encrypts.EncryptInt64(t.CreateBy, AESKey)
	td.ProjectCode, _ = encrypts.EncryptInt64(t.ProjectCode, AESKey)
	td.DoneBy, _ = encrypts.EncryptInt64(t.DoneBy, AESKey)
	td.AssignTo, _ = encrypts.EncryptInt64(t.AssignTo, AESKey)
	td.StageCode, _ = encrypts.EncryptInt64(int64(t.StageCode), AESKey)
	td.Pcode, _ = encrypts.EncryptInt64(t.Pcode, AESKey)
	td.VersionCode, _ = encrypts.EncryptInt64(t.VersionCode, AESKey)
	td.FeaturesCode, _ = encrypts.EncryptInt64(t.FeaturesCode, AESKey)
	td.ExecuteStatus = t.GetExecuteStatusStr()
	td.Code, _ = encrypts.EncryptInt64(t.Id, AESKey)
	td.CanRead = 1
	td.StatusText = t.GetStatusStr()
	td.PriText = t.GetPriStr()
	return td
}
func (t *Task) ToMyTaskDisplay(p *Project, name string, avatar string) *MyTaskDisplay {
	td := &MyTaskDisplay{}
	copier.Copy(td, p)
	copier.Copy(td, t)
	td.Executor = &Executory{
		Name:   name,
		Avatar: avatar,
	}
	td.ProjectName = p.Name
	td.CreateTime = tms.FormatByMill(t.CreateTime)
	td.DoneTime = tms.FormatByMill(t.DoneTime)
	td.BeginTime = tms.FormatByMill(t.BeginTime)
	td.EndTime = tms.FormatByMill(t.EndTime)
	td.RemindTime = tms.FormatByMill(t.RemindTime)
	td.DeletedTime = tms.FormatByMill(t.DeletedTime)
	td.CreateBy, _ = encrypts.EncryptInt64(t.CreateBy, AESKey)
	td.ProjectCode, _ = encrypts.EncryptInt64(t.ProjectCode, AESKey)
	td.DoneBy, _ = encrypts.EncryptInt64(t.DoneBy, AESKey)
	td.AssignTo, _ = encrypts.EncryptInt64(t.AssignTo, AESKey)
	td.StageCode, _ = encrypts.EncryptInt64(int64(t.StageCode), AESKey)
	td.Pcode, _ = encrypts.EncryptInt64(t.Pcode, AESKey)
	td.VersionCode, _ = encrypts.EncryptInt64(t.VersionCode, AESKey)
	td.FeaturesCode, _ = encrypts.EncryptInt64(t.FeaturesCode, AESKey)
	td.ExecuteStatus = t.GetExecuteStatusStr()
	td.Code, _ = encrypts.EncryptInt64(t.Id, AESKey)
	td.AccessControlType = p.GetAccessControlType()
	td.ArchiveTime = tms.FormatByMill(p.ArchiveTime)
	td.TemplateCode, _ = encrypts.EncryptInt64(int64(p.TemplateCode), AESKey)
	td.OrganizationCode, _ = encrypts.EncryptInt64(p.OrganizationCode, AESKey)
	return td
}
