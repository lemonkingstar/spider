package pbase

type Server interface {
	// Name server name
	Name() string
	// Start server
	Start() error
	// Stop the server
	Stop() error
	// GracefulStop the server gracefully
	GracefulStop() error
}

type Worker interface {
	// Start worker
	Start() error
	// Stop worker
	Stop() error
}

type Option struct {
	// 服务名称
	Name string
	// 绑定ip（默认 0.0.0.0）
	BindIP string
	// 绑定端口
	BindPort int
}

type OptFunc func(*Option)

func NewOption(fns ...OptFunc) *Option {
	opt := &Option{
		Name:     "default",
		BindIP:   "0.0.0.0",
		BindPort: 8080,
	}
	for _, f := range fns {
		f(opt)
	}
	return opt
}

func SetName(name string) OptFunc {
	_idx.Name = name
	return func(o *Option) {
		o.Name = name
	}
}

func SetBindIP(ip string) OptFunc {
	_idx.IP = ip
	return func(o *Option) {
		o.BindIP = ip
	}
}

func SetBindIPort(port int) OptFunc {
	_idx.Port = port
	return func(o *Option) {
		o.BindPort = port
	}
}
