package n

import (
	"encoding/json"
	"net/http"
)

type Context interface {
	Request() *http.Request
	Response() ResponseWriter

	Vars() Values
	Header() Values

	Status(code int)
	ReadJSON(a any) error
	WriteJSON(a any) error
}

type NContext struct {
	vars Values
	rq   *http.Request
	rs   ResponseWriter
}

func NewContext(vars Values, rq *http.Request, rs ResponseWriter) NContext {
	return NContext{
		vars: vars,
		rq:   rq,
		rs:   rs,
	}
}

func (N NContext) Request() *http.Request {
	return N.rq
}

func (N NContext) Response() ResponseWriter {
	return N.rs
}

func (N NContext) Vars() Values {
	return N.vars
}

func (N NContext) Header() Values {
	return HeaderValues{
		header: N.rq.Header,
	}
}

func (N NContext) Status(code int) {
	N.rs.WriteHeader(code)
}

func (N NContext) ReadJSON(a any) error {
	defer N.rq.Body.Close()
	return json.NewDecoder(N.rq.Body).Decode(a)
}

func (N NContext) WriteJSON(a any) error {
	defer func() {
		N.rs.Flush()
		_ = N.rs.Close()
	}()

	N.rs.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(N.rs).Encode(a)
}
