package p8s

import "testing"

func TestQuery(t *testing.T) {
	cli, err := NewDefault("https://localhost")
	if err != nil {
		t.Fatal(err)
	}
	lst, count, err := cli.QueryVector("100 - cpu_usage_idle{cpu=\"cpu-total\", cmdb_biz=\"pops\"} > 10", 10)
	t.Log(err, count)
	for _, v := range lst {
		t.Log(v.String())
	}
}
