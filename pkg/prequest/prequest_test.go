package prequest

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestHttpRequest(t *testing.T) {
	_, body, _ := New().Get("https://baidu.com/", nil, nil)
	t.Log(string(body))
}

func TestAddCookie(t *testing.T) {
	//header := http.Header{}
	//header.Set("Cookie", c.String())

	header := map[string]string{}
	arr := make([]string, 0)
	cs := []*http.Cookie{{Name: "auth_key", Value: "xxx"}, {Name: "token", Value: "xxx"}}
	for _, c := range cs {
		arr = append(arr, c.String())
	}
	cookie := strings.Join(arr, ";")
	header["Cookie"] = cookie
	t.Log(cookie)
}

func TestAddCookie2(t *testing.T) {
	_, body, _ := New().AddCookies([]*http.Cookie{{Name: "auth_key", Value: "xxx"}}).
		Get("https://autumnfish.cn/api/joke/list?num=1", nil, nil)
	t.Log(string(body))
}

func TestAddQuery(t *testing.T) {
	_, body, _ := New().AddQuery(map[string]string{"name": "jack"}).
		AddQuery(map[string]interface{}{"age": 18}).
		AddQuery(map[string]interface{}{"num": 10}).
		Get("https://autumnfish.cn/api/joke/list", nil, nil)
	t.Log(string(body))
}

func TestUploadFile(t *testing.T) {
	url := "http://10.10.10.10:8000/v1/upload"
	data := map[string]interface{}{
		"message": "今天是疯狂星期五",
	}
	r := New()
	// 添加文件
	// gorequest会有自动增加 file1,2序号的问题
	r.AddFile("file", "/tmp/test_file.txt")
	r.AddFile("image", "/tmp/test_image.jpg")
	_, body, _ := r.EndFile(url, data)
	t.Log(string(body))

	// 重新发送文件需要先重置
	r.Clean()
	b, _ := ioutil.ReadFile("/tmp/test_image2.jpg")
	r.AddFile("image", b, "test_image2.jpg")
	_, body, _ = r.EndFile(url, data)
	t.Log(string(body))
}

func TestUploadFile2(t *testing.T) {
	url := "http://10.10.10.10:8000/v1/upload"
	data := map[string]string{
		"message": "今天是疯狂星期五",
	}
	files := map[string]string{
		"file":  "/tmp/test_file.txt",
		"image": "/tmp/test_image.jpg",
	}
	_, body, _ := SendFile(url, files, data)
	t.Log(string(body))
}
