package plugin

import (
	"errors"
	"io"
	"log"
	"math"
	"os/exec"
	"sync"
	"syscall"

	"github.com/eqto/steker/buff"
)

var (
	logD = log.Println
	logW = log.Println
	logE = log.Println
)

//Value ...
type Value buff.Value

//SetLogger ...
func SetLogger(d func(debug ...interface{}), w func(warn ...interface{}), e func(err ...interface{})) {
	logD = d
	logW = w
	logE = e
}

//Plugin ...
type Plugin interface {
	Stop() error
	Request(name string) Request
}

type plug struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser

	lockWriter sync.Mutex

	respMap map[uint16]chan Response

	counterLock sync.Mutex
	counter     uint16
	exitCh      chan bool
}

//Request ...
func (p *plug) Request(name string) Request {
	if len(name) > math.MaxUint16 {
		name = name[:math.MaxUint16]
	}
	return &request{plugin: p, name: name}
}

//Stop ...
func (p *plug) Stop() error {
	return p.cmd.Process.Signal(syscall.SIGINT)
}

func (p *plug) id() uint16 {
	p.counterLock.Lock()
	defer p.counterLock.Unlock()
	if p.counter == math.MaxInt16 {
		p.counter = 0
	}
	p.counter++
	return p.counter
}

func (p *plug) sendRequest(req *request) (<-chan Response, error) {
	b := buff.Writer{}
	id := p.id()
	b.PutUint16(int(id))
	m, e := req.bytes()
	if e != nil {
		return nil, e
	}
	b.Put(m)

	pack := buff.Writer{}
	pack.PutBytes(b.Bytes())

	if p.respMap == nil {
		p.respMap = make(map[uint16]chan Response)
	}
	respCh := make(chan Response, 1)
	p.respMap[id] = respCh

	_, e = p.write(pack.Bytes())
	if e != nil {
		delete(p.respMap, id)
		return nil, e
	}
	return respCh, nil
}

func (p *plug) write(data []byte) (int, error) {
	p.lockWriter.Lock()
	defer p.lockWriter.Unlock()
	return p.stdin.Write(data)
}

func (p *plug) read() {
	defer p.Stop()
	bufout := buff.NewReader(p.stdout)
	for {
		pack, e := bufout.GetBytes()
		if e != nil {
			break
		}
		buf := buff.NewByteReader(pack)
		id, e := buf.GetUint16()
		if e != nil {
			logD(`drop message, unable to get id`, e)
			continue
		}
		resp := &response{}
		resp.Init()
		resp.SetID(id)

		if ch, ok := p.respMap[id]; ok {
			status, e := buf.GetByte() //DataSuccess / DataErr
			if e != nil {
				logD(`drop message, unable to get status`, e)
				continue
			}
			switch status {
			case buff.DataSuccess:
				l, e := buf.GetUint16()
				if e != nil {
					logD(`drop message, unable to get length fields`, e)
					continue
				}
				length := int(l)
				for i := 0; i < length; i++ {
					typ, e := buf.GetByte()
					if e != nil {
						logD(`drop message, unable to get field type`, e)
						continue
					}
					var val interface{} = nil
					switch typ {
					case buff.DataStringMap:
						key, e := buf.GetShortString()
						if e != nil {
							logD(`drop message, unable to get map key`, e)
							continue
						}
						val, e := buf.GetData()
						if e != nil {
							logD(`drop message, unable to get byte val of map`, e)
							continue
						}
						resp.PutValue(key, val)
					case buff.DataBytes:
						v, e := buf.GetBytes()
						if e != nil {
							logD(`drop message, unable to get field bytes`, e)
							continue
						}
						val = v
					case buff.DataString:
						v, e := buf.GetString()
						if e != nil {
							logD(`drop message, unable to get field string`, e)
							continue
						}
						val = v
					case buff.DataInt:
						v, e := buf.GetInt()
						if e != nil {
							logD(`drop message, unable to get field int`, e)
							continue
						}
						val = v
					case buff.DataFloat:
						v, e := buf.GetFloat()
						if e != nil {
							logD(`drop message, unable to get field float`, e)
							continue
						}
						val = v
					}
					if val != nil {
						resp.Put(buff.NewValue(val))
					}
				}

			case buff.DataErr:
				str, e := buf.GetString()
				if e != nil {
					logD(`drop message, unable to get error string`, e)
					continue
				}
				resp.SetErr(errors.New(str))
			default:
				logD(`drop message, unable to identify message`)
				continue
			}
			logD(resp.String())
			ch <- resp
		}
	}
}

//LoadPlugin ...
func LoadPlugin(path string) (Plugin, error) {
	cmd := exec.Command(path)
	stdin, e := cmd.StdinPipe()
	if e != nil {
		return nil, e
	}
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		return nil, e
	}

	if e := cmd.Start(); e != nil {
		return nil, e
	}
	plug := &plug{cmd: cmd, stdin: stdin, stdout: stdout, exitCh: make(chan bool)}
	go plug.read()
	return plug, nil
}
