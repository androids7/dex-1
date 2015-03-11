// string_items_test
package dex

import (
	"testing"
)

func TestStringItems(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e)
		}
	}()

	r, err := NewDexFileReader("test_data/bm_classes.dex", DETAIL)
	if err != nil {
		t.Fatal(err)
	}

	for i, item := range r.(*dexReader).string_data_items {
		t.Log(i, item.data)
	}
}
