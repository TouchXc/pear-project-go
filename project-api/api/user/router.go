package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"ms_project/project-api/api/rpc"
	"ms_project/project-api/middleware"
	"ms_project/project-api/routers"
)

type RouterUser struct {
}

func init() {
	log.Println("init user router")
	routerUser := &RouterUser{}
	routers.Register(routerUser)
}
func (*RouterUser) Router(r *gin.Engine) {
	//初始化grpc客户端连接
	rpc.InitRpcUserClient()
	h := &HandlerUser{}
	r.POST("project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/login/register", h.Register)
	r.POST("project/login", h.Login)
	org := r.Group("project/organization")
	org.Use(middleware.TokenVerify())
	{
		org.POST("/_getOrgList", h.MyOrgList)
	}
}
