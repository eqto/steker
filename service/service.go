package service

import "log"

var (
	logD = log.Println
	logW = log.Println
	logE = log.Println
)

//SetLogger ...
func SetLogger(d func(debug ...interface{}), w func(warn ...interface{}), e func(err ...interface{})) {
	logD = d
	logW = w
	logE = e
}
