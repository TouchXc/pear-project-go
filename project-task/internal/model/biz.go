package model

var (
	Normal            = 1
	Personal    int32 = 1
	AESKey            = "sxdgtyuhbklouyrfdscgyuin"
	NoDeleted         = 0
	Deleted           = 1
	NoArchive         = 0
	Archive           = 1
	Open              = 0
	Private           = 1
	Custom            = 2
	Default           = "default"
	Simple            = "simple"
	NoCollected       = 0
	Collected         = 1
	IsOwner     int32 = 1
	NotOwner    int32 = 0
)

const (
	NoExecutor = iota
	Executor
)
const (
	NoOwner = iota
	Owner
)
const (
	NoCanRead = iota
	CanRead
)
const (
	NotComment = iota
	Comment
)
