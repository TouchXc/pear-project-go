package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"ms_project/project-api/api/rpc"
	common "ms_project/project-common"
	"ms_project/project-common/e"
	"ms_project/project-grpc/user/login"
	"net/http"
	"time"
)

func TokenVerify() func(c *gin.Context) {
	return func(c *gin.Context) {
		result := &common.Response{}
		//1.从header中获取token
		token := c.GetHeader("Authorization")
		//2.调用user服务进行token认证
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		mem, err := rpc.LoginServiceClient.TokenVerify(ctx, &login.LoginMessage{Token: token})
		if err != nil {
			code := e.ParseGrpcError
			c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
			c.Abort()
			return
		}
		//3.处理结果，认证通过 将信息放入gin上下文 失败返回未登录
		c.Set("memberId", mem.Member.Id)
		c.Set("memberName", mem.Member.Name)
		c.Set("organizationCode", mem.Member.OrganizationCode)
		c.Next()
	}
}
