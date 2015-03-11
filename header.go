// header
package dex

import (
	"bytes"
	"errors"
	"fmt"
)

var MAGIC = [8]byte{0x64, 0x65, 0x78, 0x0a, 0x30, 0x33, 0x35, 0x00}

type Header struct {
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

func (self *Header) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, "DEX FILE HEADER:")
	fmt.Fprintf(buf, "\tMAGIC: %x\n", MAGIC)
	fmt.Fprintf(buf, "\tchecksum: %d\n", self.checksum)
	fmt.Fprintf(buf, "\tsignature: %x\n", self.signature)
	fmt.Fprintf(buf, "\tfile_size: %d\n", self.file_size)
	fmt.Fprintf(buf, "\theader_size: %d\n", self.header_size)
	fmt.Fprintf(buf, "\tendian_tag: %d\n", self.endian_tag)

	fmt.Fprintf(buf, "\tlink_size: %d\n", self.link_size)
	fmt.Fprintf(buf, "\tlink_off: %d\n", self.link_off)

	fmt.Fprintf(buf, "\tmap_off: %d\n", self.map_off)

	fmt.Fprintf(buf, "\tstring_ids_size: %d\n", self.string_ids_size)
	fmt.Fprintf(buf, "\tstring_ids_off: %d\n", self.string_ids_off)

	fmt.Fprintf(buf, "\ttype_ids_size: %d\n", self.type_ids_size)
	fmt.Fprintf(buf, "\ttype_ids_off: %d\n", self.type_ids_off)

	fmt.Fprintf(buf, "\tproto_ids_size: %d\n", self.proto_ids_size)
	fmt.Fprintf(buf, "\tproto_ids_off: %d\n", self.proto_ids_off)

	fmt.Fprintf(buf, "\tfield_ids_size: %d\n", self.field_ids_size)
	fmt.Fprintf(buf, "\tfield_ids_off: %d\n", self.field_ids_off)

	fmt.Fprintf(buf, "\tmethod_ids_size: %d\n", self.method_ids_size)
	fmt.Fprintf(buf, "\tmethod_ids_off: %d\n", self.method_ids_off)

	fmt.Fprintf(buf, "\tclass_defs_size: %d\n", self.class_defs_size)
	fmt.Fprintf(buf, "\tclass_defs_off: %d\n", self.class_defs_off)

	fmt.Fprintf(buf, "\tdata_size: %d\n", self.data_size)
	fmt.Fprintf(buf, "\tdata_off: %d", self.data_off)

	return buf.String()
}

func readHeader(r BaseDalvikReader) *Header {
	h := new(Header)

	var magic [8]byte
	_, err := r.Bytes(magic[:])
	if err != nil {
		panic(err)
	}

	if magic != MAGIC {
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
