package plugin

import (
	"io"
	"log"
	"math"
	"os/exec"
	"sync"
	"syscall"

	"github.com/eqto/steker/buff"
)

var (
	pluginMap  = make(map[string]*Plugin)
	pluginLock = sync.Mutex{}
)

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

//Plugin ...
type Plugin struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser

	lockWriter sync.Mutex

	respMap map[uint16]chan Response

	counterLock sync.Mutex
	counter     uint16
}

//Request ...
func (p *Plugin) Request(name string) Request {
	if len(name) > math.MaxUint16 {
		name = name[:math.MaxUint16]
	}
	return &request{plugin: p, name: name}
}

//Stop ...
func (p *Plugin) Stop() error {
	return p.cmd.Process.Signal(syscall.SIGINT)
}

func (p *Plugin) id() uint16 {
	p.counterLock.Lock()
	defer p.counterLock.Unlock()
	if p.counter == math.MaxInt16 {
		p.counter = 0
	}
	p.counter++
	return p.counter
}

func (p *Plugin) sendRequest(r Request) (<-chan Response, error) {
	b := buff.Writer{}
	id := p.id()
	b.PutUint16(int(id))
	b.Put(r.bytes())

	logD(b.Bytes())

	pack := buff.Writer{}
	pack.PutBytes(b.Bytes())

	logD(pack.Bytes())
	if p.respMap == nil {
		p.respMap = make(map[uint16]chan Response)
	}
	respCh := make(chan Response, 1)
	p.respMap[id] = respCh

	_, e := p.write(pack.Bytes())
	if e != nil {
		delete(p.respMap, id)
		return nil, e
	}
	return respCh, nil
}

func (p *Plugin) write(data []byte) (int, error) {
	p.lockWriter.Lock()
	defer p.lockWriter.Unlock()
	return p.stdin.Write(data)
}

//Get ...
func Get(path string) (*Plugin, error) {
	pluginLock.Lock()
	defer pluginLock.Unlock()
	if plugin, ok := pluginMap[path]; ok {
		return plugin, nil
	}
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
	plugin := &Plugin{cmd: cmd, stdin: stdin, stdout: stdout}
	pluginMap[path] = plugin
	return plugin, nil
}
