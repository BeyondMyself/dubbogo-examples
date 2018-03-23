/******************************************************
# DESC    : env var & configure
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# LICENCE : LGPL V3
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2016-07-21 16:41
# FILE    : config.go
******************************************************/

package main

import (
	"fmt"
	"os"
	"path"
)

import (
	"github.com/AlexStocks/goext/log"
	log "github.com/AlexStocks/log4go"
	config "github.com/koding/multiconfig"
)

import (
	"github.com/AlexStocks/dubbogo/codec"
	"github.com/AlexStocks/dubbogo/codec/jsonrpc"
	"github.com/AlexStocks/dubbogo/common"
	"github.com/AlexStocks/dubbogo/registry"
	"github.com/AlexStocks/dubbogo/registry/zk"
	"github.com/AlexStocks/dubbogo/transport"
)

const (
	APP_CONF_FILE     string = "APP_CONF_FILE"
	APP_LOG_CONF_FILE string = "APP_LOG_CONF_FILE"
)

type (
	RegistryNew  func(...registry.Option) registry.Registry
	TransportNew func(...transport.Option) transport.Transport
)

var (
	conf *ServerConfig

	DefaultRegistries = map[string]RegistryNew{
		"zookeeper": zookeeper.NewProviderZookeeperRegistry,
	}

	DefaultTransports = map[string]TransportNew{
		"http": transport.NewHTTPTransport,
	}

	// protocol:contentType
	DefaultContentTypes = map[string]string{
		"jsonrpc": "application/json",
	}

	// protocol:codecs
	DefaultCodecs = map[string]codec.NewCodec{
		"jsonrpc": jsonrpc.NewCodec,
	}
)

type (
	ServerConfig struct {
		// pprof
		Pprof_Enabled bool `default:"false"`
		Pprof_Port    int  `default:"10086"`

		// transport & registry
		Transport string `default:"http"`
		Registry  string `default:"zookeeper"`
		// application
		Application_Config common.ApplicationConfig
		// Registry_Address  string `default:"192.168.35.3:2181"`
		Registry_Config registry.RegistryConfig
		Service_List    []registry.ServiceConfig
		Server_List     []registry.ServerConfig
	}
)

func initServerConf() *ServerConfig {
	var (
		confFile string
	)

	confFile = os.Getenv(APP_CONF_FILE)
	if confFile == "" {
		panic(fmt.Sprintf("application configure file name is nil"))
		return nil
	}
	if path.Ext(confFile) != ".toml" {
		panic(fmt.Sprintf("application configure file name{%v} suffix must be .toml", confFile))
		return nil
	}

	conf = new(ServerConfig)
	config.MustLoadWithPath(confFile, conf)
	gxlog.CInfo("config{%#v}\n", conf)

	return conf
}

func configInit() error {
	var (
		confFile string
	)

	initServerConf()

	confFile = os.Getenv(APP_LOG_CONF_FILE)
	if confFile == "" {
		panic(fmt.Sprintf("log configure file name is nil"))
		return nil
	}
	if path.Ext(confFile) != ".xml" {
		panic(fmt.Sprintf("log configure file name{%v} suffix must be .xml", confFile))
		return nil
	}

	log.LoadConfiguration(confFile)

	return nil
}
