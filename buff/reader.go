package buff

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
)

//Reader ...
type Reader struct {
	r io.Reader
}

//GetByte get byte without header
func (r *Reader) GetByte() (byte, error) {
	b, e := r.GetLength(1)
	if e != nil {
		return 0, e
	}
	return b[0], nil
}

//Get ...
func (r *Reader) Get(b []byte) (int, error) {
	return r.r.Read(b)
}

//GetLength ...
func (r *Reader) GetLength(length int) ([]byte, error) {
	b := make([]byte, length)
	_, e := r.Get(b)
	if e != nil {
		return nil, e
	}
	return b, nil
}

//GetInt ...
func (r *Reader) GetInt() (int, error) {
	b, e := r.GetLength(8)
	if e != nil {
		return 0, e
	}
	return int(binary.BigEndian.Uint64(b)), nil
}

//GetFloat ...
func (r *Reader) GetFloat() (float64, error) {
	data, e := r.GetLength(8)
	if e != nil {
		return 0.0, e
	}
	return math.Float64frombits(binary.BigEndian.Uint64(data)), nil
}

//GetShortString ...
func (r *Reader) GetShortString() (string, error) {
	length, e := r.GetByte()
	if e != nil {
		return ``, e
	}
	data := make([]byte, length)
	_, e = r.r.Read(data)
	if e != nil {
		return ``, e
	}
	return string(data), nil
}

//GetString ...
func (r *Reader) GetString() (string, error) {
	l1, e := r.GetUint16()
	length := int(l1)
	if e != nil {
		return ``, e
	}
	if l1 >= math.MaxUint16 {
		l2, e := r.GetUint16()
		if e != nil {
			return ``, e
		}
		length += int(l2)
	}
	data, e := r.GetLength(length)
	if e != nil {
		return ``, e
	}
	return string(data), nil
}

//GetBytes use 2 or 4 bytes header for length
func (r *Reader) GetBytes() ([]byte, error) {
	l1, e := r.GetUint16()
	length := int(l1)
	if e != nil {
		return nil, e
	}
	if length >= math.MaxUint16 {
		l2, e := r.GetUint16()
		if e != nil {
			return nil, e
		}
		length += int(l2)
	}
	data := make([]byte, length)
	_, e = r.r.Read(data)
	if e != nil {
		return nil, e
	}
	return data, nil
}

//GetUint16 ...
func (r *Reader) GetUint16() (uint16, error) {
	data := make([]byte, 2)
	_, e := r.r.Read(data)
	if e != nil {
		return 0, e
	}
	return binary.BigEndian.Uint16(data), nil
}

//NewReader ...
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

//NewByteReader ...
func NewByteReader(b []byte) *Reader {
	return &Reader{r: bytes.NewReader(b)}
}
