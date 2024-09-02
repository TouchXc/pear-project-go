package model

type TaskMember struct {
	Id         int64
	TaskCode   int64
	IsExecutor int
	MemberCode int64
	JoinTime   int64
	IsOwner    int
}

func (*TaskMember) TableName() string {
	return "ms_task_member"
}
