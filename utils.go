package main

import (
	"io"
	"fmt"
	"bytes"
	"net/url"
	"strings"
	"strconv"
	"crypto/rand"
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

func checkPost(m url.Values, keys ...string) bool {
	ok := true

	for _, key := range keys {
		if _, ok = m[key]; !ok {
			return false
		}
	}
	return true
}

func genUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}


func explodeRequest(str string) []string {
	var arr []string

	for i := 0; i < len(str) && str[i] != 0x03; {
		i2 := strings.Index(str[i:], ":")
		if i2 == -1 {
			fmt.Println("error", i, str[i:], len(str))
			break
		}
		i2 += i
		sl := str[i:i2]
		l, _ := strconv.Atoi(sl)
		i2++
		if (i2+l) < len(str) {
			arr = append(arr, str[i2:i2+l])
			i = i2 + l
		} else {
			arr = append(arr, str[i2:])
			i = len(str)
		}
	}
	return arr
}


func implodeRequest(arr []string) string {
	var str string

	for i := 0; i < len(arr); i++ {
		str += strconv.Itoa(len(arr[i]))+":"+arr[i]
	}
	return str
}
