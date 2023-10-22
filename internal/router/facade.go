package router

import "github.com/rasha108bik/tiny_url/internal/server/handlers"

// RouterFacade struct for return constructor.
type RouterFacade struct {
	HTTPRoute HTTPRoute
	GRPCRoute GRPCRoute
}

// NewRouterFacade returns a newly *RouterFacade objects that registery pattern and middleware.
func NewRouterFacade(h handlers.Handlers) *RouterFacade {
	return &RouterFacade{
		HTTPRoute: newHTTPRoute(h),
		GRPCRoute: newGRPCRoute(h),
	}
}
