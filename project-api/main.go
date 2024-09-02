package main

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"ms_project/project-api/config"
	"ms_project/project-api/tracing"
	srv "ms_project/project-common"
	routersTask "ms_project/project-task/routers"
	routersUser "ms_project/project-user/routers"
	"net/http"
)

func main() {
	//用viper读取配置文件
	IniConfig := config.InitConfig()
	//初始化日志配置
	IniConfig.InitZapLog()
	//路由初始化
	r := gin.Default()
	//加载jaeger链路追踪
	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	r.Use(otelgin.Middleware("project-api"))
	//设置静态文件地址
	r.StaticFS("/upload", http.Dir("upload"))
	//初始化service模块路由
	routersUser.InitRouter(r)
	routersTask.InitRouter(r)
	//开启pprof  默认访问路径：/debug/pprof
	pprof.Register(r)
	srv.Run(r, IniConfig.SC.Name, IniConfig.SC.Addr, nil)
}
