package buff

import (
	"fmt"
	"math"
	"strings"
)

//Response ...
type Response struct {
	id      uint16
	err     error
	data    []Value
	dataMap map[string]Value
}

//PutString ...
func (r *Response) PutString(s string) {
	r.data = append(r.data, NewValue(s))
}

//PutInt ...
func (r *Response) PutInt(i int) {
	r.data = append(r.data, NewValue(i))
}

//PutFloat ...
func (r *Response) PutFloat(f float64) {
	r.data = append(r.data, NewValue(f))
}

//PutBytes ...
func (r *Response) PutBytes(b []byte) {
	r.data = append(r.data, NewValue(b))
}

//Put ...
func (r *Response) Put(v Value) {
	r.data = append(r.data, v)
}

//PutValue ...
func (r *Response) PutValue(key string, val interface{}) {
	r.dataMap[key] = NewValue(val)
}

//Bytes ...
func (r *Response) Bytes() []byte {
	buf := Writer{}
	buf.PutUint16(int(r.id))
	if r.err != nil {
		buf.PutByte(DataErr)
		buf.PutString(r.err.Error())
		return buf.Bytes()
	}
	if r.err != nil {
		buf.PutByte(DataErr)
		buf.PutString(r.err.Error())
		return buf.Bytes()
	}
	buf.PutByte(DataSuccess)
	length := len(r.data) + len(r.dataMap)
	buf.PutUint16(length)
	for _, val := range r.data {
		buf.Put(toBytes(val))
	}
	for key, val := range r.dataMap {
		buf.PutByte(DataStringMap)
		if len(key) > math.MaxUint8 {
			key = key[:math.MaxUint8]
		}
		buf.PutShortString(key)
		buf.Put(toBytes(val))
	}
	return buf.Bytes()
}

func (r *Response) String() string {
	buf := strings.Builder{}

	buf.WriteString(fmt.Sprintf("Response ID: %d\n", r.id))
	for key, val := range r.dataMap {
		buf.WriteString(fmt.Sprintf("  * %s: %s\n", key, val.String()))
	}
	for key, val := range r.data {
		buf.WriteString(fmt.Sprintf("  - [%d] %s\n", key, val.String()))
	}
	if r.err != nil {
		buf.WriteString(`Err: ` + r.err.Error())
	}

	return buf.String()

}

func toBytes(data Value) []byte {
	b := Writer{}
	switch data := data.Raw().(type) {
	case []byte: //TODO handle long bytes
		b.PutByte(DataBytes)
		b.PutBytes(data)
	case int:
		b.PutByte(DataInt)
		b.PutInt(data)
	case float64:
		b.PutByte(DataFloat)
		b.PutFloat(data)
	case string: //TODO handle long string
		b.PutByte(DataString)
		b.PutString(data)
	default:
		println(`type not supported, %v`, data)
	}
	return b.Bytes()
}

//SetID ...
func (r *Response) SetID(id uint16) {
	r.id = id
}

//SetErr ...
func (r *Response) SetErr(e error) {
	r.err = e
}

//Get ...
func (r *Response) Get(idx int) Value {
	if idx >= len(r.data) {
		return NewValue(nil)
	}
	return r.data[idx]
}

//Init ...
func (r *Response) Init() *Response {
	r.dataMap = make(map[string]Value)
	return r
}

//Error ..
func (r *Response) Error() error {
	return r.err
}

//NewResponse ...
func NewResponse() *Response {
	r := &Response{}
	return r.Init()
}
