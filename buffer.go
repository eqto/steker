package steker

import (
	"bytes"
	"encoding/binary"
	"errors"
)

//Buffer ...
type Buffer struct {
	data bytes.Buffer
}

// PutUint16 ...
func (b *Buffer) PutUint16(val uint16) error {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, val)
	if _, e := b.data.Write(data); e != nil {
		return e
	}
	return nil
}

// GetUint16 ...
func (b *Buffer) GetUint16() (uint16, error) {
	data := make([]byte, 2)

	if _, e := b.data.Read(data); e != nil {
		return 0, e
	}
	return binary.LittleEndian.Uint16(data), nil
}

//WriteShortString ...
func (b *Buffer) WriteShortString(s string) error {
	if len(s) > 255 {
		return errors.New(`maximum length of short string is 255 characters`)
	}
	if e := b.data.WriteByte(byte(len(s))); e != nil {
		return e
	}
	if _, e := b.data.Write([]byte(s)); e != nil {
		return e
	}
	return nil
}

//WriteString ...
func (b *Buffer) WriteString(s string) error {
	if len(s) > 65535 {
		return errors.New(`maximum length of string is 65535 characters`)
	}
	b.PutUint16(uint16(len(s)))
	b.data.Write([]byte(s))
	return nil
}

//ReadShortString ...
func (b *Buffer) ReadShortString() (string, error) {
	byteLen, e := b.data.ReadByte()
	if e != nil {
		return ``, e
	}
	len := uint8(byteLen)
	data := make([]byte, len)
	if _, e = b.data.Read(data); e != nil {
		return ``, e
	}
	return string(data), nil
}

//ReadString ...
func (b *Buffer) ReadString() (string, error) {
	data := make([]byte, 2)
	if _, e := b.data.Read(data); e != nil {
		return ``, e
	}
	byteLen, e := b.data.ReadByte()
	if e != nil {
		return ``, e
	}
	len := uint8(byteLen)
	data = make([]byte, len)
	_, e = b.data.Read(data)
	if e != nil {
		return ``, e

	}
	return string(data), nil
}
