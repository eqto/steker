package plugin

import (
	"github.com/eqto/steker/buff"
)

//Request ...
type Request interface {
	Put(data ...interface{})
	Send() (<-chan Response, error)

	bytes() []byte
}

type request struct {
	Request
	plugin *Plugin
	name   string

	data []interface{}
}

func (r *request) Put(data ...interface{}) {
	r.data = append(r.data, data...)
}

func (r *request) Send() (<-chan Response, error) {
	return r.plugin.sendRequest(r)
}

// 1 byte function name length
// function name to bytes
// bytes ...
func (r *request) bytes() []byte {
	b := buff.Writer{}

	b.PutShortString(r.name)
	b.PutByte(byte(len(r.data)))
	for _, d := range r.data {
		b.Put(toBytes(d))
	}
	return b.Bytes()
}

func toBytes(data interface{}) []byte {
	b := buff.Writer{}

	switch data := data.(type) {
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
