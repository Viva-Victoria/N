package n

import (
	"net/http"
	"net/textproto"
)

type Values interface {
	Add(key string, value any)
	Set(key string, value any)
	Get(key string, ref any) error
	GetString(key string) string
	Values(key string) []string
	Delete(key string)
}

type HeaderValues struct {
	header http.Header
}

func (h HeaderValues) Add(key string, value any) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	h.header[key] = append(h.header[key], anyToStrings(value)...)
}

func (h HeaderValues) Set(key string, value any) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	h.header[key] = anyToStrings(value)
}

func (h HeaderValues) Get(key string, ref any) error {
	return stringToAny(h.header[key], ref)
}

func (h HeaderValues) GetString(key string) string {
	return h.header.Get(key)
}

func (h HeaderValues) Values(key string) []string {
	return h.header.Values(key)
}

func (h HeaderValues) Delete(key string) {
	h.header.Del(key)
}

type MapValues map[string][]string

func (m MapValues) Add(key string, value any) {
	m[key] = append(m[key], anyToStrings(value)...)
}

func (m MapValues) Set(key string, value any) {
	m[key] = anyToStrings(value)
	l := len(m[key])
	m[key] = m[key][:l:l]
}

func (m MapValues) Get(key string, ref any) error {
	return stringToAny(m[key], ref)
}

func (m MapValues) GetString(key string) string {
	return m[key][0]
}

func (m MapValues) Values(key string) []string {
	return m[key]
}

func (m MapValues) Delete(key string) {
	delete(m, key)
}
