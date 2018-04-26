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
	"github.com/AlexStocks/goext/time"
	log "github.com/AlexStocks/log4go"
)

import (
	"github.com/AlexStocks/dubbogo/client"
	"github.com/AlexStocks/dubbogo/common"
	"github.com/AlexStocks/dubbogo/registry"
	"github.com/AlexStocks/dubbogo/selector"
	"github.com/AlexStocks/dubbogo/transport"
)

type Gender int

const (
	MAN = iota
	WOMAN
)

var genderStrings = [...]string{
	"MAN",
	"WOMAN",
}

func (g Gender) String() string {
	return genderStrings[g]
}

type (
	User struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int64  `json:"age"`
		Time int64  `json:"time"`
		Sex  string `json:"sex"`
		// Sex Gender `json:"sex"`
	}
)

func (u User) String() string {
	return fmt.Sprintf(
		"User{Id:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.Id, u.Name, u.Age, gxtime.YMDPrint(int(u.Time), 0), u.Sex,
	)
}

var (
	connectTimeout  time.Duration = 100e6
	survivalTimeout int           = 10e9
	rpcClient       map[string]client.Client
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

	go testJsonrpc("A003")
	go testJsonrpc("A000")

	initSignal()
}

func initClient() {
	var (
		ok              bool
		err             error
		ttl             time.Duration
		reqTimeout      time.Duration
		protocol        string
		registryNew     RegistryNew
		selectorNew     SelectorNew
		transportNew    TransportNew
		clientRegistry  registry.Registry
		clientSelector  selector.Selector
		clientTransport transport.Transport
		cltConfig       DubbogoClientConfig
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
			panic(fmt.Sprintf("registry.Register(service{%#v}) = error{%v}", service, err))
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
	// ttl, err = (conf.Request_Timeout)
	gxlog.CInfo("consumer retries:%d", conf.Retries)
	for idx := range conf.Service_List {
		protocol = conf.Service_List[idx].Protocol
		cltConfig = DefaultDubbogoClientConfig[protocol]

		// transport
		transportNew, ok = DefaultTransports[cltConfig.transportType]
		if !ok {
			panic(fmt.Sprintf("illegal transport conf{%v}", cltConfig.transportType))
			return
		}
		clientTransport = transportNew()
		if clientTransport == nil {
			panic(fmt.Sprintf("TransportNew(type{%s}) = nil", cltConfig.transportType))
			return
		}

		gxlog.CError("start to build %s protocol client", protocol)
		rpcClient[protocol] = client.NewClient(
			client.Retries(conf.Retries),
			client.PoolSize(conf.Pool_Size),
			client.PoolTTL(ttl),
			client.RequestTimeout(reqTimeout),
			client.Registry(clientRegistry),
			client.Selector(clientSelector),
			client.Transport(clientTransport),
			client.Codec(DefaultContentTypes[cltConfig.contentType], cltConfig.codec),
			client.ContentType(cltConfig.contentType),
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

func testJsonrpc(userKey string) {
	var (
		err error

		service string
		method  string
		user    *User
		ctx     context.Context
		req     client.Request
		clt     client.Client
	)

	// Create request
	service = string("com.ikurento.user.UserProvider")
	method = string("GetUser")
	clt = rpcClient["jsonrpc"]
	req = clt.NewJsonRequest(service, method, []string{userKey})
	// 注意这里，如果userKey是一个叫做UserKey类型的对象，则最后一个参数应该是 []UserKey{userKey}

	// Set arbitrary headers in context
	ctx = context.WithValue(context.Background(), common.DUBBOGO_CTX_KEY, map[string]string{
		"X-Proxy-Id": "dubbogo",
		"X-Services": service,
		"X-Method":   method,
	})

	user = new(User)
	// Call service
	if err = clt.Call(ctx, req, user, client.WithDialTimeout(connectTimeout)); err != nil {
		gxlog.CError("client.Call() return error:", err)
		return
	}

	gxlog.CInfo("response result:%s", user)
}
