package utils

import (
	"fmt"
	"testing"
)

func TestXrCheckSum(t *testing.T) {

	check := "55002a00000000004200200000000000000000000000000000000000000000000000557569645465737431"
	fmt.Println(XrCheckSum(check))

}
