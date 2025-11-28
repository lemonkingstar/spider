package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lemonkingstar/spider/cmd/realworld/server/base"
	"github.com/lemonkingstar/spider/cmd/realworld/service"
	"github.com/lemonkingstar/spider/pkg/pgin"
)

func NewHandler(r *gin.RouterGroup, s service.UserService) base.Handler {
	h := &handler{s: s}
	r.GET("/user/profile", h.profile)
	return h
}

type handler struct {
	s service.UserService
}

func (h *handler) profile(c *gin.Context) {
	pgin.Success(c, gin.H{"name": "Jack", "age": 18})
}
