package prequest

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/lemonkingstar/spider/pkg/pconv"
)

// WrapUrl 拼接url query
func WrapUrl(addr string, params map[string]string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func WrapUrlEx(addr string, params map[string]interface{}) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for key, obj := range params {
		value, _ := pconv.Type2Str(obj)
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func SendFile(url string, files map[string]string, form map[string]string, args ...map[string]string) (*http.Response, []byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	f := func(field, path string) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		part, err := writer.CreateFormFile(field, filepath.Base(path))
		if err != nil {
			return err
		}
		io.Copy(part, file)
		return nil
	}
	// 添加文件
	for k, v := range files {
		if err := f(k, v); err != nil {
			return nil, nil, err
		}
	}
	// 添加表单
	for k, v := range form {
		writer.WriteField(k, v)
	}
	writer.Close()
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}
	// 添加头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	headers := map[string]string{}
	if len(args) == 1 {
		headers = args[0]
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, b, nil
}
