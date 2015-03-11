// dexreader_test
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

func TestTypeItems(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e)
		}
	}()

	r, err := NewDexFileReader("test_data/bm_classes.dex", DETAIL)
	if err != nil {
		t.Fatal(err)
	}

	for i, item := range r.(*dexReader).type_items {
		t.Log(i, item)
	}
}

func TestProtoItems(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e)
		}
	}()

	r, err := NewDexFileReader("test_data/bm_classes.dex", DETAIL)
	if err != nil {
		t.Fatal(err)
	}

	for i, item := range r.(*dexReader).proto_items {
		t.Log(i, item)
	}
}
