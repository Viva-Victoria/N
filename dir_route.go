package n

import (
	"gitea.voopsen/n/log"
	"net/http"
	"path"
	"strings"
)

type DirRoute interface {
	Dir(path string) DirRoute
	Handle(path string, handler Handler) Route
	Get(path string, handler Handler) Route
	Post(path string, handler Handler) Route
	Patch(path string, handler Handler) Route
	Put(path string, handler Handler) Route
	Delete(path string, handler Handler) Route
}

type NDirRoute struct {
	basePath string
	logger   log.Logger
	newRoute func(path string, handler Handler) (Route, error)
}

var (
	_pathFixer = strings.NewReplacer("\\", "/")
)

func NewDirRoute(basePath string, logger log.Logger, newRoute func(path string, handler Handler) (Route, error)) *NDirRoute {
	return &NDirRoute{
		basePath: _pathFixer.Replace(basePath),
		logger:   logger,
		newRoute: newRoute,
	}
}

func (N NDirRoute) Dir(path string) DirRoute {
	return NewDirRoute(N.getPath(path), N.logger, N.newRoute)
}

func (N NDirRoute) Handle(path string, handler Handler) Route {
	r, err := N.newRoute(N.getPath(path), handler)
	if err != nil {
		N.logger.E(err)
	}
	return r
}

func (N NDirRoute) Get(path string, handler Handler) Route {
	return N.Handle(path, handler).Methods(http.MethodGet)
}

func (N NDirRoute) Post(path string, handler Handler) Route {
	return N.Handle(path, handler).Methods(http.MethodPost)
}

func (N NDirRoute) Patch(path string, handler Handler) Route {
	return N.Handle(path, handler).Methods(http.MethodPatch)
}

func (N NDirRoute) Put(path string, handler Handler) Route {
	return N.Handle(path, handler).Methods(http.MethodPut)
}

func (N NDirRoute) Delete(path string, handler Handler) Route {
	return N.Handle(path, handler).Methods(http.MethodDelete)
}

func (N NDirRoute) getPath(routePath string) string {
	return path.Join(N.basePath, fixPath(routePath))
}

func fixPath(path string) string {
	var (
		varOpened bool
	)

	runes := []rune(path)
	for i, r := range runes {
		if r == '{' {
			varOpened = true
			continue
		}
		if r == '}' {
			varOpened = false
			continue
		}

		if varOpened {
			continue
		}
		if r == '\\' {
			runes[i] = '/'
		}
	}

	return string(runes)
}
