package service

//Request ...
type Request interface {
	Put(data ...interface{})
	Name() string

	bytes() []byte
}

type request struct {
	Request
	name string
	id   uint16

	data []interface{}
}

func (r *request) Put(data ...interface{}) {
	r.data = append(r.data, data)
}

func (r *request) Name() string {
	return r.name
}
