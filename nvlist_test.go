package nvlist_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"testing"

	"github.com/ayang64/ztool/zfs/internal/nvlist"
)

func TestLooper(t *testing.T) {
	strOr := func(s ...string) string {
		for i := range s {
			if s[i] != "" {
				return s[i]
			}
		}
		return ""
	}
	volpath := strOr(os.Getenv("ZFSFILE"), "/obrovsky/recovery/zbackup0")
	fh, err := os.Open(volpath)
	if err != nil {
		t.Fatalf("could not open %q; %v -- must set ZFSFILE", volpath, err)
	}

	reader := func() io.Reader {
		switch path.Ext(volpath) {
		case ".cache":
			// concatinate a few bytes
			inf, err := os.Open(volpath)
			if err != nil {
				log.Fatal(err)
			}
			return io.MultiReader(bytes.NewReader([]byte{0x0, 0x1, 0x0, 0x0}), inf)
		default:
			// zfs nvlist is XDR encoded data that lives between 0x4000 - 0x20000 on the volume.
			// return a section reader covering that range.
			//
			// this could easily be:
			//
			// io.NewSectionReader(fh, 0x4000, 0x1c000)
			fh.Seek(0x4000, 0)
			return io.LimitReader(fh, 0x1c000)
		}
	}

	m, err := nvlist.Read(reader())

	if err != nil {
		t.Fatal(err)
	}

	o, err := json.MarshalIndent(m, "", "  ")

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n%s", string(o))
}
