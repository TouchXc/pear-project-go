package interceptor

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"ms_project/project-common/encrypts"
	"ms_project/project-grpc/task"
	"ms_project/project-task/internal/dao"
	"ms_project/project-task/internal/repo"
	"strings"
	"time"
)

type Interceptor struct {
	cache    repo.Cache
	cacheMap map[string]any
}

func NewInterceptor() *Interceptor {
	cacheMap := make(map[string]any)
	cacheMap["/TaskService/TaskList"] = &task.TaskListResponse{}
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
func (i *Interceptor) RedisCache() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		i = NewInterceptor()
		respType := i.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		//先查询是否有缓存 有的话 直接返回 无 先请求 然后存入缓存
		con, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		marshal, _ := json.Marshal(req)
		cacheKey := encrypts.Md5(string(marshal))
		respJson, _ := i.cache.Get(con, info.FullMethod+"::"+cacheKey)
		if respJson != "" {
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + " 走了缓存")
			return respType, nil
		}
		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)
		i.cache.Put(con, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		zap.L().Info(info.FullMethod + " 放入缓存")
		//hash key task field rediskey
		if strings.HasPrefix(info.FullMethod, "/task") {
			i.cache.HSet(con, "task", info.FullMethod+"::"+cacheKey, "")
		}
		return
	}
}
