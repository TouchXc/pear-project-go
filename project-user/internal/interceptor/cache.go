package interceptor

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"ms_project/project-common/encrypts"
	"ms_project/project-grpc/user/login"
	"ms_project/project-user/internal/dao"
	"ms_project/project-user/internal/repo"
	"time"
)

type Interceptor struct {
	cache    repo.Cache
	cacheMap map[string]any
}

func NewInterceptor() *Interceptor {
	cacheMap := make(map[string]any)
	cacheMap["/LoginService/MyOrgList"] = &login.OrgListResponse{}
	cacheMap["/LoginService/FindMemInfoById"] = &login.MemberMessage{}
	return &Interceptor{cache: dao.RC, cacheMap: cacheMap}
}

func (i *Interceptor) CacheInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		respType := i.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		//先查询是否有缓存，没有则请求 再存入缓存
		con, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		marshal, _ := json.Marshal(req)
		cacheKey := encrypts.Md5(string(marshal))
		respJson, _ := i.cache.Get(con, info.FullMethod+"::"+cacheKey)
		if respJson != "" {
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + "放入缓存")
			return respType, nil
		}
		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)
		i.cache.Put(con, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		return
	})
}
