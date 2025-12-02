package spider

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"

	"github.com/lemonkingstar/spider/pkg/iserver"
	"github.com/lemonkingstar/spider/pkg/iworker"
	"github.com/lemonkingstar/spider/pkg/pconf"
	"github.com/lemonkingstar/spider/pkg/perror"
	"github.com/lemonkingstar/spider/pkg/plog"
	"github.com/lemonkingstar/spider/pkg/psafe"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

var (
	logger = plog.WithField("[PACKET]", "spider")
)

type Application struct {
	servers       []iserver.Server
	workers       []iworker.Worker
	startups      []func() error
	beforeStart   func()
	afterStart    func()
	beforeStartup func()
	afterStartup  func()
	beforeStop    func()
	afterStop     func()
}

var (
	command = &cobra.Command{
		Use:     filepath.Base(os.Args[0]),
		Version: "",
	}

	CommandParam = struct {
		ConfigFile  string
		VersionFlag bool
		DaemonFlag  bool
	}{}
)

func (app *Application) Startup(fns ...func() error) *Application {
	app.startups = append(app.startups, fns...)
	return app
}

func (app *Application) Server(s ...iserver.Server) *Application {
	app.servers = append(app.servers, s...)
	return app
}

func (app *Application) Worker(w ...iworker.Worker) *Application {
	app.workers = append(app.workers, w...)
	return app
}

func (app *Application) BeforeStartup(fn func()) *Application {
	app.beforeStartup = fn
	return app
}

func (app *Application) AfterStartup(fn func()) *Application {
	app.afterStartup = fn
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

func (app *Application) Execute() {
	if app.beforeStart != nil {
		app.beforeStart()
	}
	command.PersistentFlags().StringVarP(&CommandParam.ConfigFile, "config-file", "c", "", "config file")
	command.PersistentFlags().BoolVarP(&CommandParam.VersionFlag, "version", "v", false, "show version")
	command.PersistentFlags().BoolVarP(&CommandParam.DaemonFlag, "daemon", "d", false, "start as daemon")
	command.SetHelpFunc(func(*cobra.Command, []string) {
		// print usage info
		command.Usage()
		os.Exit(0)
	})

	command.Run = func(cmd *cobra.Command, args []string) {
		// the work function
		if CommandParam.VersionFlag {
			// show version
			fmt.Println(iserver.GetAppName())
			fmt.Println("version:", iserver.GetVersion())
			fmt.Println("branch:", iserver.GetBuildBranch())
			fmt.Println("commit:", iserver.GetBuildCommit())
			fmt.Println("golang version:", runtime.Version())
			os.Exit(0)
		} else if CommandParam.DaemonFlag {
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
	}
	err := command.Execute()
	if err != nil {
		logger.Errorf("command execute error: %v", err)
		os.Exit(-1)
	}
	// parse config
	if CommandParam.ConfigFile != "" {
		pconf.SetConfigFile(CommandParam.ConfigFile)
	}
	if err := pconf.ReadInConfig(); err != nil {
		logger.Errorf("load config error: %v", err)
		return
	}
	app.Run()
}

func (app *Application) Run() {
	if app.beforeStartup != nil {
		app.beforeStartup()
	}
	if err := perror.SerialUntilError(app.startups...)(); err != nil {
		logger.Errorf("startup error: %v", err)
		return
	}
	if app.afterStartup != nil {
		app.afterStartup()
	}
	// server&worker
	startErr := sync.Once{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, s := range app.servers {
		go func(server iserver.Server) {
			if err := server.Start(); err != nil {
				startErr.Do(func() {
					logger.Errorf("server start error: %v", err)
					cancel()
				})
			}
		}(s)
	}
	for _, w := range app.workers {
		go func(worker iworker.Worker) {
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
