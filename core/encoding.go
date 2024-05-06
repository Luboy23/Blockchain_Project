package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
)

//	定义一个编码器接口，用于编码任意类型的数据。
type Encoder[T any] interface {
	Encode(T) error
}

//	定义了一个解码器接口，用于解码任意类型的数据。
type Decoder[T any] interface {
	Decode(T) error
}

//	实现Encoder接口，用于编码Transaction类型的数据。
type GobTxEncoder struct {
	w io.Writer	// 用于写入编码数据的io.Writer
}

//	创建一个新的GobTxEncoder实例。
func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	gob.Register(elliptic.P256()) // 注册椭圆曲线P256以支持其编码
	return &GobTxEncoder{w: w} // 返回新创建的GobTxEncoder实例
}

//	使用GobTxEncoder编码Transaction。
func (e *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(e.w).Encode(tx) // 使用gob编码器编码Transaction

}

//	实现了Decoder接口，用于解码Transaction类型的数据。
type GobTxDecoder struct {
	r io.Reader // 用于读取解码数据的io.Reader
}

//	创建一个新的GobTxDecoder实例。
func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256()) // 注册椭圆曲线P256以支持其解码
	return &GobTxDecoder{r: r} // 返回新创建的GobTxDecoder实例
}

//	使用GobTxDecoder解码Transaction。
func (e *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(e.r).Decode(tx) // 使用gob解码器解码Transaction

}
