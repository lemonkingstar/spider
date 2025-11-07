package pconsul

import (
	"testing"
)

func TestConsul(t *testing.T) {
	// 创建实例信息
	inst, _ := NewServiceInstance("test-001", "test",
		"localhost", 8080,
		map[string]string{"app": "myapp", "version": "1.0.0"},
		[]string{"test"})
	// 创建client
	api, _ := NewServiceRegistry("http://localhost:8500")
	// 注册实例
	api.Register(inst, false)
}
