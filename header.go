// header
package dex

import (
	"bytes"
	"fmt"
)

var MAGIC = [8]byte{0x64, 0x65, 0x78, 0x0a, 0x30, 0x33, 0x35, 0x00}

const ENDIAN_CONSTANT = 0x12345678
const REVERSE_ENDIAN_CONSTANT = 0x78563412

type Header struct {
	Magic       [8]byte
	Checksum    uint32
	Signature   [20]byte
	File_size   uint32
	Header_size uint32
	Endian_tag  uint32

	Link_size uint32
	Link_off  uint32

	Map_off uint32

	String_ids_size uint32
	String_ids_off  uint32

	Type_ids_size uint32
	Type_ids_off  uint32

	Proto_ids_size uint32
	Proto_ids_off  uint32

	Field_ids_size uint32
	Field_ids_off  uint32

	Method_ids_size uint32
	Method_ids_off  uint32

	Class_defs_size uint32
	Class_defs_off  uint32

	Data_size uint32
	Data_off  uint32
}

func (self *Header) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, "DEX FILE HEADER:")
	fmt.Fprintf(buf, "\tMAGIC: %x\n", MAGIC)
	fmt.Fprintf(buf, "\tchecksum: %d\n", self.Checksum)
	fmt.Fprintf(buf, "\tsignature: %x\n", self.Signature)
	fmt.Fprintf(buf, "\tfile_size: %d\n", self.File_size)
	fmt.Fprintf(buf, "\theader_size: %d\n", self.Header_size)
	fmt.Fprintf(buf, "\tendian_tag: %x\n", self.Endian_tag)

	fmt.Fprintf(buf, "\tlink_size: %d\n", self.Link_size)
	fmt.Fprintf(buf, "\tlink_off: %d\n", self.Link_off)

	fmt.Fprintf(buf, "\tmap_off: %d\n", self.Map_off)

	fmt.Fprintf(buf, "\tstring_ids_size: %d\n", self.String_ids_size)
	fmt.Fprintf(buf, "\tstring_ids_off: %d\n", self.String_ids_off)

	fmt.Fprintf(buf, "\ttype_ids_size: %d\n", self.Type_ids_size)
	fmt.Fprintf(buf, "\ttype_ids_off: %d\n", self.Type_ids_off)

	fmt.Fprintf(buf, "\tproto_ids_size: %d\n", self.Proto_ids_size)
	fmt.Fprintf(buf, "\tproto_ids_off: %d\n", self.Proto_ids_off)

	fmt.Fprintf(buf, "\tfield_ids_size: %d\n", self.Field_ids_size)
	fmt.Fprintf(buf, "\tfield_ids_off: %d\n", self.Field_ids_off)

	fmt.Fprintf(buf, "\tmethod_ids_size: %d\n", self.Method_ids_size)
	fmt.Fprintf(buf, "\tmethod_ids_off: %d\n", self.Method_ids_off)

	fmt.Fprintf(buf, "\tclass_defs_size: %d\n", self.Class_defs_size)
	fmt.Fprintf(buf, "\tclass_defs_off: %d\n", self.Class_defs_off)

	fmt.Fprintf(buf, "\tdata_size: %d\n", self.Data_size)
	fmt.Fprintf(buf, "\tdata_off: %d", self.Data_off)

	return buf.String()
}
