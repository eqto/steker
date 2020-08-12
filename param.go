package plugin

import (
	"strconv"
)

//Param ...
type Param struct {
	data interface{}
}

//IsNil ...
func (p Param) IsNil() bool {
	return p.data == nil
}

//Int ...
func (p Param) Int() int {
	switch data := p.data.(type) {
	case int64:
		return int(data)
	case string:
		int, _ := strconv.Atoi(data)
		return int
	case float64:
		return int(data)
	}
	return 0
}

//Float64 ...
func (p Param) Float64() float64 {
	switch data := p.data.(type) {
	case int64:
		return float64(data)
	case string:
		float, _ := strconv.ParseFloat(data, 64)
		return float
	case float64:
		return data
	}
	return 0
}

func (p Param) String() string {
	switch data := p.data.(type) {
	case int64:
		return strconv.FormatInt(data, 10)
	case string:
		return data
	case float64:
		return strconv.FormatFloat(data, 'f', -1, 64)
	}
	return ``
}
