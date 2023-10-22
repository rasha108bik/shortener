package router

import (
	"github.com/rasha108bik/tiny_url/internal/server/handlers"
)

// GRPCRoute struct for return constructor.
type GRPCRoute struct {
}

// newGRPCRoute returns a newly GRPCRoute objects that registery pattern and middleware.
func newGRPCRoute(s handlers.Handlers) GRPCRoute {
	return GRPCRoute{}
}
