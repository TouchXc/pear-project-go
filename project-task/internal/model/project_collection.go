package model

type CollectionProject struct {
	Id          int64
	ProjectCode int64
	MemberCode  int64
	CreateTime  int64
}

func (*CollectionProject) TableName() string {
	return "ms_project_collection"
}
