package nvlist

// List encodes the version and flags (if any) of an nvlist.
type List struct {
	Version int32
	Flags   uint32
}
