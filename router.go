package n

import (
	"context"
	"errors"
	"gitea.voopsen/OSS/n/log"
	"gitea.voopsen/OSS/n/sync"
	"gitea.voopsen/OSS/n/tree"
	"github.com/google/uuid"
	"go.uber.org/atomic"
	"net/http"
	"regexp"
)

type routeMeta struct {
	id      string
	regexp  *regexp.Regexp
	handler Handler
}

type Router interface {
	DirRoute
	http.Handler
	Close(ctx context.Context) error
}

type NRouter struct {
	NDirRoute
	config    RouterConfig
	routes    tree.Tree[map[string]routeMeta]
	logger    log.Logger
	listening atomic.Bool
	wg        sync.WaitGroup
}

func NewRouter(basePath string, logger log.Logger, options ...RouterOption) *NRouter {
	config := RouterConfig{
		BadRequestCode:    http.StatusBadRequest,
		InternalErrorCode: http.StatusInternalServerError,
		RouteNotFoundCode: http.StatusNotFound,
	}
	for _, o := range options {
		o.Apply(&config)
	}

	nr := &NRouter{
		routes:    tree.NewTree[map[string]routeMeta](),
		listening: *atomic.NewBool(true),
		wg:        sync.NewWaitGroup(1),
		config:    config,
	}

	nr.NDirRoute = *NewDirRoute(basePath, logger, func(path string, handler Handler) (Route, error) {
		global, matchers, err := ParseRoute(path)
		if err != nil {
			return nil, err
		}

		v, ok := nr.routes.Get(matchers)
		if !ok || len(v) == 0 {
			v = make(map[string]routeMeta)
			nr.routes.Add(matchers, v)
		}
		meta := routeMeta{
			id:      uuid.NewString(),
			regexp:  global,
			handler: handler,
		}

		if _, ok = v["*"]; !ok {
			v["*"] = meta
		}

		return NewRoute(handler, func(methods []string) {
			for m := range v {
				if v[m].id == meta.id {
					delete(v, m)
				}
			}

			for _, m := range methods {
				v[m] = meta
			}
		}), nil
	})

	return nr
}

func (N *NRouter) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {
	if !N.listening.Load() {
		return
	}

	_ = N.wg.Add(1)
	func() {
		defer func() {
			N.wg.Done(1)

			if panicErr := recover(); panicErr != nil {
				N.logger.P(panicErr)
			}
		}()

		routes, ok := N.routes.Find(splitPath(rq.URL.Path))
		if !ok || len(routes) == 0 {
			rs.WriteHeader(N.config.RouteNotFoundCode)
			return
		}

		route, ok := routes[rq.Method]
		if !ok || route.regexp == nil {
			route, ok = routes["*"]
		}
		if !ok || route.regexp == nil {
			rs.WriteHeader(N.config.RouteNotFoundCode)
			return
		}

		var (
			vars       = make(MapValues)
			matches    = route.regexp.FindAllStringSubmatch(rq.URL.Path, -1)
			groupNames = route.regexp.SubexpNames()
		)
		for _, match := range matches {
			for groupId, group := range match {
				name := groupNames[groupId]
				if len(name) == 0 {
					continue
				}

				vars.Set(name, group)
			}
		}

		newRS := NewResponseWriter(rs)
		err := route.handler.Handle(NewContext(vars, rq, newRS))
		if err == nil {
			return
		}

		N.logger.E(err)
		if newRS.IsCommitted() {
			return
		}

		var badRequestError BadRequestError
		if errors.As(err, &badRequestError) {
			newRS.WriteHeader(N.config.BadRequestCode)
			return
		}

		newRS.WriteHeader(N.config.InternalErrorCode)
	}()
}

func (N *NRouter) Close(ctx context.Context) error {
	N.listening.Store(false)
	N.wg.Done(1)
	return N.wg.WaitContext(ctx)
}
