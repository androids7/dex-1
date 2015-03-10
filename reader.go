// type
package dex

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

func _U2I(u uint32) int32 {
	return *((*int32)(unsafe.Pointer(&u)))
}

type DexReader interface {
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

	Utf8String() string
}

type dexReader struct {
	*bytes.Reader
	offset int
}

func New(b []byte) DexReader {
	return &dexReader{bytes.NewReader(b), 0}
}

func (self *dexReader) Byte() int8 {
	b, err := self.ReadByte()
	if err != nil {
		panic(err)
	}

	return *((*int8)(unsafe.Pointer(&b)))
}

func (self *dexReader) Ubyte() uint8 {
	b, err := self.ReadByte()
	if err != nil {
		panic(err)
	}

	return b
}

func (self *dexReader) Ushort() uint16 {
	b1, err := self.ReadByte()
	if err != nil {
		panic(err)
	}

	b2, err := self.ReadByte()
	if err != nil {
		panic(err)
	}

	return uint16(b1) | (uint16(b2) << 8)
}

func (self *dexReader) Short() int16 {
	ushort := self.Ushort()
	return *((*int16)(unsafe.Pointer(&ushort)))
}

func (self *dexReader) Uint() uint32 {
	low := self.Ushort()
	high := self.Ushort()

	return uint32(low) | (uint32(high) << 16)
}

func (self *dexReader) Int() int32 {
	i := self.Uint()
	return *((*int32)(unsafe.Pointer(&i)))
}

func (self *dexReader) Ulong() uint64 {
	low := self.Uint()
	high := self.Uint()

	return uint64(low) | (uint64(high) << 32)
}

func (self *dexReader) Long() int64 {
	l := self.Ulong()

	return *((*int64)(unsafe.Pointer(&l)))
}

func (self *dexReader) Sleb128() int32 {
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

func (self *dexReader) Uleb128() uint32 {
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

func (self *dexReader) Uleb128_p1() uint32 {
	return self.Uleb128() - 1
}

func (self *dexReader) Utf8String() string {
	var r []rune

	var ch rune
	var c, v byte

	for c = self.Ubyte(); c != 0; c = self.Ubyte() {
		if (c & 0x80) == 0x80 {
			if (c & 0xe0) == 0xc0 {
				c &= 0x1f
				v = self.Ubyte() & 0x3f
				ch = rune(c<<6) | rune(v)
			} else if (c & 0xf0) == 0xe0 {
				v = self.Ubyte() & 0x3f
				ch = rune(c<<6) | rune(v)
				v = self.Ubyte() & 0x3f
				ch = ch<<6 | rune(v)
			} else {
				panic(fmt.Errorf("Bad (point 4) UTF 8 %b", c))
			}
		} else {
			ch = rune(c)
		}

		r = append(r, ch)
	}

	return string(r)
}
