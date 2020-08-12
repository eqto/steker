package plugin

import (
	"io"
	"log"
	"os/exec"
)

var (
	pluginMap     = make(map[string]*Plugin)
	debugLogger   = log.Println
	warningLogger = log.Println
	errorLogger   = log.Println
)

//Plugin ...
type Plugin struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

//Request ...
func (p *Plugin) Request() *Message {
	return &Message{plugin: p}
}

//Get ...
func Get(path string) *Plugin {
	if p, ok := pluginMap[path]; ok {
		return p
	}
	return nil
}

//SetLogger ...
func SetLogger(d func(debug ...interface{}), w func(warn ...interface{}), e func(err ...interface{})) {
	debugLogger = d
	warningLogger = w
	errorLogger = e
}

//Load ...
func Load(path string) (*Plugin, error) {
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
