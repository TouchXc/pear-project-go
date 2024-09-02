package routers

import (
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"ms_project/project-api/api/user"
	"ms_project/project-common/discovery"
	"ms_project/project-common/logs"
	"ms_project/project-grpc/user/login"
	"ms_project/project-user/config"
	LoginServiceV1 "ms_project/project-user/pkg/service/login_service_v1"
	"net"
)

type InterfaceRouter interface {
	Router(r *gin.Engine)
}
type RegisterRouter struct {
}

func New() *RegisterRouter {
	return &RegisterRouter{}
}
func (*RegisterRouter) Router(ro InterfaceRouter, r *gin.Engine) {
	ro.Router(r)
}

func InitRouter(r *gin.Engine) {
	rg := New()
	rg.Router(&user.RouterUser{}, r)
}

type GrpcConfig struct {
	Addr         string
	RegisterFunc func(server *grpc.Server)
}

func RegisterGrpc(gc *config.GrpcConfig) *grpc.Server {
	c := GrpcConfig{
		Addr: gc.Addr,
		RegisterFunc: func(g *grpc.Server) {
			login.RegisterLoginServiceServer(g, LoginServiceV1.New())
		}}
	//cacheInterceptor := interceptor.NewInterceptor()
	s := grpc.NewServer(
		//cacheInterceptor.CacheInterceptor(),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(),
		)),
	)
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", gc.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return s
}
func RegisterEtcdServer() {
	in := config.InitConfig()
	etcdRegister := discovery.NewResolver(in.Ec.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	info := discovery.Server{
		Name:    in.GC.Name,
		Addr:    in.GC.Addr,
		Version: in.GC.Version,
		Weight:  in.GC.Weight,
	}
	r := discovery.NewRegister(in.Ec.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
