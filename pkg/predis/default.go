package predis

var (
	defaultCli Client
)

// Default return the default client
func Default() Client { return defaultCli }

// CreateDefault init the default client
func CreateDefault(cfg Config) (Client, error) {
	cli, err := New(cfg)
	if err != nil {
		return nil, err
	}
	defaultCli = cli
	return cli, nil
}

// AcquireDefault usual way to acquire client
func AcquireDefault(addr, pw string, db int) (Client, error) {
	return CreateDefault(Config{Address: addr, Password: pw, Database: db})
}
