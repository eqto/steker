package service

//Response ...
type Response interface {
	PutString(s string)
	PutInt(i int)
	PutFloat(f float64)
	PutBytes(b []byte)
	PutValue(key string, val interface{})
}
