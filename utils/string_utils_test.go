package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestSubStr(t *testing.T) {
	s := "55002a00000000004200200000000000000000000000000000000000000000000055756964546573743434f903"

	fmt.Println(SubStr(s, 16, 2))
	fmt.Println(SubStr(s, 18, 4))

	fmt.Println(SubStr(s, 22, 64))
	bytes, _ := hex.DecodeString(SubStr(s, 22, 64))
	fmt.Println(string(bytes))
	fmt.Println(strings.TrimSpace(string(bytes)))
	fmt.Println(strings.Trim(string(bytes), ""))

}

func TestHexStrFill(t *testing.T) {
	fmt.Println(HexStrFill("dfasfa", 4))
}
