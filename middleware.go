package n

type Middleware interface {
	Handle(ctx Context, handler Handler) Handler
}

type MiddlewareFunc func(ctx Context, handler Handler) Handler

func (m MiddlewareFunc) Handle(ctx Context, handler Handler) Handler {
	if m == nil {
		return handler
	}
	
	return m(ctx, handler)
}
