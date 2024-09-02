package model

import (
	"github.com/jinzhu/copier"
	"ms_project/project-common/encrypts"
	"ms_project/project-common/tms"
)

type MemberAccount struct {
	Id               int64
	OrganizationCode int64
	DepartmentCode   int64
	MemberCode       int64
	Authorize        string
	IsOwner          int
	Name             string
	Mobile           string
	Email            string
	CreateTime       int64
	LastLoginTime    int64
	Status           int
	Description      string
	Avatar           string
	Position         string
	Department       string
}

func (*MemberAccount) TableName() string {
	return "ms_member_account"
}

func (a *MemberAccount) ToDisplay() *MemberAccountDisplay {
	md := &MemberAccountDisplay{}
	copier.Copy(md, a)
	md.Code, _ = encrypts.EncryptInt64(a.Id, AESKey)
	md.MembarAccountCode = md.Code
	md.MemberCode, _ = encrypts.EncryptInt64(a.MemberCode, AESKey)
	md.OrganizationCode, _ = encrypts.EncryptInt64(a.OrganizationCode, AESKey)
	md.DepartmentCode, _ = encrypts.EncryptInt64(a.DepartmentCode, AESKey)
	md.CreateTime = tms.FormatByMill(a.CreateTime)
	md.LastLoginTime = tms.FormatByMill(a.LastLoginTime)
	md.StatusText = a.StatusText()
	md.AuthorizeArr = []string{a.Authorize}
	return md
}

func (a *MemberAccount) StatusText() string {
	if a.Status == 0 {
		return "禁用"
	}
	if a.Status == 1 {
		return "使用中"
	}
	return ""
}

type MemberAccountDisplay struct {
	Id                int64  `json:"id"`
	Code              string `json:"code"`
	MemberCode        string
	OrganizationCode  string   `json:"organization_code"`
	DepartmentCode    string   `json:"department_code"`
	Authorize         string   `json:"authorize"`
	IsOwner           int      `json:"is_owner"`
	Name              string   `json:"name"`
	Mobile            string   `json:"mobile"`
	Email             string   `json:"email"`
	CreateTime        string   `json:"create_time"`
	LastLoginTime     string   `json:"last_login_time"`
	Status            int      `json:"status"`
	Description       string   `json:"description"`
	Avatar            string   `json:"avatar"`
	Position          string   `json:"position"`
	Department        string   `json:"department"`
	MembarAccountCode string   `json:"membar_account_code"`
	Departments       string   `json:"departments"`
	StatusText        string   `json:"statusText"`
	AuthorizeArr      []string `json:"authorizeArr"`
}
