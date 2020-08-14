package plugin

import (
	"fmt"
	"math"

	"github.com/eqto/steker/buff"
)

//Request ...
type Request interface {
	Put(data interface{})
	PutValue(key string, val interface{})
	Send() (<-chan Response, error)
}

type request struct {
	Request
	plugin *Plugin
	id     uint16
	name   string

	data    []Value
	dataMap map[string]Value
}

func (r *request) Put(data interface{}) {
	r.data = append(r.data, buff.NewValue(data))
}

func (r *request) PutValue(key string, val interface{}) {
	if r.dataMap == nil {
		r.dataMap = make(map[string]Value)
	}
	r.dataMap[key] = buff.NewValue(val)
}

func (r *request) Send() (<-chan Response, error) {
	return r.plugin.sendRequest(r)
}

// 1 byte function name length
// function name to bytes
// bytes ...
func (r *request) bytes() ([]byte, error) {
	if r.dataMap == nil {
		r.dataMap = make(map[string]Value)
	}
	length := len(r.data) + len(r.dataMap)

	if length > math.MaxUint16 {
		return nil, fmt.Errorf(`too many parameters, maximum is %d parameters`, math.MaxUint16)
	}

	b := buff.Writer{}
	b.PutShortString(r.name)
	b.PutUint16(length)

	for _, d := range r.data {
		b.Put(toBytes(d))
	}
	for key, val := range r.dataMap {
		b.PutByte(buff.DataStringMap)
		if len(key) > math.MaxUint8 {
			key = key[:math.MaxUint8]
		}
		b.PutShortString(key)
		b.PutBytes(val.Bytes())
	}
	return b.Bytes(), nil
}

func toBytes(data Value) []byte {
	b := buff.Writer{}

	switch data := data.Raw().(type) {
	case []byte: //TODO handle long bytes
		b.PutByte(buff.DataBytes)
		b.PutBytes(data)
	case int:
		b.PutByte(buff.DataInt)
		b.PutInt(data)
	case float64:
		b.PutByte(buff.DataFloat)
		b.PutFloat(data)
	case string: //TODO handle long string
		b.PutByte(buff.DataString)
		b.PutString(data)
	default:
		logE(`type not supported, %v`, data)
	}
	return b.Bytes()
}
