package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// 定义了一个长度为32的uint8数组，用于表示哈希值。
type Hash [32]uint8

// 检查Hash是否为零。
// 它遍历Hash中的每个元素，如果发现任何一个元素不为零，则返回false。
// 如果所有元素都为零，则返回true。
func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

// 将Hash转换为字节切片。
// 它创建一个新的字节切片，并将Hash中的每个元素复制到切片中。
func (h Hash) ToSlice() []byte {
	b := make([]byte, 32) // 创建一个新的字节切片，长度为32

	for i := 0; i < 32; i++ { // 遍历Hash中的每个元素
		b[i] = h[i] // 将元素复制到切片中
	}
	return b // 返回转换后的字节切片
}

// 将Hash转换为字符串。
// 它首先将Hash转换为字节切片，然后使用hex.EncodeToString将切片转换为十六进制字符串。
func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice()) // 返回十六进制字符串
}

// 从字节切片创建一个Hash。
// 它检查输入切片的长度是否为32，如果不是，则抛出panic。
// 如果长度正确，则创建一个新的Hash，并将切片中的元素复制到Hash中。
func HashFromBytes(b []byte) Hash {
	if len(b) != 32 { // 检查切片的长度
		msg := fmt.Sprintf("输入的字节长度应该为 32， 而不是 %d", len(b)) // 创建错误消息
		panic(msg) // 如果长度不正确，抛出panic
	}

	var value [32]uint8 // 创建一个新的uint8数组，长度为32
	for i := 0; i < 32; i++ { // 遍历切片中的每个元素
		value[i] = b[i] // 将元素复制到数组中
	}

	return Hash(value) // 将数组转换为Hash类型并返回
}

// 生成指定大小的随机字节切片。
// 它使用crypto/rand包的Read方法生成随机数据。
func RandomBytes(size int) []byte {
	token := make([]byte, size) // 创建一个新的字节切片，长度为指定的大小
	rand.Read(token) // 使用随机数据填充切片
	return token // 返回生成的随机字节切片
}

// 生成一个随机的Hash。
// 它首先生成一个随机的字节切片，然后使用这个切片创建一个Hash。
func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32)) // 返回生成的随机Hash
}
