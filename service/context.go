package service

//Context ...
type Context interface {
	Request() Request
	Response() Response
}

type context struct {
	Context
	req  Request
	resp Response
}

func (c *context) Request() Request {
	return c.req
}

func (c *context) Response() Response {
	return c.resp
}
