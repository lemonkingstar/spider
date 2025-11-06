package predis

var (
	defaultCli Client
)

// Default return the default client
func Default() Client { return defaultCli }

// NewDefault init the default client
func NewDefault(cfg Config) (Client, error) {
	cli, err := New(cfg)
	if err != nil { return nil, err }
	defaultCli = cli
	return cli, nil
}
