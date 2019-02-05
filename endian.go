package nvlist

import (
	"encoding/binary"
	"fmt"
)

// Endian encodes an byte ordervalue in a nvlist header.
type Endian int8

const (
	// LittleEndian denotes a little-endian byte order.
	LittleEndian = Endian(iota) // 0
	// BigEndian denotes a big-endian byte order.
	BigEndian // 1
)

// ByteOrder returns a binary.ByteOrder that corresponds with the Endian value.
func (e Endian) ByteOrder() binary.ByteOrder {
	switch e {
	case BigEndian:
		return binary.BigEndian
	case LittleEndian:
		return binary.LittleEndian
	}
	// unknown byte order value
	return nil
}

func (e Endian) String() string {
	switch e {
	case BigEndian:
		return "BigEndian"
	case LittleEndian:
		return "LittleEndian"
	}
	return fmt.Sprintf("*ERROR-INVALID-ENDIAN-VALUE-%02x*", int8(e))
}
