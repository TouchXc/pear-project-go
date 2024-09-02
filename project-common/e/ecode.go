package e

const (

	//	redis/mysql数据库错误
	RedisError     = 10001
	MysqlError     = 10002
	ParseGrpcError = 10003
	// 用户登录服务错误
	InValidMobile          = 2001
	InValidCaptcha         = 2002
	EmailHasExist          = 2003
	UserHasExist           = 2004
	MobileHasExist         = 2005
	CaptchaExpired         = 2006
	AccountOrPasswordError = 2007
	TokenExpired           = 2008
	TaskNameNUllError      = 2009
	TaskStagesNullError    = 2010
	ProjectHasDeletedError = 2011
	Error                  = 500
)

var MsgFlags = map[int]string{
	InValidMobile:          "手机号不合法",
	InValidCaptcha:         "验证码错误",
	RedisError:             "redis 错误",
	MysqlError:             "数据库错误",
	ParseGrpcError:         "grpc解析token错误",
	EmailHasExist:          "邮箱已注册",
	UserHasExist:           "账号已注册",
	MobileHasExist:         "手机号已注册",
	CaptchaExpired:         "验证码已失效",
	AccountOrPasswordError: "账号密码错误",
	TokenExpired:           "Token已过期",
	Error:                  "fail",
	TaskNameNUllError:      "任务标题不能为空",
	TaskStagesNullError:    "任务步骤不能为空",
	ProjectHasDeletedError: "该项目已被删除",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[Error]
	}
	return msg
}
