package model

import (
	"ms_project/project-common/encrypts"
	"ms_project/project-common/tms"
)

type ProjectTemplate struct {
	Id               int
	Name             string
	Description      string
	Sort             int
	CreateTime       int64
	OrganizationCode int64
	Cover            string
	MemberCode       int64
	IsSystem         int
}

func (*ProjectTemplate) TableName() string {
	return "ms_project_template"
}

type ProjectTemplateAll struct {
	Id               int
	Name             string
	Description      string
	Sort             int
	CreateTime       string
	OrganizationCode string
	Cover            string
	MemberCode       string
	IsSystem         int
	TaskStages       []*TaskStagesOnlyName
	Code             string
}

func (pt *ProjectTemplate) Convert(taskStages []*TaskStagesOnlyName) *ProjectTemplateAll {
	organizationCode, _ := encrypts.EncryptInt64(pt.OrganizationCode, AESKey)
	memberCode, _ := encrypts.EncryptInt64(pt.MemberCode, AESKey)
	code, _ := encrypts.EncryptInt64(int64(pt.Id), AESKey)
	return &ProjectTemplateAll{
		Id:               pt.Id,
		Name:             pt.Name,
		Description:      pt.Description,
		Sort:             pt.Sort,
		CreateTime:       tms.FormatByMill(pt.CreateTime),
		OrganizationCode: organizationCode,
		Cover:            pt.Cover,
		MemberCode:       memberCode,
		IsSystem:         pt.IsSystem,
		TaskStages:       taskStages,
		Code:             code,
	}
}
func ToProjectTemplateIds(pts []ProjectTemplate) []int {
	var ids []int
	for _, v := range pts {
		ids = append(ids, v.Id)
	}
	return ids
}
