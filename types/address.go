package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (a Address) ToSlice() []byte {
	b := make([]byte, 20)

	for i :=0; i <20; i++ {
		b[i] = a[i]
	}
	return b
}

func AddressFromBytes(b []byte)  Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("输入的字节长度应该为 20， 而不是 %d",len(b))
		panic(msg)
	}

	var value [20]uint8
	for i :=0; i <20; i++ {
		value[i] = b[i]
	}

	return Address(value)
}

func(a Address) String() string {
	return hex.EncodeToString(a.ToSlice())

}