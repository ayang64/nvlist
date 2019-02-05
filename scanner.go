package nvlist

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type Scanner struct {
	byteOrder        binary.ByteOrder
	r                io.Reader
	header           Header
	list             List
	pair             Pair
	fieldName        string
	fieldType        Type
	fieldNumElements int
	value            interface{}
	err              error
	bytes            []byte
}

func NewScanner(r io.Reader) (rc *Scanner) {
	rc = &Scanner{
		r: r,
	}

	if err := binary.Read(r, binary.BigEndian, &rc.header); err != nil {
		rc.err = err
		return
	}

	rc.byteOrder = func() binary.ByteOrder {
		if rc.header.Endian == BigEndian {
			return binary.BigEndian
		}
		return binary.LittleEndian
	}()

	if err := binary.Read(r, rc.byteOrder, &rc.list); err != nil {
		rc.err = err
		return
	}
	return
}

func (s *Scanner) ReadValue(r io.Reader, t Type) (interface{}, error) {
	switch t {
	case DontCare:
	case Unknown:
	case Boolean:
		return true, nil

	case Byte:
	case Int16:
	case Uint16:
	case Int32:
	case Uint32:
	case Int64:
	case Uint64:
		var rc uint64
		if err := binary.Read(r, s.byteOrder, &rc); err != nil {
			return nil, err
		}
		return rc, nil

	case String:
		var length int32
		if err := binary.Read(r, s.byteOrder, &length); err != nil {
			return "", err
		}

		align4 := func(i int32) int32 {
			return (i + 3) & ^3
		}

		str := make([]byte, align4(length))

		if err := binary.Read(r, s.byteOrder, str); err != nil {
			return "", err
		}

		return string(str[:length]), nil

	case ByteArray:
	case Int16Array:
	case Uint16Array:
	case Int32Array:
	case Uint32Array:
	case Int64Array:
	case Uint64Array:
	case StringArray:
	case HRTime:
	case NVList:
		return s.ReadSub(r)

	case NVListArray:
		rc := make([]map[string]interface{}, 0, s.NumElements())
		for i := 0; i < s.NumElements(); i++ {
			v, err := s.ReadSub(r)

			if err != nil {
				log.Fatal(err)
			}
			rc = append(rc, v)
		}
		return rc, nil

	case BooleanValue:
	case Int8:
	case Uint8:
	case BooleanArray:
	case Int8Array:
	case Uint8Array:
	}

	return nil, fmt.Errorf("unknown type %q", t)
}

func (s *Scanner) ReadSub(r io.Reader) (map[string]interface{}, error) {
	rc := make(map[string]interface{})
	scn := s.NewSubScanner(r)
	if err := scn.Error(); err != nil {
		return nil, err
	}
	for scn.Next() {
		rc[scn.Name()] = scn.Value()
	}

	if err := scn.Error(); err != nil {
		return nil, err
	}

	return rc, nil
}

func (s *Scanner) ReadNumElements(r io.Reader) (int32, error) {
	var rc int32
	if err := binary.Read(r, s.byteOrder, &rc); err != nil {
		return -1, err
	}
	return rc, nil
}

func (s *Scanner) ReadType(r io.Reader) (Type, error) {
	var rc Type
	if err := binary.Read(r, s.byteOrder, &rc); err != nil {
		return -1, err
	}
	return rc, nil
}

func (s *Scanner) Bytes() []byte {
	return s.bytes
}

func (s *Scanner) ValueString() string {
	switch v := s.Value().(type) {
	case uint64:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	}
	return "*unrepresentable*"
}

func (s *Scanner) ReadString(r io.Reader) (string, error) {
	// read 4 byte length
	var length int32
	if err := binary.Read(r, s.byteOrder, &length); err != nil {
		return "", err
	}

	align4 := func(i int32) int32 {
		return (i + 3) & ^3
	}

	str := make([]byte, align4(length))

	if err := binary.Read(r, s.byteOrder, str); err != nil {
		return "", err
	}

	return string(str[:length]), nil
}

func (s *Scanner) Name() string {
	return s.fieldName
}

func (s *Scanner) Value() interface{} {
	return s.value
}

func (s *Scanner) Type() Type {
	return s.fieldType
}

func (s *Scanner) NumElements() int {
	return s.fieldNumElements
}

func (s *Scanner) Error() error {
	if s.err != nil {
		return s.err
	}
	return nil
}

func (s *Scanner) FieldSize() int {
	return int(s.pair.Size)
}

func (s *Scanner) Next() bool {
	// Do not continue if the scanner is in an errored state.
	if s.err != nil {
		return false
	}

	if s.err = binary.Read(s.r, s.byteOrder, &s.pair); s.err != nil {
		return false
	}

	if s.pair.Size == 0 && s.pair.DecodedSize == 0 {
		return false
	}

	// read entire record into a byte slice
	record := make([]byte, s.pair.Size-8)
	if s.err = binary.Read(s.r, s.byteOrder, record); s.err != nil {
		return false
	}

	// lets read from the remainding bytes
	br := bytes.NewReader(record)

	// read the name of the field
	name, err := s.ReadString(br)

	if err != nil {
		s.err = err
		return false
	}

	s.fieldName = name

	typ, err := s.ReadType(br)

	if err != nil {
		s.err = err
		return false
	}

	s.fieldType = typ

	nelements, err := s.ReadNumElements(br)

	if err != nil {
		s.err = err
		return false
	}

	s.fieldNumElements = int(nelements)

	value, err := s.ReadValue(br, s.fieldType)

	if err != nil {
		s.err = err
		return false
	}

	s.value = value

	return true
}

func (s *Scanner) NewSubScanner(r io.Reader) (rc *Scanner) {
	rc = &Scanner{
		r:         r,
		byteOrder: s.byteOrder,
	}

	// if err := binary.Read(r, rc.byteOrder, &rc.list); err != nil {
	if err := binary.Read(r, s.byteOrder, &rc.list); err != nil {
		rc.err = err
		return
	}

	return

}
