package controller

import "password-lock/service"

type Controller struct {
	service service.Service
}

func NewController(service service.Service) Controller {
	return Controller{
		service: service,
	}
}
