package n

import "net/http"

type Route interface {
	Handler() Handler
	Methods(methods ...string) Route
	Use(mw ...Middleware) Route
}

type NRoute struct {
	methods     []string
	middlewares []Middleware
	handler     Handler
}

func NewRoute(handler Handler) *NRoute {
	r := &NRoute{}

	r.middlewares = append(r.middlewares, httpMethodMiddleware(r))

	r.handler = HandlerFunc(func(ctx Context) error {
		var target = handler

		for _, mw := range r.middlewares {
			target = mw.Handle(ctx, target)
		}

		return target.Handle(ctx)
	})

	return r
}

func (N *NRoute) Handler() Handler {
	return N.handler
}

func (N *NRoute) Methods(methods ...string) Route {
	N.methods = methods
	return N
}

func (N *NRoute) Use(mw ...Middleware) Route {
	N.middlewares = append(N.middlewares, mw...)
	return N
}

func httpMethodMiddleware(r *NRoute) MiddlewareFunc {
	return func(ctx Context, handler Handler) Handler {
		if len(r.methods) == 0 {
			return handler
		}

		for _, m := range r.methods {
			if m == ctx.Request().Method {
				return handler
			}
		}

		return HandlerFunc(func(ctx Context) error {
			ctx.Status(http.StatusMethodNotAllowed)
			return nil
		})
	}
}
