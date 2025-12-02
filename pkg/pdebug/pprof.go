package pdebug

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/lemonkingstar/spider/pkg/iserver"
	"github.com/lemonkingstar/spider/pkg/plog"
)

func NewServer() *PProfServer { return &PProfServer{} }

type PProfServer struct {
	opt *iserver.Delegate
	srv *http.Server
}

func (s *PProfServer) Handler() http.Handler {
	mux := http.NewServeMux()
	// std pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// custom debug
	mux.HandleFunc("/debug/system", system)
	mux.HandleFunc("/debug/version", version)
	return mux
}

func (s *PProfServer) Init(fns ...iserver.DelegateOption) {
	s.opt = iserver.NewDelegate(fns...)
	if s.opt.Address == "" {
		s.opt.Address = ":8081"
	}
}

func (s *PProfServer) Name() string {
	return "pprof"
}

func (s *PProfServer) Start() error {
	s.srv = &http.Server{
		Addr:    s.opt.Address,
		Handler: s.Handler(),
	}
	plog.Infof("pprof server running at: %s", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil {
		plog.Errorf("pprof server listen error: %v", err)
		return err
	}
	return nil
}

func (s *PProfServer) Stop() error {
	//if s.srv != nil {
	//	s.srv.Close()
	//}
	return s.GracefulStop()
}

func (s *PProfServer) GracefulStop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		plog.Errorf("pprof server forced to shutdown: %v", err)
		return nil
	}

	plog.Info("pprof server exiting")
	return nil
}
