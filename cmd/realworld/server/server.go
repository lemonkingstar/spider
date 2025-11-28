package server

import (
	"github.com/lemonkingstar/spider/cmd/realworld/conf"
	"github.com/lemonkingstar/spider/cmd/realworld/server/user"
	"github.com/lemonkingstar/spider/cmd/realworld/service"
	"github.com/lemonkingstar/spider/pkg/iserver"
	"github.com/lemonkingstar/spider/pkg/pgin"
)

var ginServer = pgin.Default()

func Get() *pgin.GinServer { return ginServer }

func Startup() error {
	ginServer.Init(
		iserver.ServiceName(conf.GetServer().Http.Name),
		iserver.Address(conf.GetServer().Http.Addr),
	)

	eng := ginServer.Engine()
	api := eng.Group("/v1")
	{
		user.NewHandler(api, service.NewUserService())
	}
	return nil
}
