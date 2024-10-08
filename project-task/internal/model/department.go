package model

import (
	"github.com/jinzhu/copier"
	"ms_project/project-common/encrypts"
	"ms_project/project-common/tms"
)

type Department struct {
	Id               int64
	OrganizationCode int64
	Name             string
	Sort             int
	Pcode            int64
	icon             string
	CreateTime       int64
	Path             string
}

func (*Department) TableName() string {
	return "ms_department"
}

type DepartmentDisplay struct {
	Id               int64
	OrganizationCode string
	Name             string
	Sort             int
	Pcode            string
	icon             string
	CreateTime       string
	Path             string
}

func (d *Department) ToDisplay() *DepartmentDisplay {
	dp := &DepartmentDisplay{}
	copier.Copy(dp, d)
	dp.CreateTime = tms.FormatByMill(d.CreateTime)
	dp.OrganizationCode, _ = encrypts.EncryptInt64(d.OrganizationCode, AESKey)
	if d.Pcode > 0 {
		dp.Pcode, _ = encrypts.EncryptInt64(d.Pcode, AESKey)
	} else {
		dp.Pcode = ""
	}
	return dp
}
