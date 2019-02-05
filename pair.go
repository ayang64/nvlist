package nvlist

// Pair encodes data about the name/value pair. The Size is the absolute size
// of the data and DecodedSize is the size after any decoding.
type Pair struct {
	Size        int32
	DecodedSize int32
}
