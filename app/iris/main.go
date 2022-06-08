package main

import (
	"flag"
	"fmt"
	getty "github.com/AlexStocks/getty/transport"
	gxnet "github.com/AlexStocks/goext/net"
	log "github.com/sirupsen/logrus"
	conf "iris/config"
	"iris/global"
	"iris/tcp/handler"
	"iris/web/router"
	"net"
	"net/http"
	"time"
)

var (
	cfgPath        string
	irisMsgHandler = handler.NewEchoMessageHandler()
	irisPkgHandler = handler.NewEchoPackageHandler()
)

func init() {
	flag.StringVar(&cfgPath, "c", "./etc/config.yaml", "")
	flag.Parse()
}

func main() {
	println("starting.....")

	err := conf.Init(cfgPath)
	if err != nil {
		log.WithFields(log.Fields{"cfg_path": cfgPath}).WithError(err).Error("[main] config init error")
		return
	}

	log.Infof("conf:%v", conf.GConfig)

	//初始化全局变量
	global.Init()

	go startTcpServer(conf.GConfig.TcpServer.Port)

	router, err := router.InitRouter()
	if err != nil {
		log.WithError(err).Error("[main] init router error")
		return
	}
	server := &http.Server{
		Addr:           ":" + conf.GConfig.WebServer.Port,
		Handler:        router,
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.WithError(err).Error("[main] start web server error")
		return
	}
}

/**
启动tcpServer
*/
func startTcpServer(port string) {
	addr := gxnet.HostAddress2("localhost", port)
	server := getty.NewTCPServer(
		getty.WithLocalAddress(addr),
	)
	// run server
	server.RunEventLoop(newSession)
	log.Info("server bind addr{%s} ok!", port)

}

func newSession(session getty.Session) error {
	var (
		ok      bool
		tcpConn *net.TCPConn
	)

	if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
	}

	tcpConn.SetKeepAlive(true)

	session.SetPkgHandler(irisPkgHandler)
	session.SetEventListener(irisMsgHandler)

	session.SetCronPeriod(1000 * 10)
	log.Debug("client new session:%s\n", session.Stat())

	return nil
}
