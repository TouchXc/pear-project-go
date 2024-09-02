package task

import (
	"github.com/gin-gonic/gin"
	"log"
	"ms_project/project-api/api/rpc"
	"ms_project/project-api/middleware"
	"ms_project/project-api/routers"
)

type RouterTask struct {
}

func init() {
	log.Println("init task router")
	routerTask := &RouterTask{}
	routers.Register(routerTask)
}
func (*RouterTask) Router(r *gin.Engine) {
	//初始化grpc客户端连接
	rpc.InitRpcTaskClient()
	h := &HandlerTask{}
	group := r.Group("/project")
	group.Use(middleware.TokenVerify())
	group.Use(Auth())
	{
		//TODO()路由分组细分出一个task路由
		group.POST("/index", h.Index)
		group.POST("/project", h.MyProjectList)
		group.POST("/project_template", h.ProjectTemplate)
		group.POST("/project_collect/collect", h.CollectProject)
		group.POST("/task_stages", h.TaskStages)
		group.POST("/project_member/index", h.MemberProjectList)
		group.POST("/task_stages/tasks", h.TaskList)
		group.POST("/task_member", h.ListTaskMember)
		group.POST("/file/uploadFiles", h.UploadFiles)
		group.POST("/account", h.Account)
		group.POST("/department", h.ListDepartment)
		group.POST("/department/save", h.SaveDepartment)
		group.POST("/department/read", h.ReadDepartment)
		group.POST("/auth", h.AuthList)
		group.POST("/menu/menu", h.MenuList)
		group.POST("/node", h.NodeList)
		group.POST("/auth/apply", h.AuthApply)
		project := group.Group("/project")
		{
			project.POST("/selfList", h.MyProjectList)
			project.POST("/save", h.ProjectSave)
			project.POST("/read", h.ReadProject)
			project.POST("/recycle", h.DeleteProject)
			project.POST("/recovery", h.RecoveryProject)
			project.POST("/edit", h.EditProject)
			project.POST("/getLogBySelfProject", h.GetLogBySelfProject)
		}
		task := group.Group("/task")
		{
			task.POST("/save", h.SaveTask)
			task.POST("/sort", h.TaskSort)
			task.POST("/selfList", h.MyTaskList)
			task.POST("/read", h.ReadTask)
			task.POST("/taskLog", h.TaskLog)
			task.POST("/_taskWorkTimeList", h.TaskWorkTimeList)
			task.POST("/saveTaskWorkTime", h.SaveTaskWorkTime)
			task.POST("/taskSources", h.TaskSource)
			task.POST("/createComment", h.CreateComment)
		}
	}
}
