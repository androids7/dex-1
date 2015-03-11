// header_test
package dex

import (
	"testing"
)

func TestHeader(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal("Header parse error!")
		}
	}()

	r, err := NewDexFileReader("test_data/bm_classes.dex", HEADER_ONLY)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(r.HeaderInfo())
}
