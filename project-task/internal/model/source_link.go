package model

import (
	"github.com/jinzhu/copier"
	"ms_project/project-common/encrypts"
	"ms_project/project-common/tms"
)

type SourceLink struct {
	Id               int64
	SourceType       string
	SourceCode       int64
	LinkType         string
	LinkCode         int64
	OrganizationCode int64
	CreateBy         int64
	CreateTime       int64
	Sort             int
}

func (*SourceLink) TableName() string {
	return "ms_source_link"
}

type SourceLinkDisplay struct {
	Id               int64        `json:"id"`
	Code             string       `json:"code"`
	SourceType       string       `json:"source_type"`
	SourceCode       string       `json:"source_code"`
	LinkType         string       `json:"link_type"`
	LinkCode         string       `json:"link_code"`
	OrganizationCode string       `json:"organization_code"`
	CreateBy         string       `json:"create_by"`
	CreateTime       string       `json:"create_time"`
	Sort             int          `json:"sort"`
	Title            string       `json:"title"`
	SourceDetail     SourceDetail `json:"sourceDetail"`
}

type SourceDetail struct {
	Id               int64  `json:"id"`
	Code             string `json:"code"`
	PathName         string `json:"path_name"`
	Title            string `json:"title"`
	Extension        string `json:"extension"`
	Size             int    `json:"size"`
	ObjectType       string `json:"object_type"`
	OrganizationCode string `json:"organization_code"`
	TaskCode         string `json:"task_code"`
	ProjectCode      string `json:"project_code"`
	CreateBy         string `json:"create_by"`
	CreateTime       string `json:"create_time"`
	Downloads        int    `json:"downloads"`
	Extra            string `json:"extra"`
	Deleted          int    `json:"deleted"`
	FileUrl          string `json:"file_url"`
	FileType         string `json:"file_type"`
	DeletedTime      string `json:"deleted_time"`
	ProjectName      string `json:"projectName"`
	FullName         string `json:"fullName"`
}

func (s *SourceLink) ToDisplay(f *File) *SourceLinkDisplay {
	sl := &SourceLinkDisplay{}
	copier.Copy(sl, s)
	sl.SourceDetail = SourceDetail{}
	copier.Copy(&sl.SourceDetail, f)
	sl.LinkCode, _ = encrypts.EncryptInt64(s.LinkCode, AESKey)
	sl.OrganizationCode, _ = encrypts.EncryptInt64(s.OrganizationCode, AESKey)
	sl.CreateTime = tms.FormatByMill(s.CreateTime)
	sl.CreateBy, _ = encrypts.EncryptInt64(s.CreateBy, AESKey)
	sl.SourceCode, _ = encrypts.EncryptInt64(s.SourceCode, AESKey)
	sl.SourceDetail.OrganizationCode, _ = encrypts.EncryptInt64(f.OrganizationCode, AESKey)
	sl.SourceDetail.CreateBy, _ = encrypts.EncryptInt64(f.CreateBy, AESKey)
	sl.SourceDetail.CreateTime = tms.FormatByMill(f.CreateTime)
	sl.SourceDetail.DeletedTime = tms.FormatByMill(f.DeletedTime)
	sl.SourceDetail.TaskCode, _ = encrypts.EncryptInt64(f.TaskCode, AESKey)
	sl.SourceDetail.ProjectCode, _ = encrypts.EncryptInt64(f.ProjectCode, AESKey)
	sl.SourceDetail.FullName = f.Title
	sl.Title = f.Title
	return sl
}
