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
	log "github.com/AlexStocks/log4go"
)

import (
	"github.com/AlexStocks/dubbogo/client"
	"github.com/AlexStocks/dubbogo/common"
	"github.com/AlexStocks/dubbogo/registry"
	"github.com/AlexStocks/dubbogo/selector"
)

var (
	connectTimeout  time.Duration = 100e6
	requestTimeout  time.Duration = 10e6
	survivalTimeout int           = 10e9
	rpcClient                     = make(map[client.CodecType]client.Client, 8)
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
		ok             bool
		err            error
		reqTimeout     time.Duration
		codecType      client.CodecType
		newRegistry    registry.NewRegistry
		newSelector    selector.NewSelector
		clientRegistry registry.Registry
		clientSelector selector.Selector
	)

	if conf == nil {
		panic(fmt.Sprintf("conf is nil"))
		return
	}

	// registry
	newRegistry, ok = client.DefaultRegistries[conf.Registry]
	if !ok {
		panic(fmt.Sprintf("illegal registry conf{%v}", conf.Registry))
		return
	}
	clientRegistry = newRegistry(
		registry.ApplicationConf(conf.Application_Config),
		registry.RegistryConf(conf.Registry_Config),
	)
	if clientRegistry == nil {
		panic("fail to init registry.Registy")
		return
	}
	for _, service := range conf.Service_List {
		err = clientRegistry.Register(service)
		if err != nil {
			panic(fmt.Sprintf("registry.Register(service{%#v}) = error{%v}", service, err))
			return
		}
	}

	// selector
	newSelector, ok = client.DefaultSelectors[conf.Selector]
	if !ok {
		panic(fmt.Sprintf("illegal selector conf{%v}", conf.Selector))
		return
	}
	clientSelector = newSelector(
		selector.Registry(clientRegistry),
		selector.SelectMode(selector.SM_RoundRobin),
	)
	if clientSelector == nil {
		panic(fmt.Sprintf("NewSelector(opts{registry{%#v}}) = nil", clientRegistry))
		return
	}

	// consumer
	reqTimeout, err = time.ParseDuration(conf.Request_Timeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Request_Timeout{%#v}) = error{%v}", conf.Request_Timeout, err))
		return
	}
	connectTimeout, err = time.ParseDuration(conf.Connect_Timeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Connect_Timeout{%#v}) = error{%v}", conf.Connect_Timeout, err))
		return
	}
	gxlog.CInfo("consumer retries:%d", conf.Retries)
	for idx := range conf.Service_List {
		codecType = client.GetCodecType(conf.Service_List[idx].Protocol)
		if codecType == client.CODECTYPE_UNKNOWN {
			panic(fmt.Sprintf("unknown protocol %s", conf.Service_List[idx].Protocol))
		}

		rpcClient[codecType] = client.NewClient(
			client.Retries(conf.Retries),
			client.RequestTimeout(reqTimeout),
			client.Registry(clientRegistry),
			client.Selector(clientSelector),
			client.ClientCodecType(codecType),
		)
	}
}

func uninitClient() {
	for k := range rpcClient {
		rpcClient[k].Close()
	}
	rpcClient = nil
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
			uninitClient()
			fmt.Println("app exit now...")
			return
		}
	}
}
