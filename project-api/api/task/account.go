package task

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"ms_project/project-api/api/rpc"
	"ms_project/project-api/pkg/model"
	common "ms_project/project-common"
	"ms_project/project-common/e"
	"ms_project/project-grpc/account"
	"ms_project/project-grpc/task"
	"net/http"
	"time"
)

func (*HandlerTask) Account(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req model.AccountReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &account.AccountReqMessage{
		MemberId:         c.GetInt64("memberId"),
		OrganizationCode: c.GetString("organizationCode"),
		Page:             int64(req.Page),
		PageSize:         int64(req.PageSize),
		SearchType:       int32(req.SearchType),
		DepartmentCode:   req.DepartmentCode,
	}
	response, err := rpc.AccountServiceClient.Account(ctx, msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	//组装接口返回数据
	var list []*model.MemberAccount
	copier.Copy(&list, response.AccountList)
	if list == nil {
		list = []*model.MemberAccount{}
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.AuthList)
	if authList == nil {
		authList = []*model.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total":    response.Total,
		"page":     req.Page,
		"list":     list,
		"authList": authList,
	}))
}
func (*HandlerTask) ListDepartment(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req *model.DepartmentReq
	c.ShouldBind(&req)
	msg := &account.DepartmentReqMessage{
		Page:                 req.Page,
		PageSize:             req.PageSize,
		ParentDepartmentCode: req.Pcode,
		OrganizationCode:     c.GetString("organizationCode"),
	}
	listDepartmentMessage, err := rpc.AccountServiceClient.ListDepartment(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var list []*model.Department
	copier.Copy(&list, listDepartmentMessage.List)
	if list == nil {
		list = []*model.Department{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": listDepartmentMessage.Total,
		"page":  req.Page,
		"list":  list,
	}))
}

func (*HandlerTask) SaveDepartment(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req *model.DepartmentReq
	c.ShouldBind(&req)
	msg := &account.DepartmentReqMessage{
		Name:                 req.Name,
		DepartmentCode:       req.DepartmentCode,
		ParentDepartmentCode: req.ParentDepartmentCode,
		OrganizationCode:     c.GetString("organizationCode"),
	}
	departmentMessage, err := rpc.AccountServiceClient.SaveDepartment(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var res = &model.Department{}
	copier.Copy(res, departmentMessage)
	c.JSON(http.StatusOK, result.Success(res))
}
func (*HandlerTask) ReadDepartment(c *gin.Context) {
	var code int
	result := &common.Response{}
	departmentCode := c.PostForm("departmentCode")
	msg := &account.DepartmentReqMessage{
		DepartmentCode:   departmentCode,
		OrganizationCode: c.GetString("organizationCode"),
	}
	departmentMessage, err := rpc.AccountServiceClient.ReadDepartment(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var res = &model.Department{}
	copier.Copy(res, departmentMessage)
	c.JSON(http.StatusOK, result.Success(res))
}
func (*HandlerTask) AuthList(c *gin.Context) {
	var code int
	result := &common.Response{}
	organizationCode := c.GetString("organizationCode")
	var page = &model.Page{}
	page.Bind(c)
	msg := &account.AuthReqMessage{
		OrganizationCode: organizationCode,
		Page:             page.Page,
		PageSize:         page.PageSize,
	}
	response, err := rpc.AccountServiceClient.AuthList(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.List)
	if authList == nil {
		authList = []*model.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": response.Total,
		"list":  authList,
		"page":  page.Page,
	}))
}

func (*HandlerTask) MenuList(c *gin.Context) {
	var code int
	result := &common.Response{}
	msg := &account.MenuReqMessage{}
	res, err := rpc.AccountServiceClient.MenuList(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var list []*model.Menu
	copier.Copy(&list, res.List)
	if list == nil {
		list = []*model.Menu{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}

func (*HandlerTask) NodeList(c *gin.Context) {
	var code int
	result := &common.Response{}
	response, err := rpc.TaskServiceClient.NodeList(context.Background(), &task.ProjectRpcMessage{})
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var list []*model.ProjectNodeTree
	copier.Copy(&list, response.Nodes)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"nodes": list,
	}))
}
func (*HandlerTask) AuthApply(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req *model.ProjectAuthReq
	c.ShouldBind(&req)
	var nodes []string
	if req.Nodes != "" {
		json.Unmarshal([]byte(req.Nodes), &nodes)
	}
	msg := &account.AuthReqMessage{
		Action: req.Action,
		AuthId: req.Id,
		Nodes:  nodes,
	}
	applyResponse, err := rpc.AccountServiceClient.AuthApply(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var list []*model.ProjectNodeAuthTree
	copier.Copy(&list, applyResponse.List)
	var checkedList []string
	copier.Copy(&checkedList, applyResponse.CheckedList)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":        list,
		"checkedList": checkedList,
	}))
}
func (*HandlerTask) GetAuthNodes(c *gin.Context) ([]string, error) {
	memberId := c.GetInt64("memberId")
	msg := &account.AuthReqMessage{
		MemberId: memberId,
	}
	response, err := rpc.AccountServiceClient.AuthNodesByMemberId(context.Background(), msg)
	if err != nil {
		return nil, err
	}
	return response.List, err
}
