package n

import (
	"encoding/json"
	"net/http"
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
	h.header.Add(key, anyToString(value))
}

func (h HeaderValues) Set(key string, value any) {
	h.header.Set(key, anyToString(value))
}

func (h HeaderValues) Get(key string, ref any) error {
	return stringToAny(h.header.Get(key), ref)
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
	m[key] = append(m[key], anyToString(value))
}

func (m MapValues) Set(key string, value any) {
	arr, ok := m[key]
	if ok {
		arr = arr[:1:1]
	} else {
		arr = make([]string, 1, 1)
	}

	arr[0] = anyToString(value)
	m[key] = arr
}

func (m MapValues) Get(key string, ref any) error {
	return stringToAny(m[key][0], ref)
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

func anyToString(a any) string {
	if s, ok := a.(string); ok {
		return s
	}

	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func stringToAny(s string, ref any) error {
	if sp, ok := ref.(*string); ok {
		*sp = s
		return nil
	}

	return json.Unmarshal([]byte(s), ref)
}
