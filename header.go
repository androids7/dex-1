// header
package dex

import (
	"errors"
)

var MAGIC = [8]byte{0x64, 0x65, 0x78, 0x0a, 0x30, 0x33, 0x35, 0x00}

type Header struct {
	magic       [8]byte
	checksum    uint32
	signature   [20]byte
	file_size   uint32
	header_size uint32
	endian_tag  uint32

	link_size uint32
	link_off  uint32

	map_off uint32

	string_ids_size uint32
	string_ids_off  uint32

	type_ids_size uint32
	type_ids_off  uint32

	proto_ids_size uint32
	proto_ids_off  uint32

	field_ids_size uint32
	field_ids_off  uint32

	method_ids_size uint32
	method_ids_off  uint32

	class_defs_size uint32
	class_defs_off  uint32

	data_size uint32
	data_off  uint32
}

func readHeader(r BaseDalvikReader) *Header {
	h := new(Header)

	_, err := r.Bytes(h.magic[:])
	if err != nil {
		panic(err)
	}

	if h.magic != MAGIC {
		panic(errors.New("Not MAGIC ahead of file"))
	}

	h.checksum = r.Uint()
	_, err = r.Bytes(h.signature[:])
	if err != nil {
		panic(err)
	}
	h.file_size = r.Uint()
	h.header_size = r.Uint()
	h.endian_tag = r.Uint()

	h.link_size = r.Uint()
	h.link_off = r.Uint()

	h.map_off = r.Uint()

	h.string_ids_size = r.Uint()
	h.string_ids_off = r.Uint()

	h.type_ids_size = r.Uint()
	h.type_ids_off = r.Uint()

	h.proto_ids_size = r.Uint()
	h.proto_ids_off = r.Uint()

	h.field_ids_size = r.Uint()
	h.field_ids_off = r.Uint()

	h.method_ids_size = r.Uint()
	h.method_ids_off = r.Uint()

	h.class_defs_size = r.Uint()
	h.class_defs_off = r.Uint()

	h.data_size = r.Uint()
	h.data_off = r.Uint()

	return h
}
