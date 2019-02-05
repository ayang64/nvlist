package nvlist

// typedef struct nvpair {
// 	int32_t nvp_size;	/* size of this nvpair */
// 	int16_t	nvp_name_sz;	/* length of name string */
// 	int16_t	nvp_reserve;	/* not used */
// 	int32_t	nvp_value_elem;	/* number of elements for array types */
// 	data_type_t nvp_type;	/* type of value */
// 	/* name string */
// 	/* aligned ptr array for string arrays */
// 	/* aligned array of data for value */
// } nvpair_t;

// Type encodes the data type of a nvlist pair's value.
type Type int32

// nvlist data types
const (
	DontCare = Type(iota - 1)
	Unknown
	Boolean
	Byte
	Int16
	Uint16
	Int32
	Uint32
	Int64
	Uint64
	String
	ByteArray
	Int16Array
	Uint16Array
	Int32Array
	Uint32Array
	Int64Array
	Uint64Array
	StringArray
	HRTime
	NVList
	NVListArray
	BooleanValue
	Int8
	Uint8
	BooleanArray
	Int8Array
	Uint8Array
)

// Size returns the encoded size of a type. FIXME: This entire function may be unnesesary.
func (t Type) Size() int32 {
	sizes := map[Type]int32{
		DontCare:     -1, // We're not sure what DontCare types are and when they're used.
		Unknown:      -1,
		Boolean:      0, // Booleans can take up 0 space.  The existence of a key denotes a true value if the type is boolean.
		Byte:         1,
		Int16:        2,
		Uint16:       2,
		Int32:        4,
		Uint32:       4,
		Int64:        8,
		Uint64:       8,
		String:       1,
		ByteArray:    -1,
		Int16Array:   -1,
		Uint16Array:  -1,
		Int32Array:   -1,
		Uint32Array:  -1,
		Int64Array:   -1,
		Uint64Array:  -1,
		StringArray:  -1,
		HRTime:       -1,
		NVList:       -1,
		NVListArray:  -1,
		BooleanValue: -1,
		Int8:         1,
		Uint8:        1,
		BooleanArray: -1,
		Int8Array:    -1,
		Uint8Array:   -1,
	}

	if size, found := sizes[t]; found {
		return size
	}

	return -1
}

func (t Type) String() string {
	m := map[Type]string{
		Unknown:      "*UNKKNOWN-TYPE0*",
		Boolean:      "Boolean",
		Byte:         "Byte",
		Int16:        "Int16",
		Uint16:       "Uint16",
		Int32:        "Int32",
		Uint32:       "Uint32",
		Int64:        "Int64",
		Uint64:       "Uint64",
		String:       "String",
		ByteArray:    "ByteArray",
		Int16Array:   "Int16Array",
		Uint16Array:  "Uint16Array",
		Int32Array:   "Int32Array",
		Uint32Array:  "Uint32Array",
		Int64Array:   "Int64Array",
		Uint64Array:  "Uint64Array",
		StringArray:  "StringArray",
		HRTime:       "HRTime",
		NVList:       "NVList",
		NVListArray:  "NVListArray",
		BooleanValue: "BooleanValue",
		Int8:         "Int8",
		Uint8:        "Uint8",
		BooleanArray: "BooleanArray",
		Int8Array:    "Int8Array",
		Uint8Array:   "Uint8Array",
	}

	if str, found := m[t]; found {
		return str
	}

	return "*UNKNOWN-TYPE*"
}
