package pgin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lemonkingstar/spider/pkg/pbase"
	"github.com/lemonkingstar/spider/pkg/pgin/middleware"
	"github.com/lemonkingstar/spider/pkg/plog"
)

var (
	logger = plog.GetLogger()
)

func New() *GinServer {
	gin.SetMode(gin.ReleaseMode)
	return &GinServer{engine: gin.New()}
}

func Default() *GinServer {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(middleware.RecoverMiddleware())
	g.Use(middleware.CORSMiddleware(nil))
	g.Use(middleware.LayerMiddleware(true))
	return &GinServer{engine: g}
}

type GinServer struct {
	opt    *pbase.Option
	engine *gin.Engine
	srv    *http.Server
}

func (s *GinServer) Init(fns ...pbase.OptFunc) {
	s.opt = pbase.NewOption(fns...)
}

func (s *GinServer) Name() string {
	return s.opt.Name
}

func (s *GinServer) Start() error {

	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.opt.BindIP, s.opt.BindPort),
		Handler: s.engine,
	}
	logger.Infof("Server[%s] running at: %s", s.Name(), s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil {
		logger.Errorf("listen error: %v", err)
		return err
	}
	return nil
}

func (s *GinServer) Stop() error {
	//if s.srv != nil {
	//	s.srv.Close()
	//}
	return s.GracefulStop()
}

func (s *GinServer) GracefulStop() error {
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		return nil
	}

	logger.Info("Server exiting")
	return nil
}

/******************* gin function ***************/

func (s *GinServer) Engine() *gin.Engine {
	return s.engine
}

// SetMode 设置日志级别
func (s *GinServer) SetMode(debug bool) {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
