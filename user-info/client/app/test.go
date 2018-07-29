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
	jerrors "github.com/juju/errors"
)

import (
	"github.com/AlexStocks/dubbogo/client"
	"github.com/AlexStocks/dubbogo/registry"
)

func testJsonrpc(userKey string) {
	var (
		err        error
		service    string
		method     string
		serviceIdx int
		user       *JsonRPCUser
		ctx        context.Context
		conf       registry.ServiceConfig
		req        client.Request
		serviceURL *registry.ServiceURL
		clt        *client.HTTPClient
	)

	clt = client.NewHTTPClient(
		&client.HTTPOptions{
			HandshakeTimeout: clientConfig.connectTimeout,
			HTTPTimeout:      clientConfig.requestTimeout,
		},
	)

	serviceIdx = -1
	service = "com.ikurento.user.UserProvider"
	for i := range clientConfig.Service_List {
		if clientConfig.Service_List[i].Service == service && clientConfig.Service_List[i].Protocol == client.CODECTYPE_JSONRPC.String() {
			serviceIdx = i
			break
		}
	}
	if serviceIdx == -1 {
		panic(fmt.Sprintf("can not find service in config service list:%#v", clientConfig.Service_List))
	}

	// Create request
	method = string("GetUser")
	// 注意这里，如果userKey是一个叫做UserKey类型的对象，则最后一个参数应该是 []UserKey{userKey}
	gxlog.CInfo("jsonrpc selected service %#v", clientConfig.Service_List[serviceIdx])
	conf = registry.ServiceConfig{
		Group:    clientConfig.Service_List[serviceIdx].Group,
		Protocol: client.CodecType(client.CODECTYPE_JSONRPC).String(),
		Version:  clientConfig.Service_List[serviceIdx].Version,
		Service:  clientConfig.Service_List[serviceIdx].Service,
	}
	req = clt.NewRequest(conf, method, []string{userKey})

	serviceURL, err = clientRegistry.Filter(req.ServiceConfig(), 1)
	if err != nil {
		gxlog.CError("registry.Filter(conf:%#v) = error:%s", req.ServiceConfig(), jerrors.ErrorStack(err))
		return
	}
	// Set arbitrary headers in context
	ctx = context.WithValue(context.Background(), client.DUBBOGO_CTX_KEY, map[string]string{
		"X-Proxy-Id": "dubbogo",
		"X-Services": service,
		"X-Method":   method,
	})

	user = new(JsonRPCUser)
	// Call service
	if err = clt.Call(ctx, *serviceURL, req, user); err != nil {
		gxlog.CError("client.Call() return error:%+v", jerrors.ErrorStack(err))
		return
	}

	//log.Info("response result:%s", user)
	gxlog.CInfo("response result:%s", user)
}
