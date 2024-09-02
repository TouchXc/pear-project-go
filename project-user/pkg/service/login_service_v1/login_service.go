package login_service_v1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	common "ms_project/project-common"
	"ms_project/project-common/e"
	"ms_project/project-common/encrypts"
	"ms_project/project-common/jwt"
	"ms_project/project-common/tms"
	"ms_project/project-grpc/user/login"
	"ms_project/project-user/internal/dao"
	"ms_project/project-user/internal/database/gorms"
	"ms_project/project-user/internal/model"
	"ms_project/project-user/internal/repo"
	"strconv"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache            repo.Cache
	memberRepo       repo.MemberRepo
	organizationRepo repo.OrganizationRepo
}

func New() *LoginService {
	return &LoginService{
		cache:            dao.RC,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganizationDao(),
	}
}
func (ls *LoginService) GetCaptcha(c context.Context, msg *login.CaptchaMessage) (*login.CaptchaResponse, error) {
	//var code int
	//rsp := &common.Response{}
	//1.获取参数
	mobile := msg.Mobile
	//2.校验参数
	if !common.VerifyMobile(mobile) {
		//code = e.InValidMobile
		//c.JSON(http.StatusOK, rsp.Failed(e.InValidMobile, e.GetMsg(code)))
		return nil, status.Error(e.InValidMobile, e.GetMsg(e.InValidMobile))
	}
	//3.生成验证码
	//todo() 实现真正的短信验证码生成 随机4位/随机6位
	captchaCode := "123456"
	//4.调用短信平台（第三方 放入go协程中执行 接口可快速响应）
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("短信平台调用成功 Info")
		zap.L().Debug("短信平台调用成功 Debug")
		zap.L().Error("短信平台调用成功 Error")
		//log.Println("短信平台调用成功")
		//5.存储验证码 redis中 过期时间15分钟
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := ls.cache.Put(c, model.RegisterRedisKey+mobile, captchaCode, 15*time.Minute)
		if err != nil {
			log.Printf("验证码存入redist出错，cause by :%v/n", err)
		}
	}()
	//c.JSON(http.StatusOK, rsp.Success(captchaCode))
	return &login.CaptchaResponse{Code: captchaCode}, nil
}
func (ls *LoginService) Register(c context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	//todo() 重构校验判断 ——> switch case fall
	var code int
	redisCode, err := ls.cache.Get(context.Background(), model.RegisterRedisKey+msg.Mobile)
	if errors.Is(err, redis.Nil) {
		code = e.CaptchaExpired
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if err != nil {
		code = e.RedisError
		zap.L().Error("Register redis get error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if redisCode != msg.Captcha {
		code = e.InValidCaptcha
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	// 校验业务逻辑（邮箱是否注册、账号是否注册、手机号是否注册）
	//校验邮箱
	exist, err := ls.memberRepo.GetMemberByEmail(context.Background(), msg.Email)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("Register db mysql error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if exist {
		code = e.EmailHasExist
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//校验账号
	exist, err = ls.memberRepo.GetMemberByName(context.Background(), msg.Name)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("Register db mysql error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if exist {
		code = e.UserHasExist
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//校验手机号
	exist, err = ls.memberRepo.GetMemberByMobile(context.Background(), msg.Mobile)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("Register db mysql error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if exist {
		code = e.MobileHasExist
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//执行业务，数据存入member表 生成数据存入组织表
	pwd := encrypts.Md5(msg.Password)
	member := &model.Member{
		Account:       msg.Name,
		Password:      pwd,
		Name:          msg.Name,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        1,
	}
	db := gorms.NewDBClient()
	//创建用户及项目事务
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&model.Member{}).Create(member).Error; err != nil {
			zap.L().Error("Register db mysql error", zap.Error(err))
			return err
		}
		organization := &model.Organization{
			Name:       member.Name + "个人组织",
			MemberId:   member.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   1,
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		if err = tx.Model(&model.Organization{}).Create(organization).Error; err != nil {
			zap.L().Error("Register db mysql error", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		code = e.MysqlError
		zap.L().Error("Register db mysql error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	return &login.RegisterResponse{}, nil
}
func (ls *LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	var code int
	//1.验证账号密码 去数据库查询密码是否正确
	pwd := encrypts.Md5(msg.Password)
	member, err := ls.memberRepo.FindMemberByAccount(ctx, msg.Account, pwd)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("Login db FindMemberByAccount error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if member == nil {
		code = e.AccountOrPasswordError
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	encryptId, _ := encrypts.EncryptInt64(member.Id, model.AESKey)
	lst := tms.FormatByMill(member.LastLoginTime)
	memberMessage := &login.MemberMessage{
		Id:            member.Id,
		Name:          member.Name,
		Mobile:        member.Mobile,
		Realname:      member.Realname,
		Account:       member.Account,
		Status:        int32(member.Status),
		LastLoginTime: lst,
		Address:       member.Address,
		Province:      int32(member.Province),
		City:          int32(member.City),
		Area:          int32(member.Area),
		Email:         member.Email,
		Code:          encryptId,
	}
	//2.根据用户id查组织
	orgs, err := ls.organizationRepo.FindOrganizationByMemberId(ctx, member.Id)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("Login db FindOrganizationByMemberId error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	var orgsMessage []*login.OrganizationMessage
	copier.Copy(&orgsMessage, orgs)
	for _, v := range orgsMessage {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
		v.OwnerCode = memberMessage.Code
		o := model.ToMap(orgs)[v.Id]
		v.CreateTime = tms.FormatByMill(o.CreateTime)
	}
	if len(orgs) > 0 {
		memberMessage.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	//3.用jwt生成token
	memberIdStr := strconv.FormatInt(member.Id, 10)
	exp := model.AccessExp * 24 * time.Hour
	refreshExp := model.RefreshExp * 24 * time.Hour
	token := jwt.CreateToken(memberIdStr, exp, model.AccessSecret, model.RefreshSecret, refreshExp)
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenType:      "bearer",
		AccessTokenExp: token.AccessExp,
	}
	go func() {
		marshal, _ := json.Marshal(member)
		ls.cache.Put(ctx, model.MemberRedisKey+"::"+memberIdStr, string(marshal), exp)
		orgsJson, _ := json.Marshal(orgs)
		ls.cache.Put(ctx, model.MemberOrganizationRedisKey+"::"+memberIdStr, string(orgsJson), exp)
	}()
	return &login.LoginResponse{
		Member:           memberMessage,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}
func (ls *LoginService) TokenVerify(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	var code int
	token := msg.Token
	if token == "" {
		code = e.TokenExpired
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	parseToken, err := jwt.ParseToken(token, model.AccessSecret)
	if err != nil {
		code = e.TokenExpired
		zap.L().Error("TokenVerify ParseToken  error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	//从缓存中查询 如果没有 直接返回认证失败
	//TODO() 数据库查询 优化点 登陆后缓存用户信息
	memberJson, err := ls.cache.Get(ctx, model.MemberRedisKey+"::"+parseToken)
	if err != nil {
		code = e.RedisError
		zap.L().Error("TokenVerify redis cache get member error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if memberJson == "" {
		code = e.TokenExpired
		zap.L().Error("TokenVerify redis cache get member expire", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	memberById := &model.Member{}
	json.Unmarshal([]byte(memberJson), memberById)

	memberMessage := &login.MemberMessage{}
	_ = copier.Copy(&memberMessage, memberById)
	memberMessage.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKey)
	//缓存中读组织信息
	orgsJson, err := ls.cache.Get(ctx, model.MemberOrganizationRedisKey+"::"+parseToken)
	if err != nil {
		code = e.RedisError
		zap.L().Error("TokenVerify redis cache get organization error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if orgsJson == "" {
		code = e.TokenExpired
		zap.L().Error("TokenVerify redis cache get organization expire", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	var orgs []*model.Organization
	json.Unmarshal([]byte(orgsJson), &orgs)
	if len(orgs) > 0 {
		memberMessage.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
		memberMessage.CreateTime = tms.FormatByMill(memberById.CreateTime)
	}
	return &login.LoginResponse{Member: memberMessage}, nil
}
func (ls *LoginService) MyOrgList(ctx context.Context, msg *login.UserMessage) (*login.OrgListResponse, error) {
	var code int
	memId := msg.MemId
	orgs, err := ls.organizationRepo.FindOrganizationByMemberId(ctx, memId)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("MyOrgList FindOrganizationByMemberId", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	var orgsMessage []*login.OrganizationMessage
	_ = copier.Copy(&orgsMessage, orgs)
	for _, o := range orgsMessage {
		o.Code, _ = encrypts.EncryptInt64(o.Id, model.AESKey)
	}
	return &login.OrgListResponse{OrganizationList: orgsMessage}, nil
}
func (ls *LoginService) FindMemInfoById(ctx context.Context, msg *login.UserMessage) (*login.MemberMessage, error) {
	var code int
	member, err := ls.memberRepo.FindMemberById(context.Background(), msg.MemId)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("FindMemInfoById FindMemberById  error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if member == nil {
		code = e.MysqlError
		zap.L().Error("FindMemInfoById member is nil")
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	memberMessage := &login.MemberMessage{}
	_ = copier.Copy(&memberMessage, member)
	memberMessage.Code, _ = encrypts.EncryptInt64(member.Id, model.AESKey)
	memberMessage.CreateTime = tms.FormatByMill(member.CreateTime)
	orgs, err := ls.organizationRepo.FindOrganizationByMemberId(ctx, member.Id)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("TokenVerify FindOrganizationByMemberId error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if len(orgs) > 0 {
		memberMessage.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	return memberMessage, nil
}
func (ls *LoginService) FindMemInfoByIds(ctx context.Context, msg *login.UserMessage) (*login.MemberMessageList, error) {
	var code int
	memberList, err := ls.memberRepo.FindMemberByIds(ctx, msg.MIds)
	if err != nil {
		code = e.MysqlError
		zap.L().Error("FindMemInfoByIds FindMemberByIds  error", zap.Error(err))
		return nil, status.Error(codes.Code(code), e.GetMsg(code))
	}
	if memberList == nil || len(memberList) <= 0 {
		return &login.MemberMessageList{List: nil}, nil
	}
	mMap := make(map[int64]*model.Member)
	for _, v := range memberList {
		mMap[v.Id] = v
	}
	memberMessages := make([]*login.MemberMessage, 0)
	_ = copier.Copy(&memberMessages, memberList)
	for _, v := range memberMessages {
		m := mMap[v.Id]
		v.CreateTime = tms.FormatByMill(m.CreateTime)
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
	}
	return &login.MemberMessageList{List: memberMessages}, nil
}
