package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	srv "ms_project/project-common"
	"ms_project/project-user/config"
	"ms_project/project-user/internal/dao"
	"ms_project/project-user/internal/database/gorms"
	"ms_project/project-user/routers"
	"ms_project/project-user/tracing"
	"strings"
)

func main() {
	//用viper读取配置文件
	IniConfig := config.InitConfig()
	// 加载redis
	dao.InitRedis()
	//加载mysql
	mysqlConfig := IniConfig.MysqlConfig
	path := strings.Join([]string{mysqlConfig.Username, ":", mysqlConfig.Password, "@tcp(", mysqlConfig.Host, ":", mysqlConfig.Port, ")/", mysqlConfig.DbName, "?charset=utf8mb4&parseTime=true"}, "")
	gorms.Database(path)
	//初始化日志配置
	IniConfig.InitZapLog()
	//注册grpc
	gc := routers.RegisterGrpc(IniConfig.GC)
	//grpc服务注册至etcd
	IniConfig.ReadEtcdConfig()
	stop := func() {
		gc.Stop()
	}
	//路由初始化
	r := gin.Default()
	//加载jaeger链路追踪
	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	routers.InitRouter(r)
	srv.Run(r, IniConfig.SC.Name, IniConfig.SC.Addr, stop)
}
