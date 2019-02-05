package nvlist

import "io"

// Read recursively parses the nvlist stored in the supplied
// io.Reader.  It is up to the caller to ensure that the
// reader is in position to start reading the nvlist.
func Read(r io.Reader) (map[string]interface{}, error) {
	rc := make(map[string]interface{})
	scn := NewScanner(r)
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
