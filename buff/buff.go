package buff

import "encoding/binary"

const (
	DataErr = iota
	DataByte
	DataInt
	DataFloat
	DataBytes
	DataString
)

//uint16len use 2 bytes
func uint16len(u int) []byte {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, uint16(u))
	return data
}

//uint16len use 4 bytes
func uint32len(u int) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(u))
	return data
}
