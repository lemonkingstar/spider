package k8s

import "testing"

var (
	opt = Option{
		Host:        "https://localhost:6443",
		BearerToken: "dangerous",
	}
)

func TestNamespace(t *testing.T) {
	cli, err := NewClient(opt)
	if err != nil {
		t.Fatal(err)
	}
	lst, err := cli.ListNamespace()
	t.Log(err)
	for _, v := range lst {
		t.Log(v.String())
	}
}
