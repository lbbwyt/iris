package uart

import (
	"encoding/hex"
	"errors"
	"github.com/jacobsa/go-serial/serial"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

/*
定义串口客户端
*/

type SerialClient struct {
	readWriter io.ReadWriteCloser
	OutBuffer  chan string
	PortName   string
	baudRate   uint
}

func NewSerialClient(portName string, baudRate uint) (*SerialClient, error) {

	c := &SerialClient{
		readWriter: nil,
		OutBuffer:  make(chan string, 1024),
		PortName:   portName,
		baudRate:   baudRate,
	}

	r, err := c.open(portName, baudRate)
	if err != nil {
		return nil, err
	}
	c.readWriter = r
	return c, nil
}

func (c *SerialClient) Start() {
	for {
		buf := make([]byte, 1024)
		n, err := c.readWriter.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Error("Error reading from serial port:  ", err)

				//尝试重新打开
				r, err := c.open(c.PortName, c.baudRate)
				if err != nil {
					time.Sleep(2 * time.Second)
					continue
				}
				c.readWriter = r
			}
			continue
		}

		buf = buf[:n]
		res := hex.EncodeToString(buf)
		if strings.TrimSpace(res) == "" {
			continue
		}
		log.Info("Rx: ", res)
		c.OutBuffer <- res
	}
}

func (c *SerialClient) open(portName string, baudRate uint) (io.ReadWriteCloser, error) {
	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        baudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	r, err := serial.Open(options)
	if err != nil {
		return nil, err
	}
	c.readWriter = r
	return r, err
}

func (c *SerialClient) SendMsg(data string) (int, error) {

	var (
		err error
	)
	if c.readWriter == nil {
		log.Warn("")
		return 0, errors.New("serial port is closed")
	}

	log.Info("send:" + data)
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return 0, err
	}

	return c.readWriter.Write(bytes)
}
