package utils

import (
	"encoding/hex"
)

/*
异或校验
*/
func XrCheckSum(data string) (string, error) {

	bytes, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	temp := byte(00)
	for i := 0; i < len(bytes); i++ {
		temp = temp ^ bytes[i]
	}
	res := ^temp
	return hex.EncodeToString([]byte{res}), nil
}
