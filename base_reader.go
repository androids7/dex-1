// type
package dex

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf16"
	"unsafe"
)

func _U2I(u uint32) int32 {
	return *((*int32)(unsafe.Pointer(&u)))
}

func _B2Uint16(b byte) uint16 {
	return uint16(b)
}

type BaseDalvikReader interface {
	io.Seeker

	Byte() int8
	Ubyte() uint8

	Short() int16
	Ushort() uint16

	Int() int32
	Uint() uint32

	Long() int64
	Ulong() uint64

	Sleb128() int32
	Uleb128() uint32
	Uleb128_p1() uint32

	Utf8String() (string, int)

	Bytes([]byte) (int, error)

	Offset() int64
}

type baseDalvikReader struct {
	*bytes.Reader
	offset int64
}

func NewBaseDalvikReader(b []byte) BaseDalvikReader {
	return &baseDalvikReader{bytes.NewReader(b), 0}
}

func (self *baseDalvikReader) Byte() int8 {
	b, err := self.ReadByte()
	if err != nil {
		panic(err)
	}
	self.offset++

	return *((*int8)(unsafe.Pointer(&b)))
}

func (self *baseDalvikReader) Ubyte() uint8 {
	b, err := self.ReadByte()
	if err != nil {
		panic(err)
	}
	self.offset++

	return b
}

func (self *baseDalvikReader) Ushort() uint16 {
	b1, err := self.ReadByte()
	if err != nil {
		panic(err)
	}
	self.offset++

	b2, err := self.ReadByte()
	if err != nil {
		panic(err)
	}
	self.offset++

	return uint16(b1) | (uint16(b2) << 8)
}

func (self *baseDalvikReader) Short() int16 {
	ushort := self.Ushort()
	return *((*int16)(unsafe.Pointer(&ushort)))
}

func (self *baseDalvikReader) Uint() uint32 {
	low := self.Ushort()
	high := self.Ushort()

	return uint32(low) | (uint32(high) << 16)
}

func (self *baseDalvikReader) Int() int32 {
	i := self.Uint()
	return *((*int32)(unsafe.Pointer(&i)))
}

func (self *baseDalvikReader) Ulong() uint64 {
	low := self.Uint()
	high := self.Uint()

	return uint64(low) | (uint64(high) << 32)
}

func (self *baseDalvikReader) Long() int64 {
	l := self.Ulong()

	return *((*int64)(unsafe.Pointer(&l)))
}

func (self *baseDalvikReader) Sleb128() int32 {
	var r int32
	var v byte

	v = self.Ubyte()
	r = int32(v & 0x7f) // 清除高于 7 个 bit 的位

	if v >= 0x80 {
		v = self.Ubyte()
		r |= int32(v&0x7f) << 7 // 清除高于 7 个 bit 的位，然后左移 7 位
		if v >= 0x80 {
			v = self.Ubyte()
			r |= int32(v&0x7f) << 14 // 清除高于 7 个 bit 的位，然后左移 14 位
			if v >= 0x80 {
				v = self.Ubyte()
				r |= int32(v&0x7f) << 21 // 清除高于 7 个 bit 的位，然后左移 21 位
				if v >= 0x80 {
					v = self.Ubyte()
					r |= int32(v&0x7f) << 28 // 清除高于 7 个 bit 的位，然后左移 28 位
					if v >= 0x80 {
						panic(errors.New("Bad sleb128"))
					}
				} else if (v & 0x40) != 0 {
					r |= _U2I(0xf0000000)
				}
			} else if (v & 0x40) != 0 {
				r |= _U2I(0xffe00000)
			}
		} else if (v & 0x40) != 0 {
			r |= _U2I(0xffffc000)
		}
	} else if (v & 0x40) != 0 {
		r |= _U2I(0xffffff80)
	}

	return r
}

func (self *baseDalvikReader) Uleb128() uint32 {
	var r uint32
	var v byte

	v = self.Ubyte()
	r = uint32(v & 0x7f)

	if v >= 0x80 {
		v = self.Ubyte()
		r = uint32(v&0x7f) << 7
		if v >= 0x80 {
			v = self.Ubyte()
			r = uint32(v&0x7f) << 14
			if v >= 0x80 {
				v = self.Ubyte()
				r = uint32(v&0x7f) << 21
				if v >= 0x80 {
					v = self.Ubyte()
					r = uint32(v&0x7f) << 28
					if v >= 0x80 {
						panic(errors.New("Bad uleb128"))
					}
				}
			}
		}
	}

	return r
}

func (self *baseDalvikReader) Uleb128_p1() uint32 {
	return self.Uleb128() - 1
}

func (self *baseDalvikReader) Utf8String() (string, int) {
	var r []uint16

	var c, v uint16

	for c = _B2Uint16(self.Ubyte()); c != 0; c = _B2Uint16(self.Ubyte()) {
		if (c & 0x80) == 0x80 { // 字节最高位为 1
			if (c & 0xe0) == 0xc0 {
				c &= 0x1f
				v = _B2Uint16(self.Ubyte()) & 0x3f
				c = c<<6 | v
			} else if (c & 0xf0) == 0xe0 {
				v = _B2Uint16(self.Ubyte()) & 0x3f
				c = c<<6 | v
				v = _B2Uint16(self.Ubyte()) & 0x3f
				c = c<<6 | v
			} else {
				panic(fmt.Errorf("Bad (point 4) UTF 8 %b", c))
			}
		}

		r = append(r, c)
	}

	return string(utf16.Decode(r)), len(r)
}

//func (self *baseDalvikReader) Utf8StringDebug() (string, int) {
//	var r []uint16

//	var c, v uint16

//	for c = _B2Uint16(self.Ubyte()); c != 0; c = _B2Uint16(self.Ubyte()) {
//		if (c & 0x80) == 0x80 { // 字节最高位为 1
//			if (c & 0xe0) == 0xc0 {
//				c &= 0x1f
//				v = _B2Uint16(self.Ubyte()) & 0x3f
//				c = c<<6 | v
//			} else if (c & 0xf0) == 0xe0 {
//				v = _B2Uint16(self.Ubyte()) & 0x3f
//				c = c<<6 | v
//				v = _B2Uint16(self.Ubyte()) & 0x3f
//				c = c<<6 | v
//			} else {
//				panic(fmt.Errorf("Bad (point 4) UTF 8 %b", c))
//			}
//		}

//		if r != nil {
//			fmt.Printf(", ")
//		}
//		fmt.Printf("%d", c)
//		r = append(r, c)
//	}

//	return string(utf16.Decode(r)), len(r)
//}

func (self *baseDalvikReader) Bytes(b []byte) (int, error) {
	n, err := io.ReadFull(self, b)
	if n > 0 {
		self.offset += int64(n)
	}
	return n, err
}

func (self *baseDalvikReader) Offset() int64 {
	return self.offset
}
