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
	Name    string
	Address string
}

type DelegateOption func(*Delegate)

func NewDelegate(fns ...DelegateOption) *Delegate {
	opt := &Delegate{}
	for _, f := range fns {
		f(opt)
	}
	return opt
}

func WithName(name string) DelegateOption {
	return func(d *Delegate) {
		d.Name = name
	}
}

func WithAddress(addr string) DelegateOption {
	return func(o *Delegate) {
		o.Address = addr
	}
}
