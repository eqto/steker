package service

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/eqto/steker/buff"
	"gitlab.com/tuxer/go-service"
)

var (
	funcMap      map[string]func(ctx Context) error
	regexPackage = regexp.MustCompile(`(?Uis)^.*/([a-z0-9]+)/[^/]+.go$`)
	listener     net.Conn
	running      = false

	writeLock = sync.Mutex{}

	logD = log.Println
	logW = log.Println
	logE = log.Println
)

//Value ...
type Value buff.Value

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
	logD(fmt.Sprintf(`Add function: %s()`, name))
	funcMap[name] = f
	return nil
}

//Stop ...
func Stop() {
	logD(`stopping...`)
	running = false
}

//Serve ..
func Serve() error {
	running = true
	buff := buff.NewReader(os.Stdin)

	reqCh := make(chan []byte)
	go func() {
		for running {
			data, e := buff.GetBytes()
			if e != nil {
				running = false
				reqCh <- nil
			} else {
				reqCh <- data
			}
		}
	}()

	ctx := &context{}

	service.New(func(exit <-chan int) {
		select {
		case <-exit:
			os.Stdin.Close()
		case data := <-reqCh:
			if data != nil {
				go processRequest(ctx, data)
			}
		}
	}, 1)
	if e := service.Start(); e != nil {
		return e
	}
	service.Wait()
	return nil
}

func processRequest(ctx *context, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			logE(r)
		}
	}()
	req, e := translateRequest(data)
	logD(req.String())
	if e != nil {
		panic(e)
	}
	if req.name != `` {
		if f, ok := funcMap[req.name]; ok {
			resp := buff.NewResponse()
			resp.SetID(req.id)
			ctx.req = req
			ctx.resp = resp
			if e := f(ctx); e != nil {
				resp.SetErr(e)
			}
			writePack(resp.Bytes())
		} else {
			panic(fmt.Sprintf(`func %s not found`, req.name))
		}
	} else {
		panic(`empty func name`)
	}
}

func writePack(data []byte) {
	writeLock.Lock()
	defer writeLock.Unlock()
	buf := buff.Writer{}
	buf.PutBytes(data)
	os.Stdout.Write(buf.Bytes())
}

func translateRequest(data []byte) (*request, error) {
	buf := buff.NewByteReader(data)

	id, e := buf.GetUint16()
	if e != nil {
		return nil, e
	}

	name, e := buf.GetShortString()
	if e != nil {
		return nil, e
	}
	req := &request{id: id, name: name}

	numVar, e := buf.GetUint16()
	if e != nil {
		return nil, e
	}
	for i := 0; i < int(numVar); i++ {
		typ, e := buf.GetByte()
		if e != nil {
			return nil, e
		}
		switch int(typ) {
		case buff.DataStringMap:
			key, e := buf.GetShortString()
			if e != nil {
				return nil, e
			}
			val, e := buf.GetBytes()
			if e != nil {
				return nil, e
			}
			req.putValue(key, val)
		case buff.DataBytes:
			b, e := buf.GetBytes()
			if e != nil {
				return nil, e
			}
			req.Put(b)
		case buff.DataInt:
			i, e := buf.GetInt()
			if e != nil {
				return nil, e
			}
			req.Put(i)
		case buff.DataFloat:
			f, e := buf.GetFloat()
			if e != nil {
				return nil, e
			}
			req.Put(f)
		case buff.DataString:
			s, e := buf.GetString()
			if e != nil {
				return nil, e
			}
			req.Put(s)
		}
	}
	return req, nil
}

//SetLogger ...
func SetLogger(d func(debug ...interface{}), w func(warn ...interface{}), e func(err ...interface{})) {
	logD = d
	logW = w
	logE = e
}
