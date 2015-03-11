// reader
package dex

import "io/ioutil"

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

	string_id_items   []uint32
	string_data_items []string_data_item
}

func NewDexReader(b []byte, flag Flag) DexReader {
	r := &dexReader{BaseDalvikReader: NewBaseDalvikReader(b)}
	r.Header = readHeader(r.BaseDalvikReader)

	if (flag & DETAIL) != 0 {
		r.parseDetail()
	}
	return r
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
}

func (self *dexReader) parseStringItems() {
	// parse string_data_item
	if self.Header.string_ids_size > 0 {
		_, err := self.Seek(int64(self.Header.string_ids_off), 0)
		if err != nil {
			panic(err)
		}

		// fill string_id_items
		size := self.Header.string_ids_size
		// 优化
		self.string_id_items = make([]uint32, 0, size)
		self.string_data_items = make([]string_data_item, 0, size)
		for i := uint32(0); i < size; i++ {
			self.string_id_items = append(self.string_id_items, self.Uint())
		} // end fill string_id_items

		// fill string_data_items
		for _, off := range self.string_id_items {
			_, err = self.Seek(int64(off), 0)
			if err != nil {
				panic(err)
			}

			var item string_data_item
			var stringSize int
			item.utf16_size = self.Uleb128()
			item.data, stringSize = self.Utf8String()

			self.string_data_items = append(self.string_data_items, item)

			_ = stringSize
		} // end fill string_data_items
	}
}
