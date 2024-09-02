package task

import (
	"github.com/gin-gonic/gin"
	common "ms_project/project-common"
	"ms_project/project-common/e"
	"net/http"
	"strings"
)

var ignores = []string{
	"project/login/register",
	"project/login",
	"project/login/getCaptcha",
	"project/organization",
	"project/auth/apply"}

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		var code int
		result := &common.Response{}
		uri := c.Request.RequestURI
		//判断此uri是否在用户授权列表中
		h := &HandlerTask{}
		nodes, err := h.GetAuthNodes(c)
		if err != nil {
			code = e.ParseGrpcError
			c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
			c.Abort()
			return
		}
		//不需要权限认证的接口忽略拦截
		for _, v := range ignores {
			if strings.Contains(uri, v) {
				c.Next()
				return
			}
		}
		//已有权限的接口可以访问
		for _, v := range nodes {
			if strings.Contains(uri, v) {
				c.Next()
				return
			}
		}
		//无权限接口进行拦截提示
		c.JSON(http.StatusOK, result.Failed(403, "无操作权限"))
		c.Abort()
	}
}
