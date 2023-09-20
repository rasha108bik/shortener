package router

import "github.com/rasha108bik/tiny_url/internal/server/handlers"

type RouterFacade struct {
	HTTPRoute HTTPRoute
	GRPCRoute GRPCRoute
}

func NewRouterFacade(h handlers.Handlers) *RouterFacade {
	return &RouterFacade{
		HTTPRoute: newHTTPRoute(h),
		GRPCRoute: newGRPCRoute(h),
	}
}
