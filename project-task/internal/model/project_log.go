package model

import (
	"github.com/jinzhu/copier"
	"ms_project/project-common/encrypts"
	"ms_project/project-common/tms"
)

type ProjectLog struct {
	Id           int64
	MemberCode   int64
	Content      string
	Remark       string
	Type         string
	CreateTime   int64
	SourceCode   int64
	ActionType   string
	ToMemberCode int64
	IsComment    int
	ProjectCode  int64
	Icon         string
	IsRobot      int
}

func (*ProjectLog) TableName() string {
	return "ms_project_log"
}

type ProjectLogDisplay struct {
	Id           int64
	MemberCode   string
	Content      string
	Remark       string
	Type         string
	CreateTime   string
	SourceCode   string
	ActionType   string
	ToMemberCode string
	IsComment    int
	ProjectCode  string
	Icon         string
	IsRobot      int
	Member       Member
}
type Member struct {
	Id     int64
	Name   string
	Avatar string
	Code   string
}

func (l *ProjectLog) ToDisplay() *ProjectLogDisplay {
	pd := &ProjectLogDisplay{}
	copier.Copy(pd, l)
	pd.MemberCode, _ = encrypts.EncryptInt64(l.MemberCode, AESKey)
	pd.ToMemberCode, _ = encrypts.EncryptInt64(l.ToMemberCode, AESKey)
	pd.ProjectCode, _ = encrypts.EncryptInt64(l.ProjectCode, AESKey)
	pd.CreateTime = tms.FormatByMill(l.CreateTime)
	pd.SourceCode, _ = encrypts.EncryptInt64(l.SourceCode, AESKey)
	return pd
}

type IndexProjectLogDisplay struct {
	Content      string
	Remark       string
	CreateTime   string
	SourceCode   string
	IsComment    int
	ProjectCode  string
	MemberAvatar string
	MemberName   string
	ProjectName  string
	TaskName     string
}

func (l *ProjectLog) ToIndexDisplay() *IndexProjectLogDisplay {
	pd := &IndexProjectLogDisplay{}
	copier.Copy(pd, l)
	pd.ProjectCode, _ = encrypts.EncryptInt64(l.ProjectCode, AESKey)
	pd.CreateTime = tms.FormatByMill(l.CreateTime)
	pd.SourceCode, _ = encrypts.EncryptInt64(l.SourceCode, AESKey)
	return pd
}
