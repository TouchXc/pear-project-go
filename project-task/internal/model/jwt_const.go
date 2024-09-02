package model

import "time"

// 设置jwt过期时间以及加密串
const (
	AccessSecret                = "msproject"
	RefreshSecret               = "ms_project"
	AccessExp     time.Duration = 7
	RefreshExp    time.Duration = 14
)
