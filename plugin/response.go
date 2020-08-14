package plugin

import "github.com/eqto/steker/buff"

//Response ...
type Response interface {
	Get(idx int) Value
	Error() error
}

type response struct {
	buff.Response
}

func (r *response) Get(idx int) Value {
	v := r.Response.Get(idx)
	return v
}
