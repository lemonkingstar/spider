package spider

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/lemonkingstar/spider/pkg/pconf"
	"github.com/lemonkingstar/spider/pkg/perror"
	"github.com/lemonkingstar/spider/pkg/plog"
	"github.com/lemonkingstar/spider/pkg/psafe"
	"github.com/lemonkingstar/spider/pkg/server"
	"github.com/lemonkingstar/spider/pkg/worker"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

var (
	logger = plog.WithField("[PACKET]", "spider")
)

type Application struct {
	servers      []server.Server
	workers      []worker.Worker
	startups     []func() error
	beforeStart  func()
	afterStart   func()
	beforeStop   func()
	afterStop    func()
	forceRestart bool
	argsFunc     func()
}

var (
	RootCmd = &cobra.Command{
		Use:     filepath.Base(os.Args[0]),
		Version: "",
	}

	RootParam = struct {
		ConfigFile  string
		VersionFlag bool
		DaemonFlag  bool
	}{}
)

func (app *Application) Startup(fns ...func() error) *Application {
	app.startups = append(app.startups, fns...)
	return app
}

func (app *Application) Server(s ...server.Server) *Application {
	app.servers = append(app.servers, s...)
	return app
}

func (app *Application) Worker(w ...worker.Worker) *Application {
	app.workers = append(app.workers, w...)
	return app
}

func (app *Application) BeforeStart(fn func()) *Application {
	app.beforeStart = fn
	return app
}

func (app *Application) AfterStart(fn func()) *Application {
	app.afterStart = fn
	return app
}

func (app *Application) BeforeStop(fn func()) *Application {
	app.beforeStop = fn
	return app
}

func (app *Application) AfterStop(fn func()) *Application {
	app.afterStop = fn
	return app
}

func (app *Application) ArgsFunc(fn func()) *Application {
	app.argsFunc = fn
	return app
}

func (app *Application) Execute() {
	if app.beforeStart != nil {
		app.beforeStart()
	}
	RootCmd.PersistentFlags().StringVarP(&RootParam.ConfigFile, "config-file", "c", "", "config file")
	RootCmd.PersistentFlags().BoolVarP(&RootParam.VersionFlag, "version", "v", false, "show version")
	RootCmd.PersistentFlags().BoolVarP(&RootParam.DaemonFlag, "daemon", "d", false, "start as daemon")
	if app.argsFunc != nil {
		app.argsFunc()
	}
	RootCmd.SetHelpFunc(func(*cobra.Command, []string) {
		// print usage info
		RootCmd.Usage()
		os.Exit(0)
	})

	RootCmd.Run = func(cmd *cobra.Command, args []string) {
		// the work function
		if RootParam.VersionFlag {
			// show version
			fmt.Println(pconf.VERSION)
			os.Exit(0)
		} else if RootParam.DaemonFlag {
			// start as daemon
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			ctx := &daemon.Context{
				WorkDir: dir,
			}
			child, err := ctx.Reborn()
			if err != nil {
				logger.Errorf("reborn error: %v", err)
				os.Exit(-1)
			}
			if child != nil {
				return
			}
			defer ctx.Release()
		}
		app.Run()
	}

	err := RootCmd.Execute()
	if err != nil {
		logger.Errorf("command execute error: %v", err)
		os.Exit(-1)
	}
}

func (app *Application) Run() {
	if RootParam.ConfigFile != "" {
		pconf.SetConfigFile(RootParam.ConfigFile)
	}
	if err := pconf.ReadInConfig(); err != nil {
		logger.Errorf("load config error: %v", err)
		return
	}
	if err := perror.SerialUntilError(app.startups...)(); err != nil {
		logger.Errorf("startup error: %v", err)
		return
	}
	// server&worker
	startErr := sync.Once{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, s := range app.servers {
		go func(server server.Server) {
			if err := server.Start(); err != nil {
				startErr.Do(func() {
					logger.Errorf("server start error: %v", err)
					cancel()
				})
			}
		}(s)
	}
	for _, w := range app.workers {
		go func(worker worker.Worker) {
			if err := worker.Start(); err != nil {
				startErr.Do(func() {
					logger.Errorf("worker start error: %v", err)
					cancel()
				})
			}
		}(w)
	}

	if app.afterStart != nil {
		app.afterStart()
	}
	app.endingProc(ctx)
}

func (app *Application) endingProc(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-c:
		logger.Infof("stop signal caught, stopping... pid=%d", os.Getpid())
	case <-ctx.Done():
	}

	if app.beforeStop != nil {
		app.beforeStop()
	}
	// wait for stop
	pg := psafe.NewGroup()
	for _, s := range app.servers {
		pg.Run(s.Stop)
	}
	for _, w := range app.workers {
		pg.Run(w.Stop)
	}
	err := pg.WaitError()
	if err != nil {
		logger.Errorf("server stop error: %v", err)
	}
	if app.afterStop != nil {
		app.afterStop()
	}
	logger.Info("stopped")
}
