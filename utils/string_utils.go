package utils

import (
	"fmt"
	"strings"
	"unsafe"
)

func TrimHexStr(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func SubStr(s string, index int, length int) string {

	var tail = index + length

	if index > len(s) {
		return ""
	}

	if (index + length) > len(s) {
		tail = len(s)
	}

	return s[index:tail]
}

func BytesToString(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

func StringToBytes(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}

func HexStrFill(Str string, num int) string {
	if len(Str) < num {
		fillNum := num - len(Str)

		fillStr := ""
		for i := 0; i < fillNum; i++ {
			fillStr = fmt.Sprintf("%s%s", "0", fillStr)
		}

		return fmt.Sprintf("%s%s", fillStr, Str)
	}

	return Str
}
