package pssh

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/lemonkingstar/spider/pkg/psafe"
)

type ClientOptions struct {
	// Addr Host:Port
	Host string
	Port int

	// Auth information
	User          string
	Password      string
	KeyFile       string
	KeyBytes      []byte
	KeyPassphrase []byte // rsa key密码
}

func (p *ClientOptions) init() {
	// do options check
	if p.Port <= 0 {
		p.Port = 22
	}
	if p.User == "" {
		p.User = "root"
	}
}

type Client struct {
	opt *ClientOptions

	// Client session
	client  *ssh.Client
	session *ssh.Session
	sync.Mutex
}

func NewClient(opt *ClientOptions) *Client {
	opt.init()
	return &Client{
		opt: opt,
	}
}

func (c *Client) connect() error {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		config       ssh.Config
		err          error
	)
	opt := c.opt

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(opt.Password))
	if len(opt.KeyBytes) == 0 && len(opt.KeyFile) > 0 {
		opt.KeyBytes, err = ioutil.ReadFile(opt.KeyFile)
		if err != nil {
			return err
		}
	}
	if len(opt.KeyBytes) > 0 {
		var signer ssh.Signer
		if len(opt.KeyPassphrase) == 0 {
			signer, err = ssh.ParsePrivateKey(opt.KeyBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(opt.KeyBytes, opt.KeyPassphrase)
		}
		if err != nil {
			return err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	clientConfig = &ssh.ClientConfig{
		User:    opt.User,
		Auth:    auth,
		Timeout: 10 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if c.client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		clientConfig); err != nil {
		return err
	}
	if c.session, err = c.client.NewSession(); err != nil {
		return err
	}
	return nil
}

func (c *Client) close() {
	c.Lock()
	defer c.Unlock()
	if c.session != nil {
		c.session.Close()
		c.session = nil
	}
	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
}

func (c *Client) Close() {
	c.close()
}

func (c *Client) Run(cmd string) (string, error) {
	if err := c.connect(); err != nil {
		return "", err
	}
	defer c.close()
	b, err := c.session.CombinedOutput(cmd)
	return string(b), err
}

func (c *Client) Run2w(cmd string, w io.Writer) error {
	if err := c.connect(); err != nil {
		return err
	}
	defer c.close()

	c.session.Stdout = w
	c.session.Stderr = w
	return c.session.Run(cmd)
}

func (c *Client) Run2f(cmd string, f func(string)) error {
	if err := c.connect(); err != nil {
		return err
	}
	defer c.close()

	ch := make(chan error, 1)
	r, w := io.Pipe()
	c.session.Stdout = w
	c.session.Stderr = w
	psafe.Go(func() {
		defer func() {
			w.Close()
			r.Close()
		}()
		ch <- c.session.Run(cmd)
	})
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		f(scanner.Text())
	}
	return <-ch
}
