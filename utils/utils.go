package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)                        //开辟内存，存储字节集
	err := binary.Write(buff, binary.BigEndian, num) //num妆花字节集写入
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
