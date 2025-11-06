package psafe

// GoLimit 限制协程最大并发数
type GoLimit struct {
	Num int
	C   chan struct{}
}

// NewGl
// usage: psafe.NewGl(10).Run(func() {})
func NewGl(num int) *GoLimit {
	return &GoLimit{
		Num: num,
		C:   make(chan struct{}, num),
	}
}

func (g *GoLimit) Run(f func()) {
	g.C <- struct{}{}
	go func() {
		Call(f)
		<-g.C
	}()
}
