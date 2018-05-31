/******************************************************
# DESC    : client
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2016-07-29 15:50
# FILE    : client.go
******************************************************/

package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

import (
	"github.com/AlexStocks/goext/log"
	log "github.com/AlexStocks/log4go"
	jerrors "github.com/juju/errors"
)

import (
	"github.com/AlexStocks/dubbogo/client"
	"github.com/AlexStocks/dubbogo/codec"
	"github.com/AlexStocks/dubbogo/common"
	"github.com/AlexStocks/dubbogo/registry"
	"github.com/AlexStocks/dubbogo/selector"
)

type (
	User struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int64  `json:"age"`
		Time int64  `json:"time"`
		Sex  string `json:"sex"`
	}
)

var (
	connectTimeout  time.Duration = 100e6
	survivalTimeout int           = 10e9
	rpcClient                     = make(map[codec.CodecType]client.Client, 8)
)

func main() {
	var (
		err error
	)

	err = initClientConf()
	if err != nil {
		log.Error("initClientConf() = error{%#v}", jerrors.ErrorStack(err))
		return
	}
	initProfiling()
	initClient()

	go test()

	initSignal()
}

func initClient() {
	var (
		ok             bool
		err            error
		ttl            time.Duration
		reqTimeout     time.Duration
		codecType      codec.CodecType
		registryNew    RegistryNew
		selectorNew    SelectorNew
		clientRegistry registry.Registry
		clientSelector selector.Selector
	)

	if conf == nil {
		panic(fmt.Sprintf("conf is nil"))
		return
	}

	// registry
	registryNew, ok = DefaultRegistries[conf.Registry]
	if !ok {
		panic(fmt.Sprintf("illegal registry conf{%v}", conf.Registry))
		return
	}
	clientRegistry = registryNew(
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
			panic(fmt.Sprintf("registry.Register(service{%#v}) = error{%+v}", service, jerrors.ErrorStack(err)))
			return
		}
	}

	// selector
	selectorNew, ok = DefaultSelectors[conf.Selector]
	if !ok {
		panic(fmt.Sprintf("illegal selector conf{%v}", conf.Selector))
		return
	}
	clientSelector = selectorNew(
		selector.Registry(clientRegistry),
		selector.SelectMode(selector.SM_RoundRobin),
	)
	if clientSelector == nil {
		panic(fmt.Sprintf("NewSelector(opts{registry{%#v}}) = nil", clientRegistry))
		return
	}

	// consumer
	ttl, err = time.ParseDuration(conf.Pool_TTL)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Pool_TTL{%#v}) = error{%v}", conf.Pool_TTL, jerrors.ErrorStack(err)))
		return
	}
	reqTimeout, err = time.ParseDuration(conf.Request_Timeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Request_Timeout{%#v}) = error{%v}",
			conf.Request_Timeout, jerrors.ErrorStack(err)))
		return
	}
	connectTimeout, err = time.ParseDuration(conf.Connect_Timeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Connect_Timeout{%#v}) = error{%v}",
			conf.Connect_Timeout, jerrors.ErrorStack(err)))
		return
	}
	// ttl, err = (conf.Request_Timeout)
	gxlog.CInfo("consumer retries:%d", conf.Retries)
	for idx := range conf.Service_List {
		codecType = codec.GetCodecType(conf.Service_List[idx].Protocol)
		if codecType == codec.CODECTYPE_UNKNOWN {
			panic(fmt.Sprintf("unknown protocol %s", conf.Service_List[idx].Protocol))
		}

		rpcClient[codecType] = client.NewClient(
			client.Retries(conf.Retries),
			client.PoolSize(conf.Pool_Size),
			client.PoolTTL(ttl),
			client.RequestTimeout(reqTimeout),
			client.Registry(clientRegistry),
			client.Selector(clientSelector),
			client.CodecType(codecType),
		)
	}

	return
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

func test() {
	var (
		err     error
		idx     int
		service string
		method  string
		key     string
		ctx     context.Context
		req     client.Request
		wg      sync.WaitGroup
		sum     time.Duration
		diff    []time.Duration
		clt     client.Client
	)

	idx = -1
	service = string("com.ikurento.HelloService")
	for i := range conf.Service_List {
		if conf.Service_List[i].Service == service && conf.Service_List[i].Protocol == codec.CODECTYPE_JSONRPC.String() {
			idx = i
			break
		}
	}
	if idx == -1 {
		panic(fmt.Sprintf("can not find service in config service list:%#v", conf.Service_List))
	}

	key = conf.Test_String

	// Create request
	method = string("Echo")
	clt = rpcClient[codec.CODECTYPE_JSONRPC]
	req = clt.NewRequest(
		conf.Service_List[idx].Group,
		conf.Service_List[idx].Version,
		service,
		method,
		[]string{key},
	)

	// Set arbitrary headers in context
	ctx = context.WithValue(context.Background(), common.DUBBOGO_CTX_KEY, map[string]string{
		"X-Proxy-Id": "dubbogo",
		"X-Services": service,
		"X-Method":   method,
	})

	diff = make([]time.Duration, conf.Loop_Number, conf.Loop_Number)
	wg.Add(conf.Paral_Number)
	for idx = 0; idx < conf.Paral_Number; idx++ {
		go func(id int) {
			var (
				index int
				fail  int
				start time.Time
				rsp   string
			)
			// Call service
			start = time.Now()
			for index = 0; index < conf.Loop_Number; index++ {
				if err = clt.Call(ctx, req, &rsp, client.WithDialTimeout(connectTimeout)); err != nil {
					gxlog.CError("client.Call() return error:%s", err)
					fail++
					// return
				}
				if rsp != key {
					gxlog.CError("goroutine id:%d, client.Call(%s.%s{%s}) = {%s}", id, service, method, key, rsp)
					fail++
					// return
				}
				gxlog.CInfo("response result:%#v", rsp)
			}
			diff[id] = time.Now().Sub(start)
			gxlog.CInfo("after loop %d times, groutine{%d} time costs:%v, fail times:{%d}",
				conf.Loop_Number, id, diff[id].String(), fail)
			wg.Done()
			gxlog.CInfo("goroutine{%d} exit now", id)
		}(idx)
		gxlog.CInfo("loop index:%d", idx)
	}
	wg.Wait()

	for idx = 0; idx < conf.Loop_Number; idx++ {
		sum += diff[idx]
	}
	gxlog.CInfo("avg time diff:%s", time.Duration(int64(sum)/int64(conf.Loop_Number*conf.Paral_Number)).String())
}
