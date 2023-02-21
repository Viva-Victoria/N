package n

import "net/http"

type Route interface {
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

	r.middlewares = append(r.middlewares, MiddlewareFunc(func(ctx Context, handler Handler) Handler {
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
	}))

	r.handler = HandlerFunc(func(ctx Context) error {
		for _, mw := range r.middlewares {
			handler = mw.Handle(ctx, handler)
		}

		return handler.Handle(ctx)
	})

	return r
}

func (N *NRoute) Methods(methods ...string) Route {
	N.methods = methods
	return N
}

func (N *NRoute) Use(mw ...Middleware) Route {
	N.middlewares = append(N.middlewares, mw...)
	return N
}
