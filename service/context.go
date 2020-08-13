package service

import "github.com/eqto/cast"

//Context ...
type Context interface {
	GetString(idx int) string
	GetInt(idx int) int
	GetFloat(idx int) float64
	GetBytes(idx int) []byte

	setRequest(req Request)
}

type context struct {
	Context
	req Request
}

func (c *context) GetString(idx int) string {
	r := c.req.Get(idx)
	if r == nil {
		return ``
	}
	if s, e := cast.String(r); e == nil {
		return s
	}
	return ``
}

func (c *context) GetInt(idx int) int {
	r := c.req.Get(idx)
	if r == nil {
		return 0
	}
	if i, e := cast.Int(r); e == nil {
		return i
	}
	return 0
}
func (c *context) GetFloat(idx int) float64 {
	r := c.req.Get(idx)
	if r == nil {
		return 0.0
	}
	if f, e := cast.Float64(r); e == nil {
		return f
	}
	return 0.0
}
func (c *context) GetBytes(idx int) []byte {
	r := c.req.Get(idx)
	if r == nil {
		return nil
	}
	if b, e := cast.Bytes(r); e == nil {
		return b
	}
	return nil
}

func (c *context) setRequest(req Request) {
	c.req = req
}
