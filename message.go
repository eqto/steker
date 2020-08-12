package plugin

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

const (
	mErr = iota
	mInt
	mFloat
	mString
)

//Message ...
type Message struct {
	buffer   bytes.Buffer
	plugin   *Plugin
	funcName string
}

//Add ...
func (m *Message) Add(data ...interface{}) {
	for _, d := range data {
		switch d := d.(type) {
		case string:
			length := make([]byte, 2)
			binary.BigEndian.PutUint16(length, uint16(len(d)))
			m.buffer.WriteByte(mString)
			m.buffer.Write(length)
			m.buffer.WriteString(d)
		case int:
			length := make([]byte, 8)
			binary.BigEndian.PutUint64(length, uint64(d))

			m.buffer.WriteByte(mInt)
			m.buffer.Write(length)
		case float64:
			length := make([]byte, 8)
			binary.BigEndian.PutUint64(length, math.Float64bits(d))

			m.buffer.WriteByte(mFloat)
			m.buffer.Write(length)
		default:
		}
	}

}

//Send ...
func (m *Message) Send(funcName string) error {
	m.funcName = funcName
	length := 1 + len(funcName) + m.buffer.Len() // 1 = alokasi panjang funcName
	if length > 65535 {
		return errors.New(`message too long`)
	}
	buff := new(bytes.Buffer)

	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, uint16(length))
	buff.Write(data)

	buff.WriteByte(byte(len(funcName)))
	buff.WriteString(funcName)

	buff.Write(m.buffer.Bytes())
	m.plugin.stdin.Write(buff.Bytes())
	return nil
}

//MustRead ...
func (m *Message) MustRead() Param {
	p, e := m.Read()
	if e == nil {
		return p
	}
	return Param{}
}
func (m *Message) Read() (Param, error) {
	if m.buffer.Len() == 0 {
		return Param{}, nil
	}
	t, e := m.buffer.ReadByte()
	if e != nil {
		m.buffer.Reset()
		return Param{}, e
	}
	switch t {
	case mString:
		data := make([]byte, 2)
		_, e := m.buffer.Read(data)
		if e != nil {
			m.buffer.Reset()
			return Param{}, e
		}
		length := binary.BigEndian.Uint16(data)
		data = make([]byte, length)
		m.buffer.Read(data)
		return Param{data: string(data)}, nil
	case mInt:
		data := make([]byte, 8)
		m.buffer.Read(data)
		return Param{data: int64(binary.BigEndian.Uint64(data))}, nil
	case mFloat:
		data := make([]byte, 8)
		m.buffer.Read(data)

		return Param{data: math.Float64frombits(binary.BigEndian.Uint64(data))}, nil
	}
	m.buffer.Reset()
	return Param{}, fmt.Errorf(`value type not recognized %d`, t)
}

func parseMessage(reader *bufio.Reader) (*Message, error) {
	data := make([]byte, 2)
	_, e := reader.Read(data)
	if e != nil {
		return nil, e
	}
	length := binary.BigEndian.Uint16(data)
	data = make([]byte, length)
	_, e = reader.Read(data)
	if e != nil {
		return nil, e
	}
	buff := bytes.NewBuffer(data)
	lenFunc, e := buff.ReadByte()
	if e != nil {
		return nil, e
	}
	data = make([]byte, lenFunc)
	buff.Read(data)

	m := &Message{funcName: string(data), buffer: *buff}

	return m, nil
}
