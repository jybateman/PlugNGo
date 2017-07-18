package main

import (
	// "fmt"
	"bytes"
)

func bytesSum(bs []byte) int {
	sum := 1
	for idx := range bs {
		sum += int(bs[idx])
	}
	return sum
}

func checkSum(bs []byte) int {
	s := bytesSum(bs)
	return s & 0xff
}

func CreateMessage(message []byte) []byte {
	var buf bytes.Buffer

	buf.Write(StartMessage)
	buf.WriteByte(byte(len(message)+1))
	buf.Write(message)
	buf.WriteByte(byte(checkSum(message)))
	buf.Write(EndMessage)	
	return buf.Bytes()
}
