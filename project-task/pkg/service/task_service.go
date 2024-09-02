package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ms_project/project-common/e"
	"ms_project/project-common/encrypts"
	"ms_project/project-common/tms"
	"ms_project/project-grpc/task"
	"ms_project/project-grpc/user/login"
	"ms_project/project-task/internal/dao"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
	"ms_project/project-task/internal/repo"
	"ms_project/project-task/internal/rpc"
	"ms_project/project-task/kafka"
	"strconv"
	"time"
)

type TaskService struct {
	task.UnimplementedTaskServiceServer
	cache                  repo.Cache
	menuRepo               repo.MenuRepo
	projectRepo            repo.ProjectRepo
	ProjectTemplateRepo    repo.ProjectTemplateRepo
	taskStagesTemplateRepo repo.TaskStagesTemplateRepo
	taskStagesRepo         repo.TaskStagesRepo
	taskRepo               repo.TaskRepo
	projectLogRepo         repo.ProjectLogRepo
	fileRepo               repo.FileRepo
	sourceLinkRepo         repo.SourceLinkRepo
}

func NewTaskService() *TaskService {
	return &TaskService{
		cache:                  dao.RC,
		menuRepo:               dao.NewMenuDao(),
		projectRepo:            dao.NewProjectDao(),
		ProjectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
		taskStagesRepo:         dao.NewTaskStagesDao(),
		taskRepo:               dao.NewTaskDao(),
		projectLogRepo:         dao.NewProjectLogDao(),
		fileRepo:               dao.NewFileDao(),
		sourceLinkRepo:         dao.NewSourceLinkDao(),
	}
}

