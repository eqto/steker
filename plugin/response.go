package plugin

//Response ...
type Response struct {
	resp []interface{}
}

//PutString ...
func (r *Response) PutString(s string) {
	r.resp = append(r.resp, s)
}

//PutInt ...
func (r *Response) PutInt(i int) {
	r.resp = append(r.resp, i)
}

//PutFloat ...
func (r *Response) PutFloat(f float64) {
	r.resp = append(r.resp, f)
}

//PutBytes ...
func (r *Response) PutBytes(b []byte) {
	r.resp = append(r.resp, b)
}
