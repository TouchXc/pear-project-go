package task

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ms_project/project-api/api/rpc"
	"ms_project/project-api/pkg/model"
	common "ms_project/project-common"
	"ms_project/project-common/e"
	"ms_project/project-common/mio"
	"ms_project/project-common/tms"
	"ms_project/project-grpc/task"
	"net/http"
	"path"
	"strconv"
	"time"
)

type HandlerTask struct {
}

func (*HandlerTask) Index(c *gin.Context) {
	result := &common.Response{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task.IndexMessage{}
	indexResponse, err := rpc.TaskServiceClient.Index(ctx, msg)
	if err != nil {
		code := e.Error
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	menus := indexResponse.Menus
	var ms []*model.Menu
	copier.Copy(&ms, menus)
	c.JSON(http.StatusOK, result.Success(ms))
}

func (*HandlerTask) MyProjectList(c *gin.Context) {
	result := &common.Response{}
	page := &model.Page{}
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	selectBy := c.PostForm("selectBy")
	organizationCode := c.GetString("organizationCode")
	page.Bind(c)
	msg := &task.ProjectRpcMessage{
		MemberId:         memberId,
		MemberName:       memberName,
		Page:             page.Page,
		PageSize:         page.PageSize,
		SelectBy:         selectBy,
		OrganizationCode: organizationCode,
	}
	myProjectResponse, err := rpc.TaskServiceClient.FindProjectByMemId(context.Background(), msg)
	if err != nil {
		code := e.Error
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	var pmm []*model.ProjectAndMember
	_ = copier.Copy(&pmm, myProjectResponse.Pm)
	if pmm == nil {
		pmm = make([]*model.ProjectAndMember, 0)
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pmm,
		"total": myProjectResponse.Total,
	}))
}

func (*HandlerTask) ProjectTemplate(c *gin.Context) {
	result := &common.Response{}
	page := &model.Page{}
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	viewTypeStr := c.PostForm("viewType")
	viewType, _ := strconv.ParseInt(viewTypeStr, 10, 64)
	organizationCode := c.GetString("organizationCode")
	_ = c.ShouldBind(&page)
	//page.Bind(c)
	msg := &task.ProjectRpcMessage{
		MemberId:         memberId,
		MemberName:       memberName,
		Page:             page.Page,
		PageSize:         page.PageSize,
		ViewType:         int32(viewType),
		OrganizationCode: organizationCode,
	}
	templateResponse, err := rpc.TaskServiceClient.FindProjectTemplate(context.Background(), msg)
	if err != nil {
		code := e.Error
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	var pms []*model.ProjectTemplate
	_ = copier.Copy(&pms, templateResponse.Ptm)
	if pms == nil {
		pms = make([]*model.ProjectTemplate, 0)
	}
	for _, v := range pms {
		if v.TaskStages == nil {
			v.TaskStages = make([]*model.TaskStagesOnlyName, 0)
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": templateResponse.Total,
	}))
}

func (*HandlerTask) ProjectSave(c *gin.Context) {
	var code int
	result := &common.Response{}
	memberId := c.GetInt64("memberId")
	organizationCode := c.GetString("organizationCode")
	var req *model.SaveProjectRequest
	err := c.ShouldBind(&req)
	msg := &task.ProjectRpcMessage{
		MemberId:         memberId,
		OrganizationCode: organizationCode,
		TemplateCode:     req.TemplateCode,
		Description:      req.Description,
		Name:             req.Name,
		Id:               int64(req.Id),
	}
	saveProject, err := rpc.TaskServiceClient.SaveProject(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	var rsp *model.SaveProject
	_ = copier.Copy(&rsp, saveProject)
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (*HandlerTask) ReadProject(c *gin.Context) {
	var code int
	result := &common.Response{}
	projectCode := c.PostForm("projectCode")
	memberId := c.GetInt64("memberId")
	detail, err := rpc.TaskServiceClient.FindProjectDetail(context.Background(), &task.ProjectRpcMessage{ProjectCode: projectCode, MemberId: memberId})
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	pd := &model.ProjectDetail{}
	copier.Copy(&pd, detail)
	c.JSON(http.StatusOK, result.Success(pd))
}

func (*HandlerTask) DeleteProject(c *gin.Context) {
	var code int
	result := &common.Response{}
	projectCode := c.PostForm("projectCode")
	_, err := rpc.TaskServiceClient.UpdateDeletedProject(context.Background(), &task.ProjectRpcMessage{ProjectCode: projectCode, Deleted: true})
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	c.JSON(http.StatusOK, result.Success([]string{}))
}

func (*HandlerTask) RecoveryProject(c *gin.Context) {
	var code int
	result := &common.Response{}
	projectCode := c.PostForm("projectCode")
	_, err := rpc.TaskServiceClient.UpdateDeletedProject(context.Background(), &task.ProjectRpcMessage{ProjectCode: projectCode, Deleted: false})
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	c.JSON(http.StatusOK, result.Success([]string{}))
}

func (*HandlerTask) CollectProject(c *gin.Context) {
	var code int
	var Msg string
	projectCode := c.PostForm("projectCode")
	collectType := c.PostForm("type")
	memberId := c.GetInt64("memberId")
	_, err := rpc.TaskServiceClient.UpdateCollectedProject(context.Background(), &task.ProjectRpcMessage{ProjectCode: projectCode, Collected: collectType, MemberId: memberId})
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	if collectType == "cancel" {
		Msg = "取消收藏成功"
	} else {
		Msg = "加入收藏成功"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  Msg,
		"data": []string{},
	})
}

func (*HandlerTask) EditProject(c *gin.Context) {
	var req *model.ProjectReq
	var code int
	memberId := c.GetInt64("memberId")
	projectCode := c.PostForm("projectCode")
	_ = c.ShouldBind(&req)
	msg := task.UpdateProjectMessage{}
	_ = copier.Copy(&msg, req)
	msg.MemberCode = memberId
	msg.ProjectCode = projectCode
	_, err := rpc.TaskServiceClient.UpdateProject(context.Background(), &msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": []string{},
	})
}

func (*HandlerTask) TaskStages(c *gin.Context) {
	var code int
	result := &common.Response{}
	projectCode := c.PostForm("projectCode")
	page := &model.Page{}
	page.Bind(c)
	msg := &task.TaskReqMessage{
		MemberId:    c.GetInt64("memberId"),
		ProjectCode: projectCode,
		Page:        page.Page,
		PageSize:    page.PageSize,
	}
	stages, err := rpc.TaskServiceClient.TaskStages(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	var tspList []*model.TaskStagesResp
	copier.Copy(&tspList, stages.List)
	if tspList == nil {
		tspList = make([]*model.TaskStagesResp, 0)
	}
	//给返回结构体附上默认值
	for _, v := range tspList {
		v.TasksLoading = true  //任务加载状态
		v.FixedCreator = false //添加任务按钮定位
		v.ShowTaskCard = false //是否显示创建卡片
		v.Tasks = []int{}
		v.DoneTasks = []int{}
		v.UnDoneTasks = []int{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": stages.Total,
		"page":  stages.Total,
		"list":  tspList,
	}))
}

func (*HandlerTask) MemberProjectList(c *gin.Context) {
	var code int
	result := &common.Response{}
	projectCode := c.PostForm("projectCode")
	page := &model.Page{}
	page.Bind(c)
	msg := &task.TaskReqMessage{
		MemberId:    c.GetInt64("memberId"),
		ProjectCode: projectCode,
		Page:        page.Page,
		PageSize:    page.PageSize,
	}
	memberProjectList, err := rpc.TaskServiceClient.MemberProjectList(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	var list []*model.MemberProjectResp
	_ = copier.Copy(&list, memberProjectList.List)
	if list == nil {
		list = make([]*model.MemberProjectResp, 0)
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": memberProjectList.Total,
		"page":  page.Page,
		"list":  list,
	}))
}

