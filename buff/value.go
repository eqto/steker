package buff

import (
	"errors"
	"fmt"
	"strconv"
)

//Value supported type: string, int, float64, []byte
type Value interface {
	String() string
	Int() (int, error)
	Float64() (float64, error)
	Bytes() []byte

	Raw() interface{}
}

type value struct {
	val interface{}
}

func (v *value) Raw() interface{} {
	return v.val
}

//String cast to string, supported type: string, float64, []byte
func (v *value) String() string {
	if v.val == nil {
		return ``
	}
	switch val := v.val.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	}
	if i, e := toInt(v.val); e == nil {
		return strconv.Itoa(i)
	} else if f, e := toFloat(v.val); e == nil {
		return fmt.Sprintf(`%f`, f)
	}

	return ``
}

//Int cast to int, supported type: string, int, float64
func (v *value) Int() (int, error) {
	if v.val == nil {
		return 0, nil
	}
	switch val := v.val.(type) {
	case string:
		return strconv.Atoi(val)
	case []byte:
		return strconv.Atoi(string(val))
	}
	if i, e := toInt(v.val); e == nil {
		return i, nil
	} else if f, e := toFloat(v.val); e == nil {
		return int(f), nil
	}
	return 0, nil
}

//Float64 cast to float64, supported type: string, int, float64
func (v *value) Float64() (float64, error) {
	if v.val == nil {
		return 0.0, nil
	}
	switch val := v.val.(type) {
	case string:
		return strconv.ParseFloat(val, 64)
	case []byte:
		return strconv.ParseFloat(string(val), 64)
	}
	if i, e := toInt(v.val); e == nil {
		return float64(i), nil
	} else if f, e := toFloat(v.val); e == nil {
		return f, nil
	}
	return 0.0, nil
}

//Bytes cast to []byte
func (v *value) Bytes() []byte {
	if v.val == nil {
		return nil
	}
	switch val := v.val.(type) {
	case string:
		return []byte(val)
	case []byte:
		return val
	}
	if i, e := toInt(v.val); e == nil {
		return []byte(strconv.Itoa(i))
	} else if f, e := toFloat(v.val); e == nil {
		return []byte(fmt.Sprintf(`%f`, f))
	}
	return nil
}

func toInt(val interface{}) (int, error) {
	i := 0
	switch val := val.(type) {
	case uint:
		i = int(val)
	case uint8:
		i = int(val)
	case uint16:
		i = int(val)
	case uint32:
		i = int(val)
	case uint64:
		i = int(val)
	case int:
		i = int(val)
	case int8:
		i = int(val)
	case int16:
		i = int(val)
	case int32:
		i = int(val)
	case int64:
		i = int(val)
	default:
		return 0, errors.New(`not int`)
	}
	return i, nil
}
func toFloat(val interface{}) (float64, error) {
	f := 0.0
	switch val := val.(type) {
	case float32:
		f = float64(val)
	case float64:
		f = val
	default:
		return 0, errors.New(`not float`)
	}
	return f, nil
}

//NewValue ...
func NewValue(val interface{}) Value {
	return &value{val: val}
}
