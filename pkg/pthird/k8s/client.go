package k8s

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client is the interface for kubernetes client
type Client interface {
	ListNamespace() ([]v1.Namespace, error)
}

type Option struct {
	// Host must be a host string, a host:port pair, or a URL to the base of the apiserver.
	// If a URL is given then the (optional) Path of that URL represents a prefix that must
	// be appended to all request URIs used to access the apiserver. This allows a frontend
	// proxy to easily relocate all of the apiserver endpoints.
	Host string

	// Server requires Bearer authentication. This client will not attempt to use
	// refresh tokens for an OAuth2 flow.
	// TODO: demonstrate an OAuth2 compatible client.
	BearerToken string

	// Config conf file content
	Config string
}

type client struct {
	ctx       context.Context
	clientSet *kubernetes.Clientset
}

func NewClient(opt Option) (Client, error) {
	var err error
	var cs *kubernetes.Clientset
	if opt.BearerToken != "" {
		config := &rest.Config{
			Host:            opt.Host, // https://10.212.32.118:6443
			TLSClientConfig: rest.TLSClientConfig{Insecure: true},
			BearerToken:     opt.BearerToken,
		}
		cs, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
	} else {
		// get config form conf file content
		config, err := clientcmd.RESTConfigFromKubeConfig([]byte(opt.Config))
		if err != nil {
			return nil, err
		}
		cs, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
	}

	return &client{
		ctx:       context.Background(),
		clientSet: cs,
	}, nil
}

func (c *client) ListNamespace() ([]v1.Namespace, error) {
	res, err := c.clientSet.CoreV1().Namespaces().List(c.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}
