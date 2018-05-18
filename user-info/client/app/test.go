/******************************************************
# DESC    : test codec
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# LICENCE : Apache License 2.0
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2016-06-17 17:40
# FILE    : test.go
******************************************************/

package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
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
)

func testJsonrpc(userKey string) {
	var (
		err        error
		service    string
		method     string
		serviceIdx int
		user       *JsonRPCUser
		ctx        context.Context
		req        client.Request
		clt        client.Client
	)

	serviceIdx = -1
	service = "com.ikurento.user.UserProvider"
	for i := range conf.Service_List {
		if conf.Service_List[i].Service == service && conf.Service_List[i].Protocol == codec.CODECTYPE_JSONRPC.String() {
			serviceIdx = i
			break
		}
	}
	if serviceIdx == -1 {
		panic(fmt.Sprintf("can not find service in config service list:%#v", conf.Service_List))
	}

	// Create request
	method = string("GetUser")
	clt = rpcClient[codec.CODECTYPE_JSONRPC]
	// 注意这里，如果userKey是一个叫做UserKey类型的对象，则最后一个参数应该是 []UserKey{userKey}
	gxlog.CInfo("jsonrpc selected service %#v", conf.Service_List[serviceIdx])
	req = clt.NewRequest(
		conf.Service_List[serviceIdx].Group,
		conf.Service_List[serviceIdx].Version,
		service,
		method,
		[]string{userKey},
	)

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

// test int32 value
func testDubboCalc() {
	var (
		err        error
		service    string
		method     string
		serviceIdx int
		args       []interface{}
		req        client.Request
		rsp        int32
		clt        client.Client
	)

	serviceIdx = -1
	service = "com.ikurento.user.UserProvider"
	for i := range conf.Service_List {
		if conf.Service_List[i].Service == service && conf.Service_List[i].Protocol == codec.CODECTYPE_DUBBO.String() {
			serviceIdx = i
			break
		}
	}
	if serviceIdx == -1 {
		panic(fmt.Sprintf("can not find service in config service list:%+v", conf.Service_List))
	}

	// Create request
	method = string("Calc")
	args = []interface{}{int64(1), int64(2)}
	clt = rpcClient[codec.CODECTYPE_DUBBO]
	gxlog.CInfo("dubbo selected service %#v", conf.Service_List[serviceIdx])
	req = clt.NewRequest(
		conf.Service_List[serviceIdx].Group,
		conf.Service_List[serviceIdx].Version,
		service,
		method,
		args,
	)

	// Call service
	if err = clt.Call(context.Background(), req, &rsp, client.WithDialTimeout(requestTimeout)); err != nil {
		gxlog.CError("client.Call() return error:%+v", jerrors.ErrorStack(err))
		return
	}

	log.Info("response result:%s", rsp)
}

// test struct
func testDubboSum() {
	var (
		err        error
		service    string
		method     string
		serviceIdx int
		args       []interface{}
		req        client.Request
		rsp        Response
		clt        client.Client
	)

	serviceIdx = -1
	service = "com.ikurento.user.UserProvider"
	for i := range conf.Service_List {
		if conf.Service_List[i].Service == service && conf.Service_List[i].Protocol == codec.CODECTYPE_DUBBO.String() {
			serviceIdx = i
			break
		}
	}
	if serviceIdx == -1 {
		panic(fmt.Sprintf("can not find service in config service list:%+v", conf.Service_List))
	}

	// Create request
	method = string("Sum")
	args = []interface{}{int64(1), int64(2)}
	clt = rpcClient[codec.CODECTYPE_DUBBO]
	gxlog.CInfo("dubbo selected service %#v", conf.Service_List[serviceIdx])
	req = clt.NewRequest(
		conf.Service_List[serviceIdx].Group,
		conf.Service_List[serviceIdx].Version,
		service,
		method,
		args,
	)

	// Call service
	if err = clt.Call(context.Background(), req, &rsp, client.WithDialTimeout(requestTimeout)); err != nil {
		gxlog.CError("client.Call() return error:%+v", jerrors.ErrorStack(err))
		return
	}

	log.Info("response result:%s", rsp)
}

// test list
func testDubboGetUsers() {
	var (
		err        error
		service    string
		method     string
		serviceIdx int
		args       []interface{}
		req        client.Request
		rsp        []DubboUser
		clt        client.Client
	)

	serviceIdx = -1
	service = "com.ikurento.user.UserProvider"
	for i := range conf.Service_List {
		if conf.Service_List[i].Service == service && conf.Service_List[i].Protocol == codec.CODECTYPE_DUBBO.String() {
			serviceIdx = i
			break
		}
	}
	if serviceIdx == -1 {
		panic(fmt.Sprintf("can not find service in config service list:%+v", conf.Service_List))
	}

	// Create request
	method = string("GetUsers")
	args = []interface{}{[]string{"001", "003", "004"}}
	clt = rpcClient[codec.CODECTYPE_DUBBO]
	gxlog.CInfo("dubbo selected service %#v", conf.Service_List[serviceIdx])
	req = clt.NewRequest(
		conf.Service_List[serviceIdx].Group,
		conf.Service_List[serviceIdx].Version,
		service,
		method,
		args,
	)

	// Call service
	if err = clt.Call(context.Background(), req, &rsp, client.WithDialTimeout(requestTimeout)); err != nil {
		gxlog.CError("client.Call() return error:%+v", jerrors.ErrorStack(err))
		return
	}

	log.Info("response result:%s", rsp)
}

// test map
func testDubboGetUserMap() {
	var (
		err        error
		service    string
		method     string
		serviceIdx int
		args       []interface{}
		req        client.Request
		rsp        map[string]DubboUser
		clt        client.Client
	)

	serviceIdx = -1
	service = "com.ikurento.user.UserProvider"
	for i := range conf.Service_List {
		if conf.Service_List[i].Service == service && conf.Service_List[i].Protocol == codec.CODECTYPE_DUBBO.String() {
			serviceIdx = i
			break
		}
	}
	if serviceIdx == -1 {
		panic(fmt.Sprintf("can not find service in config service list:%+v", conf.Service_List))
	}

	// Create request
	method = string("GetUserMap")
	args = []interface{}{[]string{"001", "003", "004"}}
	clt = rpcClient[codec.CODECTYPE_DUBBO]
	gxlog.CInfo("dubbo selected service %#v", conf.Service_List[serviceIdx])
	req = clt.NewRequest(
		conf.Service_List[serviceIdx].Group,
		conf.Service_List[serviceIdx].Version,
		service,
		method,
		args,
	)

	// Call service
	if err = clt.Call(context.Background(), req, &rsp, client.WithDialTimeout(requestTimeout)); err != nil {
		gxlog.CError("client.Call() return error:%+v", jerrors.ErrorStack(err))
		return
	}

	log.Info("response result:%s", rsp)
}
