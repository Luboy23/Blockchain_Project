package types

import (
	"encoding/hex"
	"fmt"
)

// 定义了一个长度为20的uint8数组，用于表示以太坊地址。
type Address [20]uint8

// 将Address类型转换为字节切片。
// 它创建一个新的字节切片，并将Address中的每个元素复制到切片中。
func (a Address) ToSlice() []byte {
	b := make([]byte, 20) // 创建一个新的字节切片，长度为20

	for i := 0; i < 20; i++ { // 遍历Address中的每个元素
		b[i] = a[i] // 将元素复制到切片中
	}
	return b // 返回转换后的字节切片
}

// 从字节切片创建一个Address。
// 它检查输入切片的长度是否为20，如果不是，则抛出panic。
// 如果长度正确，则创建一个新的Address，并将切片中的元素复制到Address中。
func AddressFromBytes(b []byte) Address {
	if len(b) != 20 { // 检查切片的长度
		msg := fmt.Sprintf("输入的字节长度应该为 20， 而不是 %d", len(b)) // 创建错误消息
		panic(msg) // 如果长度不正确，抛出panic
	}

	var value [20]uint8 // 创建一个新的uint8数组，长度为20
	for i := 0; i < 20; i++ { // 遍历切片中的每个元素
		value[i] = b[i] // 将元素复制到数组中
	}

	return Address(value) // 将数组转换为Address类型并返回
}

// 将Address转换为字符串。
// 它首先将Address转换为字节切片，然后使用hex.EncodeToString将切片转换为十六进制字符串。
func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice()) // 返回十六进制字符串
}
