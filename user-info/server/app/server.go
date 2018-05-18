/******************************************************
# DESC    : provider example
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# LICENCE : Apache License 2.0
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2016-07-21 16:41
# FILE    : server.go
******************************************************/

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

import (
	// "github.com/AlexStocks/gocolor"
	log "github.com/AlexStocks/log4go"
)

import (
	"github.com/AlexStocks/dubbogo/codec"
	"github.com/AlexStocks/dubbogo/common"
	"github.com/AlexStocks/dubbogo/registry"
	"github.com/AlexStocks/dubbogo/server"
	"github.com/AlexStocks/dubbogo/transport"
)

var (
	survivalTimeout int = 3e9
	servo           server.Server
)

func main() {
	var (
		err error
	)

	err = configInit()
	if err != nil {
		log.Error("configInit() = error{%#v}", err)
		return
	}
	initProfiling()

	servo = initServer()
	err = servo.Handle(&UserProvider{})
	if err != nil {
		panic(err)
		return
	}
	servo.Start()

	initSignal()
}

func initServer() server.Server {
	var (
		ok              bool
		protocol        string
		contentType     string
		cdcNew          codec.NewCodec
		registryNew     RegistryNew
		transportNew    TransportNew
		protocolMap     map[string]string
		codecs          map[string]codec.NewCodec
		serverRegistry  registry.Registry
		serverTransport transport.Transport
		srv             server.Server
	)

	if conf == nil {
		panic(fmt.Sprintf("conf is nil"))
		return nil
	}

	// registry
	registryNew, ok = DefaultRegistries[conf.Registry]
	if !ok {
		panic(fmt.Sprintf("illegal registry conf{%v}", conf.Registry))
		return nil
	}
	serverRegistry = registryNew(
		registry.ApplicationConf(conf.Application_Config),
		registry.RegistryConf(conf.Registry_Config),
	)
	if serverRegistry == nil {
		panic("fail to init registry.Registy")
		return nil
	}

	// transport
	transportNew, ok = DefaultTransports[conf.Transport]
	if !ok {
		panic(fmt.Sprintf("illegal transport conf{%v}", conf.Transport))
		return nil
	}
	serverTransport = transportNew()
	if serverTransport == nil {
		panic(fmt.Sprintf("TransportNew(type{%s}) = nil", conf.Transport))
		return nil
	}

	// codec
	protocolMap = make(map[string]string, len(conf.Server_List))
	for _, svrConf := range conf.Server_List {
		protocolMap[svrConf.Protocol] = svrConf.Protocol
	}
	if len(protocolMap) == 0 {
		panic("server list is nil")
		return nil
	}
	codecs = make(map[string]codec.NewCodec, len(protocolMap))
	for protocol = range protocolMap {
		if contentType, ok = DefaultContentTypes[protocol]; !ok {
			// 如果将来使用tcp协议，也可以为拟定一个伪contentType
			panic(fmt.Sprintf("can not find content type of protocol{%s}", protocol))
			return nil
		}

		if cdcNew, ok = DefaultCodecs[protocol]; !ok {
			panic(fmt.Sprintf("can not find codec of protocol{%s}", protocol))
			return nil
		}
		codecs[contentType] = cdcNew
	}

	// provider
	srv = server.NewServer(
		server.Codec(codecs),
		server.Registry(serverRegistry),
		server.Transport(serverTransport),
		server.ServerConfList(conf.Server_List),
		server.ServiceConfList(conf.Service_List),
	)

	return srv
}

func uninitServer() {
	if servo != nil {
		servo.Stop()
	}
	log.Close()
}

func initProfiling() {
	if !conf.Pprof_Enabled {
		return
	}
	const (
		PprofPath = "/debug/pprof/"
	)
	var (
		err  error
		ip   string
		addr string
	)

	ip, err = common.GetLocalIP(ip)
	if err != nil {
		panic("cat not get local ip!")
	}
	addr = ip + ":" + strconv.Itoa(conf.Pprof_Port)
	log.Info("App Profiling startup on address{%v}", addr+PprofPath)

	go func() {
		log.Info(http.ListenAndServe(addr, nil))
	}()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		log.Info("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
		// reload()
		default:
			go common.Future(survivalTimeout, func() {
				log.Warn("app exit now by force...")
				os.Exit(1)
			})

			// 要么survialTimeout时间内执行完毕下面的逻辑然后程序退出，要么执行上面的超时函数程序强行退出
			uninitServer()
			fmt.Println("provider app exit now...")
			return
		}
	}
}
