package n

import (
	"errors"
	"gitea.voopsen/n/log"
	"gitea.voopsen/n/tree"
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
}

type NRouter struct {
	NDirRoute
	badRequestCode    int
	internalErrorCode int
	routes            tree.Tree[routeMeta]
	logger            log.Logger
}

func NewRouter() *NRouter {
	nr := &NRouter{
		routes: tree.NewTree[routeMeta](),
	}

	nr.handle = func(path string, handler Handler) Route {
		global, matchers, err := ParseRoute(path)
		if err != nil {

			return nil
		}

		nr.routes.Add(matchers, routeMeta{
			regexp:  global,
			handler: handler,
		})

		return NewRoute(handler)
	}

	return nr
}

func (N *NRouter) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {
	go func() {
		if panicErr := recover(); panicErr != nil {
			N.logger.P(panicErr)
		}

		route, ok := N.routes.Find(splitPath(rq.URL.Path))
		if !ok {
			rs.WriteHeader(http.StatusNotFound)
			return
		}

		var (
			vars       = make(MapValues)
			matches    = route.regexp.FindAllStringSubmatch(rq.URL.Path, -1)
			groupNames = route.regexp.SubexpNames()
		)
		for groupId, group := range matches {
			name := groupNames[groupId]
			vars.Set(name, group[0])
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
