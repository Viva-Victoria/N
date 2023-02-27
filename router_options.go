package n

type RouterConfig struct {
	BadRequestCode    int
	InternalErrorCode int
	RouteNotFoundCode int
}

type RouterOption interface {
	Apply(*RouterConfig)
}

type RouterOptionFunc func(*RouterConfig)

func (r RouterOptionFunc) Apply(config *RouterConfig) {
	r(config)
}

func WithBadRequest(code int) RouterOption {
	return RouterOptionFunc(func(config *RouterConfig) {
		config.BadRequestCode = code
	})
}

func WithInternalError(code int) RouterOption {
	return RouterOptionFunc(func(config *RouterConfig) {
		config.InternalErrorCode = code
	})
}

func WithRouteNotFound(code int) RouterOption {
	return RouterOptionFunc(func(config *RouterConfig) {
		config.RouteNotFoundCode = code
	})
}
