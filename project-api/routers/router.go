package routers

import (
	"github.com/gin-gonic/gin"
)

// todo 路由分组 ：将所有路由结合至一个routers包中进行统一接口分组 分离出路由与接口实现文件

type InterfaceRouter interface {
	Router(r *gin.Engine)
}
type RegisterRouter struct {
}

var routers []InterfaceRouter

func (*RegisterRouter) Router(ro InterfaceRouter, r *gin.Engine) {
	ro.Router(r)
}

func (rr *RegisterRouter) InitRouter(r *gin.Engine) {
	//rg := New()
	//rr.Router(&user.RouterUser{}, r)
	for _, ro := range routers {
		ro.Router(r)
	}
}
func Register(ro ...InterfaceRouter) {
	routers = append(routers, ro...)
}
