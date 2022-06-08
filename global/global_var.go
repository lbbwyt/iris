package global

import (
	getty "github.com/AlexStocks/getty/transport"
	conf "iris/config"
	"iris/uart"
)

type GlobalVar struct {
	UartClient    *uart.SerialClient //串口
	Conn          getty.Session
	StopCycleChan chan struct{}
}

var GVars GlobalVar

func Init() error {
	client, err := uart.NewSerialClient(conf.GConfig.Uart.PortName, conf.GConfig.Uart.BaudRate)
	if err != nil {
		return err
	}
	go client.Start()
	GVars.UartClient = client
	GVars.StopCycleChan = make(chan struct{}, 4)
	return nil
}
