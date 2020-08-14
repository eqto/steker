package service

import (
	"fmt"
	"strings"

	"github.com/eqto/steker/buff"
)

//Request ...
type Request interface {
	Put(data interface{})
	Get(idx int) Value
	GetValue(key string) Value
}

type request struct {
	Request
	name string
	id   uint16

	data    []Value
	dataMap map[string]Value
}

func (r *request) Put(data interface{}) {
	r.data = append(r.data, buff.NewValue(data))
}

func (r *request) putValue(key string, val interface{}) {
	if r.dataMap == nil {
		r.dataMap = make(map[string]Value)
	}
	r.dataMap[key] = buff.NewValue(val)
}

func (r *request) Get(idx int) Value {
	if idx >= len(r.data) {
		return buff.NewValue(nil)
	}
	return r.data[idx]
}

func (r *request) GetValue(key string) Value {
	if r.dataMap == nil {
		r.dataMap = make(map[string]Value)
	}
	if val, ok := r.dataMap[key]; ok {
		return val
	}
	return buff.NewValue(nil)
}

func (r *request) String() string {
	buf := strings.Builder{}

	buf.WriteString(fmt.Sprintf("Request ID: %d\nName: %s\n", r.id, r.name))
	if r.dataMap != nil {
		for key, val := range r.dataMap {
			buf.WriteString(fmt.Sprintf("- %s: %s\n", key, val.String()))
		}
	}
	for key, val := range r.data {
		buf.WriteString(fmt.Sprintf("- [%d] %s\n", key, val.String()))
	}

	return buf.String()
}
