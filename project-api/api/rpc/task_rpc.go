package rpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"ms_project/project-api/config"
	"ms_project/project-common/discovery"
	"ms_project/project-common/logs"
	"ms_project/project-grpc/account"
	"ms_project/project-grpc/task"
)

var TaskServiceClient task.TaskServiceClient
var AccountServiceClient account.AccountServiceClient

func InitRpcTaskClient() {
	etcdRegister := discovery.NewResolver(config.InConf.Ec.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.NewClient("127.0.0.1:8882",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	//conn, err := grpc.NewClient("etcd:///task", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	TaskServiceClient = task.NewTaskServiceClient(conn)
	AccountServiceClient = account.NewAccountServiceClient(conn)
}
