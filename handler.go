package n

type Handler interface {
	Handle(ctx Context) error
}

type HandlerFunc func(ctx Context) error

func (h HandlerFunc) Handle(ctx Context) error {
	if h == nil {
		return nil
	}

	return h(ctx)
}
