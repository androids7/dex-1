// dex_types
package dex

import (
	"bytes"
	"fmt"
)

type string_data_item struct {
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
