package prequest

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	timeout = 60 * time.Second
)

type HttpRequest struct{
	agent *gorequest.SuperAgent

	bounceToRawString    	bool
	disableSecurityCheck	bool
}

func New() *HttpRequest {
	client := &HttpRequest{
		agent: gorequest.New(),
		bounceToRawString: false,
		disableSecurityCheck: false,
	}
	client.agent.Client.Timeout = timeout
	return client
}

func (p *HttpRequest) Clean() *HttpRequest {
	p.agent.ClearSuperAgent()
	return p
}

func (p *HttpRequest) SetBounceToRawString(v bool) *HttpRequest {
	p.bounceToRawString = v
	return p
}

// DisableSecurityCheck 禁用https安全校验
func (p *HttpRequest) DisableSecurityCheck(v bool) *HttpRequest {
	p.disableSecurityCheck = v
	return p
}

func (p *HttpRequest) SetTimeout(d time.Duration) *HttpRequest {
	p.agent.Client.Timeout = d
	return p
}

// AddQuery 添加查询参数
// AddQuery(map[string]string{"name": "jack"})
// AddQuery(map[string]interface{"age": 18})
func (p *HttpRequest) AddQuery(params interface{}) *HttpRequest {
	if value, ok := params.(map[string]string); ok {
		for k, v := range value { p.agent.Param(k, v) }
	} else if value, ok := params.(map[string]interface{}); ok {
		for k, v := range value { p.agent.Param(k, fmt.Sprintf("%v", v)) }
	}
	return p
}

func (p *HttpRequest) AddCookies(cookies []*http.Cookie) *HttpRequest {
	p.agent.AddCookies(cookies)
	return p
}

func (p *HttpRequest) send(method, url string, headers map[string]string, data interface{}) {
	p.agent.BounceToRawString = p.bounceToRawString
	if p.disableSecurityCheck {
		p.agent.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	p.agent.Method = method
	p.agent.Url = url
	for k, v := range headers {
		p.agent = p.agent.Set(k, v)
	}
	p.agent.Send(data)
}

func (p *HttpRequest) Request(method, url string, headers map[string]string, data interface{}) (*http.Response, []byte, error) {
	p.send(method, url, headers, data)
	resp, body, errs := p.agent.EndBytes()
	p.Clean()
	if len(errs) > 0 {
		return nil, nil, errs[0]
	}

	return resp, body, nil
}

func (p *HttpRequest) RequestStruct(method, url string, headers map[string]string, data interface{}, v interface{}) (*http.Response, []byte, error) {
	p.send(method, url, headers, data)
	resp, body, errs := p.agent.EndStruct(v)
	p.Clean()
	if len(errs) > 0 {
		return nil, nil, errs[0]
	}

	return resp, body, nil
}

func (p *HttpRequest) Get(url string, headers map[string]string, data interface{}) (*http.Response, []byte, error) {
	return p.Request("GET", url, headers, data)
}

func (p *HttpRequest) GetStruct(url string, headers map[string]string, data interface{}, v interface{}) (*http.Response, []byte, error) {
	return p.RequestStruct("GET", url, headers, data, v)
}

func (p *HttpRequest) Post(url string, headers map[string]string, data interface{}) (*http.Response, []byte, error) {
	return p.Request("POST", url, headers, data)
}

func (p *HttpRequest) Put(url string, headers map[string]string, data interface{}) (*http.Response, []byte, error) {
	return p.Request("PUT", url, headers, data)
}

func (p *HttpRequest) Patch(url string, headers map[string]string, data interface{}) (*http.Response, []byte, error) {
	return p.Request("PATCH", url, headers, data)
}

func (p *HttpRequest) Delete(url string, headers map[string]string, data interface{}) (*http.Response, []byte, error) {
	return p.Request("DELETE", url, headers, data)
}

// AddFile 添加文件
func (p *HttpRequest) AddFile(field string, file interface{}, args ...string) *HttpRequest {
	switch file.(type) {
	case string:
		// file: 文件路径 filename: filepath.Base(file)
		p.agent.SendFile(file, "", field)
	case []byte:
		filename := ""
		if len(args) == 1 { filename = args[0] }
		// file: 文件内容 filename:
		p.agent.SendFile(file, filename, field)
	}
	return p
}

func (p *HttpRequest) EndFile(url string, form map[string]interface{}, args ...map[string]string) (*http.Response, []byte, error) {
	p.agent.Method = "POST"
	p.agent.Url = url
	headers := map[string]string{}
	if len(args) == 1 { headers = args[0] }
	for k, v := range headers {
		p.agent = p.agent.Set(k, v)
	}
	if len(form) > 0 {
		p.agent.SendMap(form)
	}

	resp, body, errs := p.agent.Type("multipart").EndBytes()
	if len(errs) > 0 {
		return nil, nil, errs[0]
	}

	return resp, body, nil
}
