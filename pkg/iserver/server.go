package iserver

type Server interface {
	// Name is the server name.
	Name() string
	// Start starts server.
	Start() error
	// Stop stops server.
	Stop() error
	// GracefulStop stops the server gracefully.
	GracefulStop() error
}

// Delegate contains server option.
// Can be expanded to include more functions.
type Delegate struct {
	ServerName string
	Address    string
}

type DelegateOption func(*Delegate)

func NewDelegate(fns ...DelegateOption) *Delegate {
	opt := &Delegate{
		Address: "0.0.0.0:8080",
	}
	for _, f := range fns {
		f(opt)
	}
	return opt
}

func ServerName(name string) DelegateOption {
	AppName = name
	return func(d *Delegate) {
		d.ServerName = name
	}
}

func Address(addr string) DelegateOption {
	return func(o *Delegate) {
		o.Address = addr
	}
}
