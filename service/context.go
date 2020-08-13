package service

//Context ...
type Context interface {
	// GetString(idx int) string
	// GetInt(idx int) int
	// GetFloat(idx int) float64
	// GetBytes(idx int) []byte

	setRequest(req Request)
}

type context struct {
	Context
	req Request
}

// func (c *context) GetString(idx int) string {
// 	return c.msg.GetString(idx)
// }
// func (c *context) GetInt(idx int) int {
// 	return c.msg.GetInt(idx)
// }
// func (c *context) GetFloat(idx int) float64 {
// 	return c.msg.GetFloat(idx)
// }
// func (c *context) GetBytes(idx int) []byte {
// 	return c.msg.GetBytes(idx)
// }

func (c *context) setRequest(req Request) {
	c.req = req
}
