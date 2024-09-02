package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/status"
	"ms_project/project-api/api/rpc"
	"ms_project/project-api/pkg/model"
	common "ms_project/project-common"
	"ms_project/project-common/e"
	"ms_project/project-grpc/user/login"
	"net/http"
	"time"
)

type HandlerUser struct {
}

func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	result := &common.Response{}
	mobile := ctx.PostForm("mobile")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	captchaResponse, err := rpc.LoginServiceClient.GetCaptcha(c, &login.CaptchaMessage{Mobile: mobile})
	if err != nil {
		fromError, _ := status.FromError(err)
		ctx.JSON(http.StatusOK, result.Failed(int(fromError.Code()), fromError.Message()))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(captchaResponse.Code))
}

func (*HandlerUser) Register(c *gin.Context) {
	result := &common.Response{}
	//1.接收参数
	var RegisterService model.RegisterService
	err := c.ShouldBind(&RegisterService)
	if err != nil {
		c.JSON(http.StatusOK, result.Failed(http.StatusBadRequest, "参数有误"))
		return
	}
	//2.校验参数  参数是否合法
	if err = RegisterService.Verify(); err != nil {
		c.JSON(http.StatusOK, result.Failed(http.StatusBadRequest, "参数有误"))
		return
	}
	//3.调用user grpc服务 获取响应
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &login.RegisterMessage{
		Email:    RegisterService.Email,
		Name:     RegisterService.Name,
		Password: RegisterService.Password,
		Mobile:   RegisterService.Mobile,
		Captcha:  RegisterService.Captcha,
	}
	if _, err = rpc.LoginServiceClient.Register(ctx, msg); err != nil {
		fromError, _ := status.FromError(err)
		c.JSON(http.StatusOK, result.Failed(int(fromError.Code()), fromError.Message()))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(""))
}

func (*HandlerUser) Login(c *gin.Context) {
	result := &common.Response{}
	var loginReq model.LoginReq
	if err := c.ShouldBind(&loginReq); err != nil {
		c.JSON(http.StatusOK, result.Failed(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc模块
	//3.调用user grpc服务 获取响应
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &login.LoginMessage{
		Account:  loginReq.Account,
		Password: loginReq.Password,
	}
	memberRsp, err := rpc.LoginServiceClient.Login(ctx, msg)
	if err != nil {
		fromError, _ := status.FromError(err)
		c.JSON(http.StatusOK, result.Failed(int(fromError.Code()), fromError.Message()))
		return
	}
	rsp := &login.LoginResponse{
		Member:           memberRsp.Member,
		OrganizationList: memberRsp.OrganizationList,
		TokenList:        memberRsp.TokenList,
	}
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (*HandlerUser) MyOrgList(c *gin.Context) {
	var code int
	result := &common.Response{}
	memIdStr, _ := c.Get("memberId")
	memId := memIdStr.(int64)
	list, err := rpc.LoginServiceClient.MyOrgList(context.Background(), &login.UserMessage{MemId: memId})
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
		return
	}
	if list.OrganizationList == nil {
		c.JSON(http.StatusOK, result.Success(make([]*model.OrganizationList, 0)))
	}
	var orgs []*model.OrganizationList
	copier.Copy(&orgs, list.OrganizationList)
	c.JSON(http.StatusOK, result.Success(orgs))
}
