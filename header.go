package nvlist

// Header encodes an nvlist header which includeas an encoding type (either XDR
// or native) and the byte order of the values it stores.
type Header struct {
	Encoding  Encoding
	Endian    Endian
	Reserved1 int8
	Reserved2 int8
}
