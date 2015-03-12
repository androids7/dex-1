// reader
package dex

import (
	"encoding/binary"
	"io/ioutil"
)

type Flag int

const (
	HEADER_ONLY Flag = 0
	DETAIL      Flag = 1 << iota
)

type DexReader interface {
	HeaderInfo() *Header

	BaseReader() BaseDalvikReader
}

type dexReader struct {
	BaseDalvikReader
	*Header

	order binary.ByteOrder

	string_items []string_item
	type_items   []string
	proto_items  []proto_item
	field_items  []field_item
	method_items []method_item
}

func NewDexReader(b []byte, flag Flag) DexReader {
	r := &dexReader{BaseDalvikReader: NewBaseDalvikReader(b)}
	r.readHeader()

	if (flag & DETAIL) != 0 {
		r.parseDetail()
	}
	return r
}

func (self *dexReader) readHeader() {
	_, err := self.Seek(32, 0)
	if err != nil {
		panic(err)
	}

	var n uint32
	err = binary.Read(self, binary.LittleEndian, &n)
	if err != nil {
		panic(err)
	}

	if n == ENDIAN_CONSTANT {
		self.order = binary.LittleEndian
	} else {
		self.order = binary.BigEndian
	}

	self.Header = new(Header)
	_, err = self.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	err = binary.Read(self.BaseDalvikReader, binary.LittleEndian, self.Header)
	if err != nil {
		panic(err)
	}
}

func NewDexFileReader(filename string, flag Flag) (DexReader, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewDexReader(content, flag), nil
}

func (self *dexReader) HeaderInfo() *Header {
	return self.Header
}

func (self *dexReader) BaseReader() BaseDalvikReader {
	return self.BaseDalvikReader
}

func (self *dexReader) parseDetail() {
	self.parseStringItems()
	self.parseTypeItems()
	self.parseProtoItems()
	self.parseFieldItems()
	self.parseMethodItems()
}

func (self *dexReader) parseStringItems() {
	// parse string_data_item
	if self.String_ids_size > 0 {
		size := self.String_ids_size
		self.string_items = make([]string_item, 0, size)
		var item string_item
		var stringSize int
		var err error

		for i := uint32(0); i < size; i++ {
			_, err = self.Seek(int64(self.String_ids_off+i*4), 0)
			if err != nil {
				panic(err)
			}

			_, err = self.Seek(int64(self.Uint()), 0)
			if err != nil {
				panic(err)
			}

			item.utf16_size = self.Uleb128()
			item.data, stringSize = self.Utf8String()

			self.string_items = append(self.string_items, item)

			// TODO 比较 item.utf16_size 和 stringSize 是否相同？
			_ = stringSize
		}
	}
}

func (self *dexReader) parseTypeItems() {
	if self.Type_ids_size > 0 {
		_, err := self.Seek(int64(self.Type_ids_off), 0)
		if err != nil {
			panic(err)
		}

		size := self.Type_ids_size
		self.type_items = make([]string, 0, size)
		for i := uint32(0); i < size; i++ {
			self.type_items = append(self.type_items,
				self.string_items[self.Uint()].data)
		}
	}
}

func (self *dexReader) parseProtoItems() {
	if self.Proto_ids_size > 0 {
		size := self.Proto_ids_size
		self.proto_items = make([]proto_item, size)
		var params_size uint32
		var params_off uint32

		for i := uint32(0); i < size; i++ {
			_, err := self.Seek(int64(self.Proto_ids_off+i*12), 0)
			if err != nil {
				panic(err)
			}

			self.proto_items[i].shorty_desc = self.string_items[self.Uint()].data
			self.proto_items[i].return_type = self.type_items[self.Uint()]

			params_off = self.Uint()
			if params_off == 0 {
				continue
			}

			_, err = self.Seek(int64(params_off), 0)
			if err != nil {
				panic(err)
			}

			params_size = self.Uint()
			for j := uint32(0); j < params_size; j++ {
				self.proto_items[i].param_types = append(self.proto_items[i].param_types,
					self.type_items[self.Ushort()])
			}
		}
	}
}

func (self *dexReader) parseFieldItems() {
	if self.Field_ids_size > 0 {
		_, err := self.Seek(int64(self.Field_ids_off), 0)
		if err != nil {
			panic(err)
		}

		size := self.Field_ids_size
		self.field_items = make([]field_item, size)

		for i := uint32(0); i < size; i++ {
			self.field_items[i].class = self.type_items[self.Ushort()]
			self.field_items[i].type_name = self.type_items[self.Ushort()]
			self.field_items[i].name = self.string_items[self.Uint()].data
		}
	}
}

func (self *dexReader) parseMethodItems() {
	if self.Method_ids_size > 0 {
		_, err := self.Seek(int64(self.Method_ids_off), 0)
		if err != nil {
			panic(err)
		}

		size := self.Method_ids_size
		self.method_items = make([]method_item, size)

		for i := uint32(0); i < size; i++ {
			self.method_items[i].class = self.type_items[self.Ushort()]
			self.method_items[i].proto = self.proto_items[self.Ushort()]
			self.method_items[i].name = self.string_items[self.Uint()].data
		}
	}
}
