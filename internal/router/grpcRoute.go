package router

import "github.com/rasha108bik/tiny_url/internal/server/handlers"

type GRPCRoute struct {
}

func newGRPCRoute(s handlers.Handlers) GRPCRoute {
	return GRPCRoute{}
}