func (*HandlerTask) TaskList(c *gin.Context) {
	var code int
	result := &common.Response{}
	stageCodeStr := c.PostForm("stageCode")
	list, err := rpc.TaskServiceClient.TaskList(context.Background(), &task.TaskReqMessage{StageCode: stageCodeStr, MemberId: c.GetInt64("memberId")})
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, status.Error(codes.Code(code), e.GetMsg(code)))
	}
	var taskDisplayList []*model.TaskDisplay
	_ = copier.Copy(&taskDisplayList, list.List)
	if taskDisplayList == nil {
		taskDisplayList = make([]*model.TaskDisplay, 0)
	}
	for _, v := range taskDisplayList {
		if v.Tags == nil {
			v.Tags = []int{}
		}
		if v.ChildCount == nil {
			v.ChildCount = []int{}
		}
	}
	c.JSON(http.StatusOK, result.Success(taskDisplayList))
}

func (*HandlerTask) SaveTask(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req *model.TaskSaveReq
	_ = c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task.TaskReqMessage{
		ProjectCode: req.ProjectCode,
		Name:        req.Name,
		StageCode:   req.StageCode,
		AssignTo:    req.AssignTo,
		MemberId:    c.GetInt64("memberId"),
	}
	taskMessage, err := rpc.TaskServiceClient.SaveTask(ctx, msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	td := &model.TaskDisplay{}
	_ = copier.Copy(td, taskMessage)
	if td != nil {
		if td.Tags == nil {
			td.Tags = []int{}
		}
		if td.ChildCount == nil {
			td.ChildCount = []int{}
		}
	}
	c.JSON(http.StatusOK, result.Success(td))
}

func (*HandlerTask) TaskSort(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req *model.TaskSortReq
	_ = c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task.TaskReqMessage{
		PreTaskCode:  req.PreTaskCode,
		NextTaskCode: req.NextTaskCode,
		ToStageCode:  req.ToStageCode,
	}
	_, err := rpc.TaskServiceClient.TaskSort(ctx, msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
		return
	}
	c.JSON(http.StatusOK, result.Success([]string{}))
}

func (*HandlerTask) MyTaskList(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req *model.MyTaskReq
	_ = c.ShouldBind(&req)
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task.TaskReqMessage{
		MemberId: memberId,
		TaskType: int32(req.TaskType),
		Type:     int32(req.Type),
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	myTaskListResponse, err := rpc.TaskServiceClient.MyTaskList(ctx, msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
		return
	}
	var myTaskList []*model.MyTaskDisplay
	copier.Copy(&myTaskList, myTaskListResponse.List)
	if myTaskList == nil {
		myTaskList = make([]*model.MyTaskDisplay, 0)
	}
	for _, v := range myTaskList {
		v.ProjectInfo = model.ProjectInfo{
			Name: v.ProjectName,
			Code: v.ProjectCode,
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  myTaskList,
		"total": myTaskListResponse.Total,
	}))
}

func (*HandlerTask) ReadTask(c *gin.Context) {
	var code int
	result := &common.Response{}
	taskCode := c.PostForm("taskCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
	}
	taskMessage, err := rpc.TaskServiceClient.ReadTask(ctx, msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	td := &model.TaskDisplay{}
	copier.Copy(td, taskMessage)
	if td != nil {
		if td.Tags == nil {
			td.Tags = []int{}
		}
		if td.ChildCount == nil {
			td.ChildCount = []int{}
		}
	}
	c.JSON(200, result.Success(td))
}

func (*HandlerTask) ListTaskMember(c *gin.Context) {
	var code int
	result := &common.Response{}
	taskCode := c.PostForm("taskCode")
	page := &model.Page{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task.TaskReqMessage{
		Page:     page.Page,
		PageSize: page.PageSize,
		MemberId: c.GetInt64("memberId"),
		TaskCode: taskCode,
	}
	taskMemberResponse, err := rpc.TaskServiceClient.ListTaskMember(ctx, msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var tms []*model.TaskMember
	copier.Copy(&tms, taskMemberResponse)
	if tms == nil {
		tms = []*model.TaskMember{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  tms,
		"total": taskMemberResponse.Total,
		"page":  page.Page,
	}))
}

func (*HandlerTask) TaskLog(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req *model.TaskLogReq
	c.ShouldBind(&req)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	msg := &task.TaskReqMessage{
		TaskCode: req.TaskCode,
		MemberId: c.GetInt64("memberId"),
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		All:      int32(req.All),
		Comment:  int32(req.Comment),
	}
	taskLogResponse, err := rpc.TaskServiceClient.TaskLog(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var tms []*model.ProjectLogDisplay
	copier.Copy(&tms, taskLogResponse.List)
	if tms == nil {
		tms = []*model.ProjectLogDisplay{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  tms,
		"total": taskLogResponse.Total,
		"page":  req.Page,
	}))
}

func (*HandlerTask) TaskWorkTimeList(c *gin.Context) {
	var code int
	taskCode := c.PostForm("taskCode")
	result := &common.Response{}
	msg := &task.TaskReqMessage{
		TaskCode: taskCode,
		MemberId: c.GetInt64("memberId"),
	}
	taskWorkTimeResponse, err := rpc.TaskServiceClient.TaskWorkTimeList(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var tms []*model.TaskWorkTime
	copier.Copy(&tms, taskWorkTimeResponse.List)
	if tms == nil {
		tms = []*model.TaskWorkTime{}
	}
	c.JSON(http.StatusOK, result.Success(tms))
}

func (*HandlerTask) SaveTaskWorkTime(c *gin.Context) {
	var code int
	result := &common.Response{}
	var req *model.SaveTaskWorkTimeReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &task.TaskReqMessage{
		TaskCode:  req.TaskCode,
		MemberId:  c.GetInt64("memberId"),
		Content:   req.Content,
		Num:       int32(req.Num),
		BeginTime: tms.ParseTime(req.BeginTime),
	}
	_, err := rpc.TaskServiceClient.SaveTaskWorkTime(ctx, msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

//文件上传接口

func (*HandlerTask) UploadFiles(c *gin.Context) {
	result := common.Response{}
	var code int
	var err error
	var req *model.UploadFileReq
	c.ShouldBind(&req)
	//处理文件
	multiFile, _ := c.MultipartForm()
	file := multiFile.File
	//假设只上传一个文件
	uploadFile := file["file"][0]
	//情况1  文件无需分片
	key := "msproject" + req.Filename
	minioClient, err := mio.NewMinioClient(
		"localhost:9009",
		"pFC9fBQycWspRtsOTs4G",
		"UWa5RuLDpQZf1HAhv1Mn1qA7z1ZgD4DNFUbzmj7Q",
		false)
	if err != nil {
		code = e.Error
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
		return
	}
	bucketName := "msproject"
	if req.TotalChunks == 1 {
		//代表不分片，直接上传
		//path := "upload/" + req.ProjectCode + "/" + req.TaskCode + "/" + tms.FormatYMD(time.Now())
		//if !fs.IsExist(path) {
		//	os.MkdirAll(path, os.ModePerm)
		//}
		//dst := path + "/" + req.Filename
		//key = dst
		//header := file["file"][0]
		//err := c.SaveUploadedFile(header, dst)
		//if err != nil {
		//	c.JSON(http.StatusOK, result.Failed(-999, err.Error()))
		//	return
		//}
		//修改后：上传至Minio服务
		//open, err := file["file"][0].Open()
		open, err := uploadFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 300,
				"data": []string{},
				"msg":  "文件上传出错",
			})
			return
		}
		defer open.Close()
		buf := make([]byte, req.CurrentChunkSize)
		open.Read(buf)
		_, err = minioClient.Put(context.Background(), bucketName, req.Filename, buf, int64(req.TotalSize), uploadFile.Header.Get("Content-Type"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 300,
				"data": []string{},
				"msg":  "文件上传出错",
			})
			return
		}
	}
	if req.TotalChunks > 1 {
		//分片上传 先把每次存储起来，再追加上传
		//path := "upload/" + req.ProjectCode + "/" + req.TaskCode + "/" + tms.FormatYMD(time.Now())
		//if !fs.IsExist(path) {
		//	os.MkdirAll(path, os.ModePerm)
		//}
		//fileName := path + "/" + req.Identifier
		//openFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		//if err != nil {
		//	c.JSON(http.StatusOK, result.Failed(-999, err.Error()))
		//	return
		//}
		open, err := uploadFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 300,
				"data": []string{},
				"msg":  "文件上传出错",
			})
			return
		}
		defer open.Close()
		buf := make([]byte, req.CurrentChunkSize)
		open.Read(buf)
		formatInt := strconv.FormatInt(int64(req.ChunkNumber), 10)
		//openFile.Write(buf)
		//openFile.Close()
		//newpath := path + "/" + req.Filename
		//key = newpath
		//改造后：使用minio分片上传
		_, err = minioClient.Put(context.Background(), bucketName, req.Filename+formatInt, buf, int64(req.CurrentChunkSize), uploadFile.Header.Get("Content-Type"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 300,
				"data": []string{},
				"msg":  "文件上传出错",
			})
			return
		}
		if req.TotalChunks == req.ChunkNumber {
			//最后一块 重命名文件名
			//err = os.Rename(fileName, newpath)
			//fmt.Println(err)
			//使用minio改造后：上传最后一个分片文件后，合并全部文件
			_, err = minioClient.Compose(context.Background(), bucketName, req.Filename, req.TotalChunks)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 300,
					"data": []string{},
					"msg":  "文件上传出错",
				})
				return
			}
		}
	}
	//调用服务 存入file表
	fileUrl := "http://localhost:9009/" + key
	msg := &task.TaskFileReqMessage{
		TaskCode:         req.TaskCode,
		ProjectCode:      req.ProjectCode,
		OrganizationCode: c.GetString("organizationCode"),
		PathName:         key,
		FileName:         req.Filename,
		Size:             int64(req.TotalSize),
		Extension:        path.Ext(key),
		FileUrl:          fileUrl,
		FileType:         file["file"][0].Header.Get("Content-Type"),
		MemberId:         c.GetInt64("memberId"),
	}
	if req.TotalChunks == req.ChunkNumber {
		_, err = rpc.TaskServiceClient.SaveTaskFile(context.Background(), msg)
		if err != nil {
			code = e.ParseGrpcError
			c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"file":        key,
		"hash":        "",
		"key":         key,
		"url":         "http://localhost:9009/" + key,
		"projectName": req.ProjectName,
	}))
	return
}

