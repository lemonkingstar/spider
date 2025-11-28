package service

import "github.com/lemonkingstar/spider/cmd/realworld/data"

type UserService interface{}

type user struct {
	userDao *data.UserStorage
}

func NewUserService() UserService {
	return &user{
		userDao: &data.UserStorage{},
	}
}
