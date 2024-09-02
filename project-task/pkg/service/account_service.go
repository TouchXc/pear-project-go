package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ms_project/project-common/e"
	"ms_project/project-common/encrypts"
	"ms_project/project-grpc/account"
	"ms_project/project-grpc/user/login"
	"ms_project/project-task/internal/dao"
	"ms_project/project-task/internal/gorms"
	"ms_project/project-task/internal/model"
	"ms_project/project-task/internal/repo"
	"ms_project/project-task/internal/rpc"
	"strconv"
	"time"
)

type AccountService struct {
	account.UnimplementedAccountServiceServer
	cache           repo.Cache
	accountRepo     repo.AccountRepo
	projectAuthRepo repo.ProjectAuthRepo
	departmentRepo  repo.DepartmentRepo
	menuRepo        repo.MenuRepo
	projectRepo     repo.ProjectRepo
}

func NewAccountService() *AccountService {
	return &AccountService{
		cache:           dao.RC,
		accountRepo:     dao.NewAccountDao(),
		projectAuthRepo: dao.NewProjectAuthDao(),
		departmentRepo:  dao.NewDepartmentDao(),
		menuRepo:        dao.NewMenuDao(),
		projectRepo:     dao.NewProjectDao(),
	}
}
func (as *AccountService) Account(ctx context.Context, msg *account.AccountReqMessage) (*account.AccountResponse, error) {
	//获取参数
	var code int
	var condition string
	//ProjectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	//ProjectCode, _ := strconv.ParseInt(ProjectCodeStr, 10, 64)
	OrganizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	OrganizationCode, _ := strconv.ParseInt(OrganizationCodeStr, 10, 64)
	DepartmentCodeStr, _ := encrypts.Decrypt(msg.DepartmentCode, model.AESKey)
	DepartmentCode, _ := strconv.ParseInt(DepartmentCodeStr, 10, 64)
	//MemberId := msg.MemberId
	SearchType := msg.SearchType
	Page := msg.Page
	PageSize := msg.PageSize
	//查询member_account数据
	switch SearchType {
	case 1:
		condition = "status = 1"
	case 2:
		condition = "department_code = NULL"
	case 3:
		condition = "status = 0"
	case 4:
		condition = fmt.Sprintf("status = 1 and department_code = %d", DepartmentCode)
	default:
		condition = "status = 1"
	}
	accountList, total, err := as.accountRepo.FindAccountList(ctx, condition, OrganizationCode, DepartmentCode, Page, PageSize)
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" Account FindAccountList error", zap.Error(err))
		return &account.AccountResponse{}, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//数据转换 model_MemberAccount -> MemberAccountDisplay
	var dList []*model.MemberAccountDisplay
	for _, v := range accountList {
		display := v.ToDisplay()
		memberInfo, _ := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: v.MemberCode})
		display.Avatar = memberInfo.Avatar
		if v.DepartmentCode > 0 {
			department, err := as.departmentRepo.FindDepartmentById(v.DepartmentCode)
			if err != nil {
				code = e.MysqlError
				zap.L().Error(" Account FindDepartmentById error", zap.Error(err))
				return &account.AccountResponse{}, status.Error(codes.Code(code), e.GetMsg(code))
			}
			display.Departments = department.Name
		}
		dList = append(dList, display)
	}
	//查询project_auth数据
	authList, err := as.projectAuthRepo.FindAuthList(ctx, OrganizationCode)
	if err != nil {
		code = e.MysqlError
		return &account.AccountResponse{}, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//转换数据 model_ProjectAuth -> ProjectAuthDisplay
	var pdList []*model.ProjectAuthDisplay
	for _, v := range authList {
		display := v.ToDisplay()
		pdList = append(pdList, display)
	}
	//组装数据
	var maList []*account.MemberAccount
	copier.Copy(&maList, dList)
	var prList []*account.ProjectAuth
	copier.Copy(&prList, pdList)
	return &account.AccountResponse{
		AccountList: maList,
		AuthList:    prList,
		Total:       total,
	}, nil
}
func (as *AccountService) SaveDepartment(ctx context.Context, msg *account.DepartmentReqMessage) (*account.DepartmentMessage, error) {
	//获取参数
	var code int
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	//var page = msg.Page
	//var pageSize = msg.PageSize
	var name = msg.Name
	//var departmentCode int64
	//if msg.DepartmentCode != "" {
	//	departmentCodeStr, _ := encrypts.Decrypt(msg.DepartmentCode, model.AESKey)
	//	departmentCode, _ = strconv.ParseInt(departmentCodeStr, 10, 64)
	//}
	var parentDepartmentCode int64
	if msg.ParentDepartmentCode != "" {
		parentDepartmentCodeStr, _ := encrypts.Decrypt(msg.ParentDepartmentCode, model.AESKey)
		parentDepartmentCode, _ = strconv.ParseInt(parentDepartmentCodeStr, 10, 64)
	}
	//找到对应部门
	department, err := as.departmentRepo.FindDepartment(ctx, organizationCode, parentDepartmentCode, name)
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" SaveDepartment FindDepartment error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if department == nil {
		department = &model.Department{
			Name:             name,
			CreateTime:       time.Now().UnixMilli(),
			OrganizationCode: organizationCode,
		}
		if parentDepartmentCode > 0 {
			department.Pcode = parentDepartmentCode
		}
		err = as.departmentRepo.SaveDepartment(ctx, department)
		if err != nil {
			code = e.MysqlError
			zap.L().Error(" SaveDepartment SaveDepartment error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
	}
	//转化数据
	dpd := department.ToDisplay()
	var res *account.DepartmentMessage
	_ = copier.Copy(&res, &dpd)
	return res, nil
}
func (as *AccountService) ListDepartment(ctx context.Context, msg *account.DepartmentReqMessage) (*account.ListDepartmentMessage, error) {
	//获取参数
	var code int
	var page = msg.Page
	var pageSize = msg.PageSize
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	var parentDepartmentCode int64
	if msg.ParentDepartmentCode != "" {
		parentDepartmentCodeStr, _ := encrypts.Decrypt(msg.ParentDepartmentCode, model.AESKey)
		parentDepartmentCode, _ = strconv.ParseInt(parentDepartmentCodeStr, 10, 64)
	}
	//查询数据
	departmentList, total, err := as.departmentRepo.ListDepartment(organizationCode, parentDepartmentCode, page, pageSize)
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" ListDepartment ListDepartment error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//转换并组装数据
	var dpList []*model.DepartmentDisplay
	for _, v := range departmentList {
		dpList = append(dpList, v.ToDisplay())
	}
	var list []*account.DepartmentMessage
	_ = copier.Copy(&list, dpList)
	for _, v := range list {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
	}
	return &account.ListDepartmentMessage{List: list, Total: total}, nil
}

func (as *AccountService) ReadDepartment(ctx context.Context, msg *account.DepartmentReqMessage) (*account.DepartmentMessage, error) {
	var code int
	departmentCodeStr, _ := encrypts.Decrypt(msg.DepartmentCode, model.AESKey)
	departmentCode, _ := strconv.ParseInt(departmentCodeStr, 10, 64)
	dp, err := as.departmentRepo.FindDepartmentById(departmentCode)
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" ReadDepartment FindDepartmentById error", zap.Error(err))
		return &account.DepartmentMessage{}, status.Error(codes.Code(code), e.GetMsg(code))
	}
	var res = &account.DepartmentMessage{}
	copier.Copy(res, dp.ToDisplay())
	res.Code, _ = encrypts.EncryptInt64(res.Id, model.AESKey)
	return res, nil
}
func (as *AccountService) AuthList(ctx context.Context, msg *account.AuthReqMessage) (*account.ListAuthMessage, error) {
	//获取参数
	var code int
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	//查询project_auth
	list, total, err := as.projectAuthRepo.FindAuthListPage(ctx, organizationCode, msg.Page, msg.PageSize)
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" AuthList FindAuthListPage error", zap.Error(err))
		return &account.ListAuthMessage{}, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//转换数据
	var paList []*model.ProjectAuthDisplay
	for _, v := range list {
		display := v.ToDisplay()
		paList = append(paList, display)
	}
	//组装返回数据
	var prList []*account.ProjectAuth
	copier.Copy(&prList, paList)
	return &account.ListAuthMessage{List: prList, Total: total}, nil
}
func (as *AccountService) MenuList(ctx context.Context, msg *account.MenuReqMessage) (*account.MenuResponseMessage, error) {
	//查询菜单列表
	var code int
	menus, err := as.menuRepo.FindMenus(ctx)
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" MenuList FindMenus error", zap.Error(err))
		return &account.MenuResponseMessage{}, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//数据转换
	menuChlid := model.CovertChild(menus)
	var list []*account.MenuMessageAccount
	_ = copier.Copy(&list, menuChlid)
	return &account.MenuResponseMessage{List: list}, nil
}
func (as *AccountService) AuthApply(ctx context.Context, msg *account.AuthReqMessage) (*account.ApplyResponse, error) {
	var code int
	var err error
	authId := msg.AuthId
	if msg.Action == "getnode" {
		//1.查询全部project_node节点
		projectNodeAll, err := as.projectRepo.FindProjectNodeAll(ctx)
		if err != nil {
			code = e.MysqlError
			zap.L().Error(" AuthApply FindProjectNodeAll error", zap.Error(err))
			return &account.ApplyResponse{}, status.Error(codes.Code(code), e.GetMsg(code))
		}
		//2.查询auth_id授权的project_node节点
		nodeList, err := as.projectAuthRepo.FindProjectAuthNodeByAuthId(ctx, authId)
		if err != nil {
			code = e.MysqlError
			zap.L().Error(" AuthApply FindProjectAuthNodeByAuthId error", zap.Error(err))
			return &account.ApplyResponse{}, status.Error(codes.Code(code), e.GetMsg(code))
		}
		//3.转化数据
		list := model.ToAuthNodeTreeList(projectNodeAll, nodeList)
		//4.组装接口返回数据
		var prList []*account.ProjectNodeAuthMessage
		copier.Copy(&prList, list)
		return &account.ApplyResponse{List: prList, CheckedList: nodeList}, nil
	}
	nodes := msg.Nodes
	db := gorms.NewDBClient()
	if msg.Action == "save" {
		//先删除project_auth表相关信息  再添加新的数据
		//需要开启事务
		err = db.Transaction(func(tx *gorm.DB) error {
			//先根据auth_id删除node信息
			if err = tx.Table("ms_project_auth_node").Where("auth = ?", authId).Delete(&model.ProjectAuthNode{}).Error; err != nil {
				zap.L().Error(" AuthApply Transaction error", zap.Error(err))
				return err
			}
			//再由node信息去创建数据
			var list []*model.ProjectAuthNode
			for _, v := range nodes {
				pn := &model.ProjectAuthNode{}
				pn.Auth = authId
				pn.Node = v
				list = append(list, pn)
			}
			//这里直接传切片会报错
			for _, v := range list {
				if err = tx.Create(&v).Error; err != nil {
					zap.L().Error(" AuthApply Transaction error", zap.Error(err))
					return err
				}
			}
			return nil
		})
		if err != nil {
			code = e.MysqlError
			zap.L().Error(" AuthApply FindProjectAuthNodeByAuthId error", zap.Error(err))
			return nil, status.Error(codes.Code(code), e.GetMsg(code))
		}
	}
	return &account.ApplyResponse{}, nil
}
func (as *AccountService) AuthNodesByMemberId(ctx context.Context, msg *account.AuthReqMessage) (*account.AuthNodesResponse, error) {
	var code int
	//获取参数
	memberId := msg.MemberId
	//1.根据memberId去member_account表查角色（authorize）字段
	authId, err := as.accountRepo.FindAuthIdByMemberId(ctx, memberId)
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" AuthNodesByMemberId FindAuthIdByMemberId error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//2.再由角色字段查auth_node表中node信息返回
	nodeList, err := as.projectAuthRepo.FindProjectAuthNodeByAuthId(ctx, authId)
	if err != nil {
		code = e.MysqlError
		zap.L().Error(" AuthNodesByMemberId FindProjectAuthNodeByAuthId error", zap.Error(err))
		return &account.AuthNodesResponse{}, status.Error(codes.Code(code), e.GetMsg(code))
	}
	return &account.AuthNodesResponse{List: nodeList}, nil
}
