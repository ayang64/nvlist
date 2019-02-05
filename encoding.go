package nvlist

import "fmt"

// Encoding is a type that stores the name/valure pair encoding type.  This can
// be either  EncodingNative or EncodingXDR.
//
// At the moment, the only encoding seen in the wild is XDR.
type Encoding uint8

const (
	// EncodingNative denotes the native nvlist value encoding. NEVER SEEN IN USE
	EncodingNative = Encoding(iota)
	// EncodingXDR denotes nvlist value encoding.  Default.
	EncodingXDR
)

func (e Encoding) String() string {
	switch e {
	case EncodingNative:
		return "EncodingNative"
	case EncodingXDR:
		return "EncodingXDR"
	}
	return fmt.Sprintf("*UNKNOWN-ENCODING-%02x*", uint8(e))
}
