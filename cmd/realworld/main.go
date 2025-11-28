package main

import (
	"github.com/lemonkingstar/spider"
	"github.com/lemonkingstar/spider/cmd/realworld/conf"
	"github.com/lemonkingstar/spider/cmd/realworld/data"
	"github.com/lemonkingstar/spider/cmd/realworld/server"
	"github.com/lemonkingstar/spider/pkg/plog"
)

func main() {
	plog.SetLevel(plog.InfoLevel)
	plog.SetReportCaller(true)
	app := new(spider.Application)
	app.Startup(
		conf.Startup,
		data.Startup,
		server.Startup,
	).Server(
		server.Get(),
	).Execute()
}
