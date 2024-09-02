package model

type Project struct {
	Id                 int64
	Cover              string
	Name               string
	Description        string
	AccessControlType  int
	WhiteList          string
	Sort               int
	Deleted            int
	TemplateCode       int
	Schedule           float64
	CreateTime         int64
	OrganizationCode   int64
	DeletedTime        int64
	Private            int
	Prefix             string
	OpenPrefix         int
	Archive            int
	ArchiveTime        int64
	OpenBeginTime      int
	OpenTaskPrivate    int
	TaskBoardTheme     string
	BeginTime          int64
	EndTime            int64
	AutoUpdateSchedule int
}

func (*Project) TableName() string {
	return "ms_project"
}

func ToProjectMap(list []*Project) map[int64]*Project {
	m := make(map[int64]*Project, len(list))
	for _, v := range list {
		m[v.Id] = v
	}
	return m
}

func (p *Project) GetAccessControlType() string {
	if p.AccessControlType == 0 {
		return "open"
	}
	if p.AccessControlType == 1 {
		return "private"
	}
	if p.AccessControlType == 2 {
		return "custom"
	}
	return ""
}
