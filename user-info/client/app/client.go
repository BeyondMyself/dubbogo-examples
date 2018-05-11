/******************************************************
# DESC    : client
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# LICENCE : LGPL V3
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2016-06-17 17:40
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
	"github.com/AlexStocks/dubbogo/codec/hessian"
	"github.com/AlexStocks/dubbogo/common"
	"github.com/AlexStocks/dubbogo/registry"
	"github.com/AlexStocks/dubbogo/selector"
)

var (
	connectTimeout  time.Duration = 100e6
	requestTimeout  time.Duration = 10e6
	survivalTimeout int           = 10e9
	rpcClient                     = make(map[codec.CodecType]client.Client, 8)
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

	time.Sleep(3e9) // wait for selector
	go testJsonrpc("A003")
	go testHessianGetUsers()

	initSignal()
}

func initClient() {
	var (
		ok             bool
		err            error
		ttl            time.Duration
		reqTimeout     time.Duration
		protocol       string
		codecType      codec.CodecType
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
	ttl, err = time.ParseDuration(conf.Pool_TTL)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Pool_TTL{%#v}) = error{%v}", conf.Pool_TTL, err))
		return
	}
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
		codecType = codec.GetCodecType(conf.Service_List[idx].Protocol)
		if codecType == codec.CODECTYPE_UNKNOWN {
			panic(fmt.Sprintf("unknown protocol %s", conf.Service_List[idx].Protocol))
		}

		gxlog.CInfo("start to build %s protocol client", codecType)
		rpcClient[codecType] = client.NewClient(
			client.Retries(conf.Retries),
			client.PoolSize(conf.Pool_Size),
			client.PoolTTL(ttl),
			client.RequestTimeout(reqTimeout),
			client.Registry(clientRegistry),
			client.Selector(clientSelector),
			client.CodecType(codecType),
		)
		gxlog.CInfo("protocol:%s, client:%+v", protocol, rpcClient[codecType])
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

func testJsonrpc(userKey string) {
	var (
		err error

		service string
		method  string
		user    *JsonRPCUser
		ctx     context.Context
		req     client.Request
		clt     client.Client
	)

	// Create request
	service = string("com.ikurento.user.UserProvider")
	method = string("GetUser")
	clt = rpcClient[codec.CODECTYPE_JSONRPC]
	req = clt.NewRequest(service, method, []string{userKey})
	// 注意这里，如果userKey是一个叫做UserKey类型的对象，则最后一个参数应该是 []UserKey{userKey}

	// Set arbitrary headers in context
	ctx = context.WithValue(context.Background(), common.DUBBOGO_CTX_KEY, map[string]string{
		"X-Proxy-Id": "dubbogo",
		"X-Services": service,
		"X-Method":   method,
	})

	user = new(JsonRPCUser)
	// Call service
	if err = clt.Call(ctx, req, user, client.WithDialTimeout(connectTimeout)); err != nil {
		gxlog.CError("client.Call() return error:%+v", jerrors.ErrorStack(err))
		return
	}

	log.Info("response result:%s", user)
	gxlog.CInfo("response result:%s", user)
}

func testHessianGetUsers() {
	var (
		err error

		service string
		method  string
		args    []interface{}
		ctx     context.Context
		req     client.Request
		rsp     []HessianUser
		clt     client.Client
	)

	hessian.RegisterJavaEnum(Gender(MAN))
	hessian.RegisterPOJO(&HessianUser{})

	// Create request
	service = string("com.ikurento.user.UserProvider")
	method = string("GetUsers")
	args = []interface{}{[]string{"001", "003", "004"}}

	clt = rpcClient[codec.CODECTYPE_DUBBO]
	req = clt.NewRequest(service, method, args)
	// Set arbitrary headers in context
	ctx = context.WithValue(context.Background(), common.DUBBOGO_CTX_KEY, map[string]string{
		"X-Proxy-Id": "dubbogo",
		"X-Services": service,
		"X-Method":   method,
	})

	// Call service
	if err = clt.Call(ctx, req, &rsp, client.WithDialTimeout(requestTimeout)); err != nil {
		gxlog.CError("client.Call() return error:%+v", jerrors.ErrorStack(err))
		return
	}

	log.Info("response result:%s", rsp)
	gxlog.CInfo("response result:%s", rsp)
}