func (*HandlerTask) TaskSource(c *gin.Context) {
	var code int
	result := &common.Response{}
	taskCode := c.PostForm("taskCode")
	sources, err := rpc.TaskServiceClient.TaskSources(context.Background(), &task.TaskReqMessage{TaskCode: taskCode})
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var slList []*model.SourceLink
	copier.Copy(&slList, sources.List)
	if slList == nil {
		slList = []*model.SourceLink{}
	}
	c.JSON(http.StatusOK, result.Success(slList))
}

func (*HandlerTask) CreateComment(c *gin.Context) {
	var code int
	result := &common.Response{}
	req := model.CommentReq{}
	c.ShouldBind(&req)
	msg := &task.TaskReqMessage{
		TaskCode:       req.TaskCode,
		CommentContent: req.Comment,
		Mentions:       req.Mentions,
		MemberId:       c.GetInt64("memberId"),
	}
	_, err := rpc.TaskServiceClient.CreateComment(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	c.JSON(http.StatusOK, result.Success(true))
}

func (*HandlerTask) GetLogBySelfProject(c *gin.Context) {
	result := &common.Response{}
	var code int
	var page = &model.Page{}
	page.Bind(c)
	msg := &task.ProjectRpcMessage{
		MemberId: c.GetInt64("memberId"),
		Page:     page.Page,
		PageSize: page.PageSize,
	}
	projectLogResponse, err := rpc.TaskServiceClient.GetLogBySelfProject(context.Background(), msg)
	if err != nil {
		code = e.ParseGrpcError
		c.JSON(http.StatusOK, result.Failed(code, e.GetMsg(code)))
	}
	var list []*model.ProjectLog
	copier.Copy(&list, projectLogResponse.List)
	if list == nil {
		list = []*model.ProjectLog{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}
