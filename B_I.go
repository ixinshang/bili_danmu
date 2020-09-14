package bili_danmu

import (
	"bytes"
	"encoding/binary"

	p "github.com/qydysky/part"
)
/*
	整数 字节转换区
	32 4字节
	16 2字节
*/
func Itob32(num int32) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {p.Logf().E(err)}
	return buffer.Bytes()
}

func Itob16(num int16) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {p.Logf().E(err)}
	return buffer.Bytes()
}

func Btoi32(b []byte) int32 {
	var buffer int32
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &buffer)
	if err != nil {p.Logf().E(err)}
	return buffer
}

func Btoi16(b []byte) int16 {
	var buffer int16
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &buffer)
	if err != nil {p.Logf().E(err)}
	return buffer
}