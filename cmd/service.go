package cmd

import "github.com/adialaleal/odins/internal/service"

var serviceFactory = func() *service.Manager {
	return service.New(service.DefaultRuntime())
}