func (ts *TaskService) Index(ctx context.Context, msg *task.IndexMessage) (*task.IndexResponse, error) {
	var code int
	pms, err := ts.menuRepo.FindMenus(context.Background())
	if err != nil {
		code = e.MysqlError
		zap.L().Error("FindMenus db mysql error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	childs := model.CovertChild(pms)
	var mms []*task.MenuMessage
	copier.Copy(&mms, childs)
	return &task.IndexResponse{Menus: mms}, nil
}
func (ts *TaskService) FindProjectByMemId(ctx context.Context, msg *task.ProjectRpcMessage) (*task.MyProjectResponse, error) {
	var code int
	var err error
	var total int64
	var pms []*model.ProjectAndMember
	//TODO() 优化代码
	if msg.SelectBy == "" || msg.SelectBy == "my" {
		pms, total, err = ts.projectRepo.FindProjectByMemberId(ctx, msg.MemberId, "and deleted = 0", msg.Page, msg.PageSize)
	}
	if msg.SelectBy == "archive" {
		pms, total, err = ts.projectRepo.FindProjectByMemberId(ctx, msg.MemberId, "and archive = 1", msg.Page, msg.PageSize)
	}
	if msg.SelectBy == "deleted" {
		pms, total, err = ts.projectRepo.FindProjectByMemberId(ctx, msg.MemberId, "and deleted = 1", msg.Page, msg.PageSize)
	}
	if msg.SelectBy == "collect" {
		pms, total, err = ts.projectRepo.FindCollectProjectByMemberId(ctx, msg.MemberId, msg.Page, msg.PageSize)
		for _, v := range pms {
			v.Collected = model.Collected
		}
	} else {
		collectPms, _, err := ts.projectRepo.FindCollectProjectByMemberId(ctx, msg.MemberId, msg.Page, msg.PageSize)
		if err != nil {
			code = e.MysqlError
			zap.L().Error("FindProjectByMemId FindCollectProjectByMemberId error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
		cMap := make(map[int64]*model.ProjectAndMember)
		for _, v := range collectPms {
			cMap[v.Id] = v
		}
		for _, v := range pms {
			if cMap[v.ProjectCode] != nil {
				v.Collected = model.Collected
			}
		}
	}
	//switch msg.SelectBy {
	//case "":
	//	pms, total, err = ts.projectRepo.FindProjectByMemberId(ctx, msg.MemberId, "and deleted = 0", msg.Page, msg.PageSize)
	//case "my":
	//	pms, total, err = ts.projectRepo.FindProjectByMemberId(ctx, msg.MemberId, "and deleted = 0", msg.Page, msg.PageSize)
	//case "archive":
	//	pms, total, err = ts.projectRepo.FindProjectByMemberId(ctx, msg.MemberId, "and archive = 1", msg.Page, msg.PageSize)
	//case "deleted":
	//	pms, total, err = ts.projectRepo.FindProjectByMemberId(ctx, msg.MemberId, "and deleted = 1", msg.Page, msg.PageSize)
	//case "collect":
	//	pms, total, err = ts.projectRepo.FindCollectProjectByMemberId(ctx, msg.MemberId, msg.Page, msg.PageSize)
	//	for _, v := range pms {
	//		v.Collected = model.Collected
	//	}
	//default:
	//
	//}
	if err != nil {
		code = e.MysqlError
		zap.L().Error("Menu Find all error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if pms == nil {
		return &task.MyProjectResponse{Pm: make([]*task.ProjectMessage, 0), Total: total}, nil
	}
	var pmm []*task.ProjectMessage
	copier.Copy(&pmm, pms)
	for _, v := range pmm {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
		pam := model.ToMap(pms)[v.Id]
		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, model.AESKey)
		v.JoinTime = tms.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Sort)
		v.CreateTime = tms.FormatByMill(pam.CreateTime)
	}
	return &task.MyProjectResponse{Pm: pmm, Total: total}, nil
}
func (ts *TaskService) FindProjectTemplate(ctx context.Context, msg *task.ProjectRpcMessage) (*task.ProjectTemplateResponse, error) {
	var pts []model.ProjectTemplate
	var code int
	var total int64
	var err error
	//1.根据viewType查询项目模板 得到list
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	page := msg.Page
	pageSize := msg.PageSize
	if msg.ViewType == -1 {
		pts, total, err = ts.ProjectTemplateRepo.FindProjectTemplateAll(ctx, organizationCode, page, pageSize)
	}
	if msg.ViewType == 0 {
		pts, total, err = ts.ProjectTemplateRepo.FindProjectTemplateCustom(ctx, msg.MemberId, organizationCode, page, pageSize)
	}
	if msg.ViewType == 1 {
		pts, total, err = ts.ProjectTemplateRepo.FindProjectTemplateSystem(ctx, page, pageSize)
	}
	//switch msg.ViewType {
	//case -1:
	//	pts, total, err = ts.ProjectTemplateRepo.FindProjectTemplateAll(ctx, organizationCode, page, pageSize)
	//case 0:
	//	pts, total, err = ts.ProjectTemplateRepo.FindProjectTemplateCustom(ctx, msg.MemberId, organizationCode, page, pageSize)
	//case 1:
	//	pts, total, err = ts.ProjectTemplateRepo.FindProjectTemplateSystem(ctx, page, pageSize)
	//}

	//	2. 模型转换，拿到模板id列表 去任务步骤模板表 查询
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" FindProjectTemplate error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	tsts, err := ts.taskStagesTemplateRepo.FindInProTemIds(ctx, model.ToProjectTemplateIds(pts))
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" FindProjectTemplate FindInProTemIds error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	var ptas []*model.ProjectTemplateAll
	for _, v := range pts {
		ptas = append(ptas, v.Convert(model.CovertProjectMap(tsts)[v.Id]))
	}
	//	3.组装数据
	var pmMsgs []*task.ProjectTemplateMessage
	copier.Copy(&pmMsgs, ptas)
	return &task.ProjectTemplateResponse{Ptm: pmMsgs, Total: total}, nil
}
func (ts *TaskService) SaveProject(ctx context.Context, msg *task.ProjectRpcMessage) (*task.SaveProjectMessage, error) {
	var code int
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	templateCodeStr, _ := encrypts.Decrypt(msg.TemplateCode, model.AESKey)
	templateCode, _ := strconv.ParseInt(templateCodeStr, 10, 64)
	//获取模板信息
	stagesTemplateList, err := ts.taskStagesTemplateRepo.FindStagesByProjectTemplateCode(ctx, int(templateCode))
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" SaveProject FindStagesByProjectTemplateCode error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//1.保存项目表
	pr := &model.Project{
		Name:              msg.Name,
		Description:       msg.Description,
		TemplateCode:      int(templateCode),
		CreateTime:        time.Now().UnixMilli(),
		Cover:             "https://img2.baidu.com/it/u=792555388,2449797505&fm=253&fmt=auto&app=138&f=JPEG?w=667&h=500",
		Deleted:           model.NoDeleted,
		Archive:           model.NoArchive,
		OrganizationCode:  organizationCode,
		AccessControlType: model.Open,
		TaskBoardTheme:    model.Simple,
	}
	var rsp *task.SaveProjectMessage
	db := gorms.NewDBClient()
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(&pr).Error; err != nil {
			zap.L().Error("Transaction SaveProject error", zap.Error(err))
			return err
		}
		//2.保存项目和成员关联表
		pm := &model.MemberProject{
			ProjectCode: pr.Id,
			MemberCode:  msg.MemberId,
			JoinTime:    time.Now().UnixMilli(),
			IsOwner:     msg.MemberId,
			Authorize:   "",
		}
		if err = tx.Save(&pm).Error; err != nil {
			zap.L().Error("Transaction SaveProjectMember error", zap.Error(err))
			return err
		}
		//3.生成任务步骤
		for index, v := range stagesTemplateList {
			taskStage := &model.TaskStages{
				Name:        v.Name,
				ProjectCode: pr.Id,
				Sort:        index + 1,
				Description: "",
				CreateTime:  time.Now().UnixMilli(),
				Deleted:     model.NoDeleted,
			}
			if err = tx.Save(&taskStage).Error; err != nil {
				zap.L().Error("Transaction SaveProjectMember error", zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

//1.查项目详情 2.项目和成员关联表 查项目拥有者 去member表查名字 3.查收藏表

func (ts *TaskService) FindProjectDetail(ctx context.Context, msg *task.ProjectRpcMessage) (*task.ProjectDetailMessage, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	memberId := msg.MemberId
	projectAndMember, err := ts.projectRepo.FindProjectByPidAndMemberId(ctx, projectCode, memberId)
	if err != nil {
		code := e.MysqlError
		zap.L().Error(" FindProjectDetail FindProjectByPidAndMemberId error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	ownerId := projectAndMember.IsOwner
	//去user模块找
	memInfoById, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: ownerId})
	if memInfoById == nil {
		memInfoById = &login.MemberMessage{}
	}
	if err != nil {
		zap.L().Error(" FindProjectDetail FindMemInfoById error", zap.Error(err))
		return nil, err
	}
	// TODO 优化 收藏时放入redis
	isCollected, err := ts.projectRepo.FindCollectProjectByPidAndMemberId(context.Background(), projectCode, memberId)
	if err != nil {
		code := e.MysqlError
		zap.L().Error(" FindProjectDetail FindCollectProjectByPidAndMemberId error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if isCollected {
		projectAndMember.Collected = model.Collected
	}
	var detailMsg = &task.ProjectDetailMessage{}
	copier.Copy(&detailMsg, projectAndMember)
	detailMsg.OwnerAvatar = memInfoById.Avatar
	detailMsg.OwnerName = memInfoById.Name
	detailMsg.Code, _ = encrypts.EncryptInt64(projectAndMember.Id, model.AESKey)
	detailMsg.AccessControlType = projectAndMember.GetAccessControlType()
	detailMsg.OrganizationCode, _ = encrypts.EncryptInt64(projectAndMember.OrganizationCode, model.AESKey)
	detailMsg.Order = int32(projectAndMember.Sort)
	detailMsg.CreateTime = tms.FormatByMill(projectAndMember.CreateTime)
	return detailMsg, nil
}
func (ts *TaskService) UpdateDeletedProject(ctx context.Context, msg *task.ProjectRpcMessage) (*task.DeletedProjectResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	err := ts.projectRepo.UpdateDeleteProject(ctx, projectCode, msg.Deleted)
	if err != nil {
		code := e.MysqlError
		zap.L().Error("RecycleProject DeleteProject error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	return &task.DeletedProjectResponse{}, nil
}
func (ts *TaskService) UpdateCollectedProject(ctx context.Context, msg *task.ProjectRpcMessage) (*task.CollectedProjectResponse, error) {
	var err error
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	if msg.Collected == "collect" {
		pc := &model.CollectionProject{
			ProjectCode: projectCode,
			MemberCode:  msg.MemberId,
			CreateTime:  time.Now().UnixMilli(),
		}
		err = ts.projectRepo.SaveProjectCollect(ctx, pc)
	} else if msg.Collected == "cancel" {
		err = ts.projectRepo.DeleteProjectCollect(ctx, msg.MemberId, projectCode)
	}
	if err != nil {
		code := e.MysqlError
		zap.L().Error("UpdateCollectedProject UpdateCollectedProject error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	return &task.CollectedProjectResponse{}, nil
}
func (ts *TaskService) UpdateProject(ctx context.Context, msg *task.UpdateProjectMessage) (*task.UpdateProjectResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	project := &model.Project{
		Id:                 projectCode,
		Name:               msg.Name,
		Description:        msg.Description,
		Cover:              msg.Cover,
		TaskBoardTheme:     msg.TaskBoardTheme,
		Private:            int(msg.Private),
		Prefix:             msg.Prefix,
		OpenPrefix:         int(msg.OpenPrefix),
		OpenBeginTime:      int(msg.OpenBeginTime),
		OpenTaskPrivate:    int(msg.OpenTaskPrivate),
		Schedule:           msg.Schedule,
		AutoUpdateSchedule: int(msg.AutoUpdateSchedule),
	}
	err := ts.projectRepo.UpdateProject(ctx, project)
	if err != nil {
		code := e.MysqlError
		zap.L().Error("UpdateProject UpdateProject error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	return &task.UpdateProjectResponse{}, nil
}
func (ts *TaskService) TaskStages(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskStagesResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	page := msg.Page
	pageSize := msg.PageSize
	stages, total, err := ts.taskStagesRepo.FindStagesByProjectId(ctx, projectCode, page, pageSize)
	if err != nil {
		code := e.MysqlError
		zap.L().Error("TaskStages FindStagesByProjectId error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}

	var tsMessage []*task.TaskStagesMessage
	copier.Copy(&tsMessage, stages)
	if tsMessage == nil {
		return &task.TaskStagesResponse{List: tsMessage, Total: 0}, nil
	}
	stagesMap := model.ToTaskStagesMap(stages)
	for _, v := range tsMessage {
		taskStages := stagesMap[int(v.Id)]
		v.Code, _ = encrypts.EncryptInt64(int64(v.Id), model.AESKey)
		v.CreateTime = tms.FormatByMill(taskStages.CreateTime)
		v.ProjectCode = msg.ProjectCode
	}
	return &task.TaskStagesResponse{
		Total: total,
		List:  tsMessage,
	}, nil
}
func (ts *TaskService) MemberProjectList(ctx context.Context, msg *task.TaskReqMessage) (*task.MemberProjectResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	ProjectMemberList, total, err := ts.projectRepo.FindProjectByPid(context.Background(), projectCode)
	if err != nil {
		code := e.MysqlError
		zap.L().Error("MemberProjectList FindProjectByPid error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if ProjectMemberList == nil || len(ProjectMemberList) <= 0 {
		return &task.MemberProjectResponse{
			Total: 0,
			List:  nil,
		}, nil
	}
	var mIds []int64
	pmMap := make(map[int64]*model.MemberProject)
	for _, v := range ProjectMemberList {
		mIds = append(mIds, v.MemberCode)
		pmMap[v.MemberCode] = v
	}
	memberMessageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIds})
	if err != nil {
		code := e.MysqlError
		zap.L().Error("MemberProjectList FindMemInfoByIds error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	var list []*task.MemberProjectMessage
	for _, v := range memberMessageList.List {
		owner := pmMap[v.Id].IsOwner
		mpm := &task.MemberProjectMessage{
			Name:       v.Name,
			Avatar:     v.Avatar,
			MemberCode: v.Id,
			Email:      v.Email,
			Code:       v.Code,
		}
		if v.Id == owner {
			mpm.IsOwner = model.IsOwner
		} else {
			mpm.IsOwner = model.NotOwner
		}
		list = append(list, mpm)
	}
	return &task.MemberProjectResponse{
		Total: total,
		List:  list,
	}, nil
}
func (ts *TaskService) TaskList(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskListResponse, error) {
	stageCodeStr, _ := encrypts.Decrypt(msg.StageCode, model.AESKey)
	stageCode, _ := strconv.ParseInt(stageCodeStr, 10, 64)
	taskList, err := ts.taskRepo.FindTaskByStageCode(ctx, stageCode)
	if err != nil {
		code := e.MysqlError
		zap.L().Error("TaskList FindTaskByStageCode error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	var taskDisplayList []*model.TaskDisplay
	var mIds []int64
	for _, v := range taskList {
		display := v.ToTaskDisplay()
		if v.Private == 1 {
			//代表隐私模式
			taskMember, err := ts.taskRepo.FindTaskMemberByTaskId(ctx, v.Id, msg.MemberId)
			if err != nil {
				code := e.MysqlError
				zap.L().Error("TaskList FindTaskMemberByTaskId error", zap.Error(err))
				return nil, status.Error(codes.Code(code), e.GetMsg(code))
			}
			if taskMember != nil {
				display.CanRead = model.CanRead
			} else {
				display.CanRead = model.NoCanRead
			}
		}
		taskDisplayList = append(taskDisplayList, display)
		mIds = append(mIds, v.AssignTo)
	}
	if len(mIds) <= 0 || mIds == nil {
		return &task.TaskListResponse{List: nil}, nil
	}
	memInfoList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIds})
	if err != nil {
		zap.L().Error("TaskList FindMemInfoByIds error", zap.Error(err))
		return nil, err
	}

	memberMap := make(map[int64]*login.MemberMessage)
	for _, v := range memInfoList.List {
		memberMap[v.Id] = v
	}
	for _, v := range taskDisplayList {
		assTo, _ := encrypts.Decrypt(v.AssignTo, model.AESKey)
		id, _ := strconv.ParseInt(assTo, 10, 64)
		message := memberMap[id]
		ex := model.Executors{
			Name:   message.Name,
			Avatar: message.Avatar,
		}
		v.Executor = ex
		fmt.Println(message)
	}
	var taskListRsp []*task.TaskMessage
	_ = copier.Copy(&taskListRsp, taskDisplayList)
	return &task.TaskListResponse{List: taskListRsp}, nil
}
func (ts *TaskService) SaveTask(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	var code int
	//1.检查业务数据
	if msg.Name == "" {
		code = e.TaskNameNUllError
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//taskStages 检查任务步骤是否存在
	stageCodeStr, _ := encrypts.Decrypt(msg.StageCode, model.AESKey)
	stageCode, _ := strconv.ParseInt(stageCodeStr, 10, 64)
	taskStages, err := ts.taskStagesRepo.FindById(ctx, stageCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("SaveTask FindById error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if taskStages == nil {
		code = e.TaskStagesNullError
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//project 检查对应项目是否存在
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	project, err := ts.projectRepo.FindProjectById(ctx, projectCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("SaveTask projectRepo.FindProjectById error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if project.Deleted == model.Deleted {
		code = e.ProjectHasDeletedError
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	maxIdNum, err := ts.taskRepo.FindTaskMaxIdNum(ctx, projectCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task SaveTask taskRepo.FindTaskMaxIdNum error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	maxSort, err := ts.taskRepo.FindTaskSort(ctx, projectCode, stageCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task SaveTask taskRepo.FindTaskSort error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	assignToStr, _ := encrypts.Decrypt(msg.AssignTo, model.AESKey)
	assignTo, _ := strconv.ParseInt(assignToStr, 10, 64)
	t := &model.Task{
		Name:        msg.Name,
		CreateTime:  time.Now().UnixMilli(),
		CreateBy:    msg.MemberId,
		AssignTo:    assignTo,
		ProjectCode: projectCode,
		StageCode:   int(stageCode),
		IdNum:       maxIdNum + 1,
		Private:     project.OpenTaskPrivate,
		Sort:        maxSort + 1,
		BeginTime:   time.Now().UnixMilli(),
		EndTime:     time.Now().Add(2 * 24 * time.Hour).UnixMilli(),
	}
	db := gorms.NewDBClient()
	//开启事务
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&model.Task{}).Save(t).Error; err != nil {
			zap.L().Error("Transaction Save Task mysql error", zap.Error(err))
			return err
		}
		tm := &model.TaskMember{
			TaskCode:   t.Id,
			MemberCode: assignTo,
			JoinTime:   time.Now().UnixMilli(),
			IsOwner:    model.Owner,
		}
		if assignTo == msg.MemberId {
			tm.IsExecutor = model.Executor
		}
		if err = tx.Model(&model.TaskMember{}).Save(tm).Error; err != nil {
			zap.L().Error("Transaction Save TaskMember mysql error", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	display := t.ToTaskDisplay()
	memInfo, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: assignTo})
	if err != nil {
		return nil, err
	}
	display.Executor = model.Executors{
		Name:   memInfo.Name,
		Avatar: memInfo.Avatar,
		Code:   memInfo.Code,
	}
	//添加任务动态
	createProjectLog(ts.projectLogRepo, projectCode, t.Id, t.Name, t.AssignTo, "create", "task")
	tm := &task.TaskMessage{}
	_ = copier.Copy(tm, display)
	//发送kafka缓存删除
	kafka.SendCache([]byte("task"))
	return tm, nil
}
func createProjectLog(logRepo repo.ProjectLogRepo, projectCode int64, taskCode int64, taskName string, toMemberCode int64, logType string, actionType string) {
	remark := ""
	if logType == "create" {
		remark = "创建了任务"
	}
	pl := &model.ProjectLog{
		MemberCode:  toMemberCode,
		SourceCode:  taskCode,
		Content:     taskName,
		Remark:      remark,
		ProjectCode: projectCode,
		CreateTime:  time.Now().UnixMilli(),
		Type:        logType,
		ActionType:  actionType,
		Icon:        "plus",
		IsComment:   0,
		IsRobot:     0,
	}
	logRepo.SaveProjectLog(pl)
}
func (ts *TaskService) TaskSort(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskSortResponse, error) {
	var code int
	//拿到preTaskCode 、toStageCode
	preTaskCodeStr, _ := encrypts.Decrypt(msg.PreTaskCode, model.AESKey)
	preTaskCode, _ := strconv.ParseInt(preTaskCodeStr, 10, 64)
	toStageCodeStr, _ := encrypts.Decrypt(msg.ToStageCode, model.AESKey)
	toStageCode, _ := strconv.ParseInt(toStageCodeStr, 10, 64)
	taskCur, err := ts.taskRepo.FindTaskById(ctx, preTaskCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("TaskSort FindTaskById error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	isChanged := false
	//任务步骤变化 移动到其他步骤
	if taskCur.StageCode != int(toStageCode) {
		taskCur.StageCode = int(toStageCode)
		isChanged = true
	}
	//开启事务操作
	db := gorms.NewDBClient()
	//TODO() 完成排序功能，目前功能不完善  先查出所有需要排序的task，再将其进行重新排序
	err = db.Transaction(func(tx *gorm.DB) error {
		var nextTask = &model.Task{}
		if msg.NextTaskCode != "" {
			nextTaskCodeStr, _ := encrypts.Decrypt(msg.NextTaskCode, model.AESKey)
			NextTaskCode, _ := strconv.ParseInt(nextTaskCodeStr, 10, 64)
			//查出移动的task id
			if err = tx.Model(&model.Task{}).Where("id = ?", NextTaskCode).First(nextTask).Error; err != nil {
				zap.L().Error("nextTask TaskSort FindTaskById error", zap.Error(err))
				return err
			}
			taskCur.Sort, nextTask.Sort = nextTask.Sort, taskCur.Sort
			isChanged = true
			if err = tx.Model(&model.Task{}).Where("id = ?", NextTaskCode).Select("sort", "stage_code").Updates(&nextTask).Error; err != nil {
				zap.L().Error("nextTask TaskSort UpdateTaskSort error", zap.Error(err))
				return err
			}
		}
		if isChanged {
			if err = tx.Model(&model.Task{}).Select("sort", "stage_code").Where("id = ?", preTaskCode).Updates(&taskCur).Error; err != nil {
				zap.L().Error("nextTask TaskSort UpdateTaskSort error", zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &task.TaskSortResponse{}, nil
}
func (ts *TaskService) MyTaskList(ctx context.Context, msg *task.TaskReqMessage) (*task.MyTaskListResponse, error) {
	var taskList []*model.Task
	var err error
	var total int64
	var code int
	//根据api接口的信息查Task内容，查出来时list
	switch msg.TaskType {
	case 1: //我执行的
		taskList, total, err = ts.taskRepo.FindTaskByAssignTo(ctx, msg.MemberId, msg.Type, msg.Page, msg.PageSize)
		if err != nil {
			code = e.MysqlError
			zap.L().Error("MyTaskList FindTaskByAssignTo error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
	case 2: //我参与的
		taskList, total, err = ts.taskRepo.FindTaskByMemberId(ctx, msg.MemberId, msg.Type, msg.Page, msg.PageSize)
		if err != nil {
			code = e.MysqlError
			zap.L().Error("MyTaskList FindTaskByAssignTo error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
	case 3: //我创建的
		taskList, total, err = ts.taskRepo.FindTaskByCreateBy(ctx, msg.MemberId, msg.Type, msg.Page, msg.PageSize)
		if err != nil {
			code = e.MysqlError
			zap.L().Error("MyTaskList FindTaskByAssignTo error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
	}
	//查不到就返回nil
	if taskList == nil || len(taskList) <= 0 {
		return &task.MyTaskListResponse{List: nil, Total: 0}, nil
	}
	//把TaskList构造成接口文档 返回值的样子
	var pids []int64
	var mids []int64
	//收集project_code 和 assignTo 集合
	for _, v := range taskList {
		pids = append(pids, v.ProjectCode)
		mids = append(mids, v.AssignTo)
	}
	//根据id表查project和member信息
	projectList, err := ts.projectRepo.FindProjectByPids(ctx, pids)
	projectMap := model.ToProjectMap(projectList)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("MyTaskList FindProjectByPids error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	memberList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{
		MIds: mids,
	})
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range memberList.List {
		mMap[v.Id] = v
	}
	if err != nil {
		code = e.MysqlError
		zap.L().Error("MyTaskList FindMemInfoByIds error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//组装需要返回的结构体值
	var mtdList []*model.MyTaskDisplay
	for _, v := range taskList {
		memberMessage := mMap[v.AssignTo]
		name := memberMessage.Name
		avatar := memberMessage.Avatar
		mtd := v.ToMyTaskDisplay(projectMap[v.ProjectCode], name, avatar)
		mtdList = append(mtdList, mtd)
	}
	var myMsgs []*task.MyTaskMessage
	_ = copier.Copy(&myMsgs, mtdList)
	return &task.MyTaskListResponse{List: myMsgs, Total: total}, nil
}
func (ts *TaskService) ReadTask(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	var code int
	taskCodeStr, _ := encrypts.Decrypt(msg.TaskCode, model.AESKey)
	taskCode, _ := strconv.ParseInt(taskCodeStr, 10, 64)
	taskInfo, err := ts.taskRepo.FindTaskById(ctx, taskCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("ReadTask FindTaskById error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if taskInfo == nil {
		return &task.TaskMessage{}, nil
	}
	display := taskInfo.ToTaskDisplay()
	if display.Private == 1 {
		//代表隐私模式
		taskMember, err := ts.taskRepo.FindTaskMemberByTaskId(ctx, taskCode, msg.MemberId)
		if err != nil {
			code = e.MysqlError
			zap.L().Error("ReadTask FindTaskMemberByTaskId error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
		if taskMember != nil {
			display.CanRead = model.CanRead
		} else {
			display.CanRead = model.NoCanRead
		}
	}
	//查项目名称
	project, err := ts.projectRepo.FindProjectById(ctx, taskInfo.ProjectCode)
	display.ProjectName = project.Name
	//查项目步骤名称
	taskStages, err := ts.taskStagesRepo.FindById(ctx, int64(taskInfo.StageCode))
	display.StageName = taskStages.Name
	//查用户信息,补充service返回的信息
	memberMessage, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: taskInfo.AssignTo})
	if err != nil {
		code = e.MysqlError
		zap.L().Error("ReadTask FindMemInfoById error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	display.Executor = model.Executors{
		Name:   memberMessage.Name,
		Avatar: memberMessage.Avatar,
	}
	var taskMessage = &task.TaskMessage{}
	copier.Copy(taskMessage, display)
	return taskMessage, nil
}
func (ts *TaskService) ListTaskMember(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMemberList, error) {
	var code int
	taskCodeStr, _ := encrypts.Decrypt(msg.TaskCode, model.AESKey)
	taskCode, _ := strconv.ParseInt(taskCodeStr, 10, 64)
	//查task-member关联表
	taskMemberPage, total, err := ts.taskRepo.FindTaskMemberPage(ctx, taskCode, msg.Page, msg.PageSize)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task ListTaskMember taskRepo.FindTaskMemberPage error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	var mids []int64
	for _, v := range taskMemberPage {
		mids = append(mids, v.MemberCode)
	}
	//根据用户id数组去查用户信息
	memberInfoList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mids})
	mMap := make(map[int64]*login.MemberMessage, len(memberInfoList.List))
	for _, v := range memberInfoList.List {
		mMap[v.Id] = v
	}
	//组装数据
	var taskMemberMessageList []*task.TaskMemberMessage
	for _, v := range taskMemberPage {
		tm := &task.TaskMemberMessage{}
		tm.Code, _ = encrypts.EncryptInt64(v.MemberCode, model.AESKey)
		tm.Id = v.Id
		message := mMap[v.MemberCode]
		tm.Name = message.Name
		tm.Avatar = message.Avatar
		tm.IsExecutor = int32(v.IsExecutor)
		tm.IsOwner = int32(v.IsOwner)
		taskMemberMessageList = append(taskMemberMessageList, tm)
	}
	return &task.TaskMemberList{List: taskMemberMessageList, Total: total}, nil
}
func (ts *TaskService) TaskLog(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskLogList, error) {
	var code int
	taskCodeStr, _ := encrypts.Decrypt(msg.TaskCode, model.AESKey)
	taskCode, _ := strconv.ParseInt(taskCodeStr, 10, 64)
	all := msg.All
	var list []*model.ProjectLog
	var total int64
	var err error
	if all == 1 {
		//显示全部
		list, total, err = ts.projectLogRepo.FindLogByTaskCode(ctx, taskCode, int(msg.Comment))
		if err != nil {
			code = e.MysqlError
			zap.L().Error("project task TaskLog projectLogRepo.FindLogByTaskCode error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
	}
	if all == 0 {
		//分页查
		list, total, err = ts.projectLogRepo.FindLogByTaskCodePage(ctx, taskCode, int(msg.Comment), int(msg.Page), int(msg.PageSize))
		if err != nil {
			code = e.MysqlError
			zap.L().Error("project task TaskLog projectLogRepo.FindLogByTaskCodePage error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
	}
	//if err != nil {
	//	code = e.MysqlError
	//	zap.L().Error("project task TaskLog projectLogRepo.FindLogByTaskCodePage error", zap.Error(err))
	//	return nil, status.Error(codes.Code(code), e.GetMsg(code))
	//}
	if total == 0 {
		return &task.TaskLogList{}, nil
	}
	//组装数据
	var displayList []*model.ProjectLogDisplay
	var mIdList []int64
	for _, v := range list {
		mIdList = append(mIdList, v.MemberCode)
	}
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIdList})
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		mMap[v.Id] = v
	}
	for _, v := range list {
		display := v.ToDisplay()
		message := mMap[v.MemberCode]
		m := model.Member{}
		m.Name = message.Name
		m.Id = message.Id
		m.Avatar = message.Avatar
		m.Code = message.Code
		display.Member = m
		displayList = append(displayList, display)
	}
	var l []*task.TaskLog
	_ = copier.Copy(&l, displayList)
	return &task.TaskLogList{List: l, Total: total}, nil
}
func (ts *TaskService) TaskWorkTimeList(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskWorkTimeResponse, error) {
	taskCodeStr, _ := encrypts.Decrypt(msg.TaskCode, model.AESKey)
	taskCode, _ := strconv.ParseInt(taskCodeStr, 10, 64)
	var code int
	var list []*model.TaskWorkTime
	var err error
	list, err = ts.taskRepo.FindWorkTimeList(ctx, taskCode)
	if err != nil {
		code = e.ParseGrpcError
		zap.L().Error("project task TaskWorkTimeList taskWorkTimeRepo.FindWorkTimeList error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if len(list) == 0 {
		return &task.TaskWorkTimeResponse{}, nil
	}
	var displayList []*model.TaskWorkTimeDisplay
	var mIdList []int64
	for _, v := range list {
		mIdList = append(mIdList, v.MemberCode)
	}
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIdList})
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		mMap[v.Id] = v
	}
	//组装数据
	for _, v := range list {
		display := v.ToDisplay()
		message := mMap[v.MemberCode]
		m := model.Member{}
		m.Name = message.Name
		m.Id = message.Id
		m.Avatar = message.Avatar
		m.Code = message.Code
		display.Member = m
		displayList = append(displayList, display)
	}
	var l []*task.TaskWorkTime
	_ = copier.Copy(&l, displayList)
	return &task.TaskWorkTimeResponse{List: l, Total: int64(len(l))}, nil
}
func (ts *TaskService) SaveTaskWorkTime(ctx context.Context, msg *task.TaskReqMessage) (*task.SaveTaskWorkTimeResponse, error) {
	var code int
	taskCodeStr, _ := encrypts.Decrypt(msg.TaskCode, model.AESKey)
	taskCode, _ := strconv.ParseInt(taskCodeStr, 10, 64)
	tmt := &model.TaskWorkTime{
		BeginTime:  msg.BeginTime,
		Num:        int(msg.Num),
		Content:    msg.Content,
		TaskCode:   taskCode,
		MemberCode: msg.MemberId,
	}
	err := ts.taskRepo.SaveTaskWorkTime(ctx, tmt)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task SaveTaskWorkTime taskWorkTimeRepo.Save error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	return &task.SaveTaskWorkTimeResponse{}, nil
}
func (ts *TaskService) SaveTaskFile(ctx context.Context, msg *task.TaskFileReqMessage) (*task.TaskFileResponse, error) {
	var code int
	//获取参数
	taskCodeStr, _ := encrypts.Decrypt(msg.TaskCode, model.AESKey)
	taskCode, _ := strconv.ParseInt(taskCodeStr, 10, 64)
	OrganizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	OrganizationCode, _ := strconv.ParseInt(OrganizationCodeStr, 10, 64)
	ProjectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	ProjectCode, _ := strconv.ParseInt(ProjectCodeStr, 10, 64)
	//存file表
	f := &model.File{
		PathName:         msg.PathName,
		Title:            msg.FileName,
		Extension:        msg.Extension,
		Size:             int(msg.Size),
		ObjectType:       "",
		OrganizationCode: OrganizationCode,
		TaskCode:         taskCode,
		ProjectCode:      ProjectCode,
		CreateBy:         msg.MemberId,
		CreateTime:       time.Now().UnixMilli(),
		Downloads:        0,
		Extra:            "",
		Deleted:          model.NoDeleted,
		FileType:         msg.FileType,
		FileUrl:          msg.FileUrl,
		DeletedTime:      0,
	}
	err := ts.fileRepo.Save(ctx, f)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task SaveTaskFile fileRepo.Save error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//存入source_link表
	sl := &model.SourceLink{
		SourceType:       "file",
		SourceCode:       f.Id,
		LinkType:         "task",
		LinkCode:         taskCode,
		OrganizationCode: OrganizationCode,
		CreateBy:         msg.MemberId,
		CreateTime:       time.Now().UnixMilli(),
		Sort:             0,
	}
	err = ts.sourceLinkRepo.Save(context.Background(), sl)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task SaveTaskFile sourceLinkRepo.Save error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	return &task.TaskFileResponse{}, nil
}

func (ts *TaskService) TaskSources(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskSourceResponse, error) {
	//获取参数
	var code int
	taskCodeStr, _ := encrypts.Decrypt(msg.TaskCode, model.AESKey)
	taskCode, _ := strconv.ParseInt(taskCodeStr, 10, 64)
	//根据taskCode查source_link表
	sourceLinks, err := ts.sourceLinkRepo.FindByTaskCode(context.Background(), taskCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task SaveTaskFile sourceLinkRepo.FindByTaskCode error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if len(sourceLinks) == 0 {
		return &task.TaskSourceResponse{}, nil
	}
	//把fileId提出来
	var fIdList []int64
	for _, v := range sourceLinks {
		fIdList = append(fIdList, v.SourceCode)
	}
	//用fileId列表去查file表的文件
	files, err := ts.fileRepo.FindByIds(context.Background(), fIdList)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task SaveTaskFile fileRepo.FindByIds error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	fMap := make(map[int64]*model.File)
	for _, v := range files {
		fMap[v.Id] = v
	}
	//组装数据
	var list []*model.SourceLinkDisplay
	for _, v := range sourceLinks {
		list = append(list, v.ToDisplay(fMap[v.SourceCode]))
	}
	var slMsg []*task.TaskSourceMessage
	copier.Copy(&slMsg, list)
	return &task.TaskSourceResponse{List: slMsg}, nil
}
func (ts *TaskService) CreateComment(ctx context.Context, msg *task.TaskReqMessage) (*task.CreateCommentResponse, error) {
	var code int
	//获取参数
	taskCodeStr, _ := encrypts.Decrypt(msg.TaskCode, model.AESKey)
	taskCode, _ := strconv.ParseInt(taskCodeStr, 10, 64)
	taskById, err := ts.taskRepo.FindTaskById(context.Background(), taskCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task CreateComment fileRepo.FindTaskById error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//组装数据
	pl := &model.ProjectLog{
		MemberCode:   msg.MemberId,
		Content:      msg.CommentContent,
		Remark:       msg.CommentContent,
		Type:         "createComment",
		CreateTime:   time.Now().UnixMilli(),
		SourceCode:   taskCode,
		ActionType:   "task",
		ToMemberCode: 0,
		IsComment:    model.Comment,
		ProjectCode:  taskById.ProjectCode,
		Icon:         "plus",
		IsRobot:      0,
	}
	//保存log数据
	ts.projectLogRepo.SaveProjectLog(pl)
	return &task.CreateCommentResponse{}, nil
}

func (ts *TaskService) GetLogBySelfProject(ctx context.Context, msg *task.ProjectRpcMessage) (*task.ProjectLogResponse, error) {
	var code int
	//获取参数
	page := msg.Page
	pageSize := msg.PageSize
	//根据用户Id查项目日志表
	projectLogs, total, err := ts.projectLogRepo.FindLogByMemberCode(ctx, msg.MemberId, page, pageSize)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task GetLogBySelfProject FindLogByMemberCode error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//构建project_code、member_code、task_code列表用于查询
	pIdList := make([]int64, len(projectLogs))
	mIdList := make([]int64, len(projectLogs))
	taskIdList := make([]int64, len(projectLogs))
	for _, v := range projectLogs {
		pIdList = append(pIdList, v.ProjectCode)
		mIdList = append(mIdList, v.MemberCode)
		taskIdList = append(taskIdList, v.SourceCode)
	}
	//查询项目信息
	projectList, err := ts.projectRepo.FindProjectByPids(ctx, pIdList)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task GetLogBySelfProject FindProjectByPids error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	pMap := make(map[int64]*model.Project)
	for _, v := range projectList {
		pMap[v.Id] = v
	}
	//查用户信息
	memberInfoList, _ := rpc.LoginServiceClient.FindMemInfoByIds(context.Background(), &login.UserMessage{MIds: mIdList})
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range memberInfoList.List {
		mMap[v.Id] = v
	}
	//查task任务信息
	taskList, err := ts.taskRepo.FindTaskByIds(context.Background(), taskIdList)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("project task GetLogBySelfProject FindProjectByPids error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	tMap := make(map[int64]*model.Task)
	for _, v := range taskList {
		tMap[v.Id] = v
	}
	//组装数据
	var list []*model.IndexProjectLogDisplay
	for _, v := range projectLogs {
		display := v.ToIndexDisplay()
		display.ProjectName = pMap[v.ProjectCode].Name
		display.MemberAvatar = mMap[v.MemberCode].Avatar
		display.MemberName = mMap[v.MemberCode].Name
		display.TaskName = tMap[v.SourceCode].Name
		list = append(list, display)
	}
	var msgList []*task.ProjectLogMessage
	copier.Copy(&msgList, list)
	return &task.ProjectLogResponse{List: msgList, Total: total}, nil
}
func (ts *TaskService) NodeList(ctx context.Context, msg *task.ProjectRpcMessage) (*task.ProjectNodeResponseMessage, error) {
	var code int
	nodes, err := ts.projectRepo.FindProjectNodeAll(ctx)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("NodeList FindProjectNodeAll error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//转化数据
	treeList := model.ToNodeTreeList(nodes)
	//组装数据
	var projectNodes []*task.ProjectNodeMessage
	_ = copier.Copy(&projectNodes, treeList)
	return &task.ProjectNodeResponseMessage{Nodes: projectNodes}, nil
}
