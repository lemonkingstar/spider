package server

import "github.com/lemonkingstar/spider/pkg/pgin"

var ginServer = pgin.Default()

func Get() *pgin.GinServer { return ginServer }

func Startup() error {
	ginServer.Init(
		common.SetName("phoenix-demo"),
		common.SetBindIP("0.0.0.0"),
		common.SetBindIPort(setting.C.Port),
	)

	core := service.NewService(setting.C)
	eng := ginServer.Engine()
	api := eng.Group("/api/v1")
	{
		// routers可以不用保存
		// 成员函数的指针有被引用
		// routers = append(routers, user.NewUserRouter(api, core))
		user.NewUserRouter(api, core)
	}
	return nil
}
