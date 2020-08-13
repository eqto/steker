package service

//Request ...
type Request interface {
	Put(data ...interface{})
	Name() string
	Get(idx int) interface{}

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

func (r *request) Get(idx int) interface{} {
	if idx > len(r.data) {
		return nil
	}
	return r.data[idx]
}
