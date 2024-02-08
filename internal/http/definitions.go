package http

import "github.com/melisource/fury_go-core/pkg/web"

type Controller interface {
	AddRoutes(Router)
}

type Router interface {
	Get(pattern string, handler web.Handler, mw ...web.Middleware)
	Patch(pattern string, handler web.Handler, mw ...web.Middleware)
	Post(pattern string, handler web.Handler, mw ...web.Middleware)
}
