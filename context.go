package plugin

import "errors"

//Context ...
type Context struct {
	params []Param
}

//Params ...
func (c *Context) Params() []Param {
	return c.params
}

//Param ...
func (c *Context) Param(i int) (Param, error) {
	if i < len(c.params) {
		return c.params[i], nil
	}
	return Param{}, errors.New(`param not found`)
}

//MustParam ...
func (c *Context) MustParam(i int) Param {
	p, _ := c.Param(i)
	return p
}
