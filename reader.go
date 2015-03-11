// reader
package dex

import "io/ioutil"

type DexReader interface {
	HeaderInfo() *Header

	BaseReader() BaseDalvikReader
}

type dexReader struct {
	BaseDalvikReader
	*Header
}

func NewDexReader(b []byte) DexReader {
	r := &dexReader{NewBaseDalvikReader(b), nil}
	r.Header = readHeader(r.BaseDalvikReader)
	return r
}

func NewDexFileReader(filename string) (DexReader, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewDexReader(content), nil
}

func (self *dexReader) HeaderInfo() *Header {
	return self.Header
}

func (self *dexReader) BaseReader() BaseDalvikReader {
	return self.BaseDalvikReader
}
