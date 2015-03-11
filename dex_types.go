// dex_types
package dex

import (
	"bytes"
	"fmt"
)

type string_item struct {
	utf16_size uint32
	data       string
}

type proto_item struct {
	shorty_desc string // 用于重载判断
	return_type string
	param_types []string
}

func (self proto_item) String() string {
	param_list := bytes.NewBufferString("")

	for _, param := range self.param_types {
		param_list.WriteString(param)
	}

	return fmt.Sprintf("%s method_name(%s) => %s", self.return_type,
		param_list.String(), self.shorty_desc)
}

type field_item struct {
	class     string // 所属类
	type_name string
	name      string
}

func (self field_item) String() string {
	return fmt.Sprintf("%s %s#%s", self.type_name, self.class, self.name)
}

type method_item struct {
	class string
	proto proto_item
	name  string
}

func (self method_item) String() string {
	param_list := bytes.NewBufferString("")

	for _, param := range self.proto.param_types {
		param_list.WriteString(param)
	}

	return fmt.Sprintf("%s %s#%s(%s)", self.proto.return_type,
		self.class, self.name, param_list.String())
}
