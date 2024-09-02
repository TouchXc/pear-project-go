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
	"ms_project/project-grpc/user/login"
)

var LoginServiceClient login.LoginServiceClient

func InitRpcUserClient() {
	etcdRegister := discovery.NewResolver(config.InConf.Ec.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.NewClient("127.0.0.1:8881",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	//conn, err := grpc.NewClient("etcd:///user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	LoginServiceClient = login.NewLoginServiceClient(conn)
}
