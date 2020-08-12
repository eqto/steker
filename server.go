package plugin

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	funcMap      map[string]func(ctx Context) error
	regexPackage = regexp.MustCompile(`(?Uis)^.*/([a-z0-9]+)/[^/]+.go$`)
)

//AddFunc ...
func AddFunc(f func(ctx Context) error) error {
	ptr := reflect.ValueOf(f).Pointer()
	name := runtime.FuncForPC(ptr).Name()
	if strings.Count(name, `.`) > 1 {
		return errors.New(`unsupported add inline function`)
	}
	name = name[strings.IndexRune(name, '.')+1:]

	if funcMap == nil {
		funcMap = make(map[string]func(ctx Context) error)
	}
	debugLogger(fmt.Sprintf(`Add function: %s()`, name))
	funcMap[name] = f
	return nil
}

//ServeUnix serve unix socket (not implemented yet)
func ServeUnix() {

}

//Serve ...
func Serve() {
	reader := bufio.NewReader(os.Stdin)

	var ctx Context

	for {
		req, e := parseMessage(reader)
		if e != nil {
			errorLogger(e)
			time.Sleep(1 * time.Second)
		} else {
			if f, ok := funcMap[req.funcName]; ok {
				params := []Param{}
				for {
					param, e := req.Read()
					if e != nil {
						errorLogger(e)
						break
					} else {
						if param.IsNil() {
							break
						}
						params = append(params, param)
					}
				}
				ctx.params = params

				if e := f(ctx); e != nil {
					debugLogger(e)
				}
			} else {
				warningLogger(fmt.Sprintf(`unable to find function named %s`, req.funcName))
			}
		}
	}
}
