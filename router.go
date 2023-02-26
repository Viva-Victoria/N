package n

import (
	"context"
	"errors"
	"gitea.voopsen/n/log"
	"gitea.voopsen/n/sync"
	"gitea.voopsen/n/tree"
	"go.uber.org/atomic"
	"net/http"
	"regexp"
)

type routeMeta struct {
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
	badRequestCode    int
	internalErrorCode int
	routeNotFoundCode int
	routes            tree.Tree[routeMeta]
	logger            log.Logger
	listening         atomic.Bool
	wg                sync.WaitGroup
}

func NewRouter(basePath string, logger log.Logger) *NRouter {
	nr := &NRouter{
		routes:            tree.NewTree[routeMeta](),
		listening:         *atomic.NewBool(true),
		wg:                sync.NewWaitGroup(1),
		badRequestCode:    http.StatusBadRequest,
		internalErrorCode: http.StatusInternalServerError,
		routeNotFoundCode: http.StatusNotFound,
	}

	nr.NDirRoute = *NewDirRoute(basePath, logger, func(path string, handler Handler) (Route, error) {
		global, matchers, err := ParseRoute(path)
		if err != nil {
			return nil, err
		}

		nr.routes.Add(matchers, routeMeta{
			regexp:  global,
			handler: handler,
		})

		return NewRoute(handler), nil
	})

	return nr
}

func (N *NRouter) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {
	if !N.listening.Load() {
		return
	}

	_ = N.wg.Add(1)
	go func() {
		defer func() {
			N.wg.Done(1)

			if panicErr := recover(); panicErr != nil {
				N.logger.P(panicErr)
			}
		}()

		route, ok := N.routes.Find(splitPath(rq.URL.Path))
		if !ok || route.regexp == nil {
			rs.WriteHeader(N.routeNotFoundCode)
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
			newRS.WriteHeader(N.badRequestCode)
			return
		}

		newRS.WriteHeader(N.internalErrorCode)
	}()
}

func (N *NRouter) Close(ctx context.Context) error {
	N.listening.Store(false)
	N.wg.Done(1)
	return N.wg.WaitContext(ctx)
}
