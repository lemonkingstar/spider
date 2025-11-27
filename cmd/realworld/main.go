package main

import (
	"github.com/go-kit/kit/transport/grpc"
	"github.com/lemonkingstar/spider"
	"github.com/lemonkingstar/spider/cmd/realworld/conf"
	"github.com/lemonkingstar/spider/pkg/plog"
)

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *spider.Application {
	app := new(spider.Application)
	app.Startup(
		conf.Startup,
		model.Startup,
		//demo.Startup,
		admin.Startup,
		router.Startup,
	).Server(
		//demo.Get(),
		admin.New(),
		router.New(),
	).ArgsFunc(parseCmd).Execute()
}

func main() {
	plog.SetLevel(plog.InfoLevel)
	plog.SetReportCaller(true)

}
