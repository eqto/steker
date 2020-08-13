package buff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Writer struct {
	buff   bytes.Buffer
	writer io.Writer
}

//PutUint16 ...
func (w *Writer) PutUint16(i int) (int, error) {
	return w.Put(uint16len(i))
}

//Put ...
func (w *Writer) Put(b []byte) (int, error) {
	if w.writer == nil {
		return w.buff.Write(b)
	}
	return w.writer.Write(b)
}

//PutByte write byte without header
func (w *Writer) PutByte(b byte) error {
	if w.writer == nil {
		return w.buff.WriteByte(b)
	}
	_, e := w.writer.Write([]byte{b})
	return e
}

//PutInt ...
func (w *Writer) PutInt(data int) (int, error) {
	length := make([]byte, 8)
	binary.BigEndian.PutUint64(length, uint64(data))
	return w.Put(length)
}

//PutFloat ...
func (w *Writer) PutFloat(data float64) (int, error) {
	length := make([]byte, 8)
	binary.BigEndian.PutUint64(length, math.Float64bits(data))
	return w.Put(length)
}

//PutBytes use 2 or 4 bytes header for length
func (w *Writer) PutBytes(data []byte) (int, error) {
	length := 0
	if len(data) >= math.MaxUint16 {
		i1, e := w.Put(uint16len(math.MaxUint16))
		if e != nil {
			return i1, e
		}
		i2, e := w.Put(uint16len(len(data) - math.MaxUint16))
		if e != nil {
			return i1 + i2, e
		}
		length = i1 + i2
	} else {
		length, e := w.Put(uint16len(len(data)))
		if e != nil {
			return length, e
		}
	}
	l2, e := w.Put(data)
	return length + l2, e
}

//PutShortString ...
func (w *Writer) PutShortString(s string) (int, error) {
	if len(s) > math.MaxUint8 {
		return 0, fmt.Errorf(`maximum length of short string is %d characters`, math.MaxUint8)
	}
	e := w.PutByte(byte(len(s)))
	if e != nil {
		return 0, e
	}
	return w.Put([]byte(s))
}

//PutString ...
func (w *Writer) PutString(s string) (int, error) {
	length := 0
	if len(s) >= math.MaxUint16 {
		i1, e := w.Put(uint16len(math.MaxUint16))
		if e != nil {
			return i1, e
		}
		i2, e := w.Put(uint16len(len(s) - math.MaxUint16))
		length := i1 + i2
		if e != nil {
			return length, e
		}
	} else {
		i, e := w.Put(uint16len(len(s)))
		if e != nil {
			return i, e
		}
		length = i
	}
	i, e := w.Put([]byte(s))
	return length + i, e
}

//Bytes ...
func (w *Writer) Bytes() []byte {
	return w.buff.Bytes()
}
