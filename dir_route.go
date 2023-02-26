package n

import (
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
	newRoute func(path string, handler Handler) Route
}

var (
	_pathFixer = strings.NewReplacer("\\", "/")
)

func NewDirRoute(basePath string, newRoute func(path string, handler Handler) Route) *NDirRoute {
	return &NDirRoute{
		basePath: _pathFixer.Replace(basePath),
		newRoute: newRoute,
	}
}

func (N NDirRoute) Dir(path string) DirRoute {
	return NewDirRoute(N.getPath(path), N.newRoute)
}

func (N NDirRoute) Handle(path string, handler Handler) Route {
	return N.newRoute(N.getPath(path), handler)
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
	return path.Join(N.basePath, _pathFixer.Replace(routePath))
}
