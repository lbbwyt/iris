package uart

import (
	"iris/utils"
	"testing"
	"time"
)

func TestNewSerialClient(t *testing.T) {

	c, err := NewSerialClient("COM4", 1500000)
	if err != nil {
		panic(err)
	}
	go c.Start()

	c.SendMsg(utils.TrimHexStr("55 00 0A 00 00 00 00 00 42 00 00 E2 03"))

	go func() {
		tick := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-tick.C:
				data := "55 00 0E 00 00 00 00 00 7A 00 04 00 00 00 00 DA 03"
				c.SendMsg(utils.TrimHexStr(data))
			}
		}
	}()

	for {
		select {
		case res := <-c.OutBuffer:
			println("received data " + res)
		default:

		}
	}

}
