/******************************************************
# DESC    : client
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# LICENCE : Apache License 2.0
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2016-06-17 17:40
# FILE    : client.go
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
	"time"
)

import (
	"github.com/AlexStocks/goext/log"
	"github.com/AlexStocks/goext/net"
	log "github.com/AlexStocks/log4go"
)

import (
	"github.com/AlexStocks/dubbogo/client"
	"github.com/AlexStocks/dubbogo/registry"
	"github.com/AlexStocks/dubbogo/registry/zk"
	"github.com/AlexStocks/dubbogo/selector"
	"github.com/AlexStocks/dubbogo/selector/cache"
)

var (
	survivalTimeout int = 10e9
	clientSelector  selector.Selector
)

func main() {
	var (
		err error
	)

	err = initClientConfig()
	if err != nil {
		log.Error("initClientConfig() = error{%#v}", err)
		return
	}
	initProfiling()
	initClient()

	time.Sleep(1e9) // wait for selector

	gxlog.CInfo("\n\n\nstart to test jsonrpc")
	testJsonrpc("A003")

	initSignal()
}

func initClient() {
	var (
		err            error
		codecType      client.CodecType
		clientRegistry registry.Registry
	)

	if clientConfig == nil {
		panic(fmt.Sprintf("clientConfig is nil"))
		return
	}

	// registry
	clientRegistry = zookeeper.NewConsumerZookeeperRegistry(
		registry.ApplicationConf(clientConfig.Application_Config),
		registry.RegistryConf(clientConfig.Registry_Config),
	)
	if clientRegistry == nil {
		panic("fail to init registry.Registy")
		return
	}
	for _, service := range clientConfig.Service_List {
		err = clientRegistry.Register(service)
		if err != nil {
			panic(fmt.Sprintf("registry.Register(service{%#v}) = error{%v}", service, err))
			return
		}
	}

	// selector
	clientSelector = cache.NewSelector(
		selector.Registry(clientRegistry),
		selector.SelectMode(selector.SM_RoundRobin),
	)
	if clientSelector == nil {
		panic(fmt.Sprintf("NewSelector(opts{registry{%#v}}) = nil", clientRegistry))
		return
	}

	// consumer
	clientConfig.requestTimeout, err = time.ParseDuration(clientConfig.Request_Timeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Request_Timeout{%#v}) = error{%v}", clientConfig.Request_Timeout, err))
		return
	}
	clientConfig.connectTimeout, err = time.ParseDuration(clientConfig.Connect_Timeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Connect_Timeout{%#v}) = error{%v}", clientConfig.Connect_Timeout, err))
		return
	}

	for idx := range clientConfig.Service_List {
		codecType = client.GetCodecType(clientConfig.Service_List[idx].Protocol)
		if codecType == client.CODECTYPE_UNKNOWN {
			panic(fmt.Sprintf("unknown protocol %s", clientConfig.Service_List[idx].Protocol))
		}
	}
}

func uninitClient() {
	log.Close()
}

func initProfiling() {
	if !clientConfig.Pprof_Enabled {
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

	ip, err = gxnet.GetLocalIP()
	if err != nil {
		panic("cat not get local ip!")
	}
	addr = ip + ":" + strconv.Itoa(clientConfig.Pprof_Port)
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
			go time.AfterFunc(time.Duration(survivalTimeout)*time.Second, func() {
				log.Warn("app exit now by force...")
				os.Exit(1)
			})

			// 要么survialTimeout时间内执行完毕下面的逻辑然后程序退出，要么执行上面的超时函数程序强行退出
			uninitClient()
			fmt.Println("app exit now...")
			return
		}
	}
}
