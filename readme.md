# dubbogo #
a golang micro-service framework compatible with alibaba dubbo. just using jsonrpc 2.0 protocol over http now.

## 说明 ##
---
- 1 dubbogo 目前版本(0.1.1)支持的codec 是jsonrpc 2.0，transport protocol是http。
- 2 只要你的java程序支持jsonrpc 2.0 over http，那么dubbogo程序就能调用它。
- 3 dubbogo自己的server端也已经实现，即dubbogo既能调用java service也能调用dubbogo实现的service。
- 4 使用的时候请先下载https://github.com/AlexStocks/dubbogo，然后放在路径$/gopath}/github.com/AlexStocks/下面。

## dubbogo examples ##
---
*dubbogo examples是基于dubbogo的实现的代码示例，目前提供echo和user-info两个例子*

> dubbogo-examples借鉴java的编译思路，提供了区别于一般的go程序的而类似于java的独特的编译脚本系统。

### dubogo example1: user-info ###
---
*从这个程序可以看出dubbogo程序(user-info/client & user-info/server) 如何与 java(dubbo)的程序(user-info/java-client & user-info/java-server)进行双向互调(测试的时候一定注意修改配置文件中服务端和zk的ip&port).*

> 1 部署zookeeper服务；
>
> 2 请编译并部署dubbogo-examples/user-info/java-server，注意修改zk地址(conf/dubbo.properties:line6:"dubbo.registry.address")和监听端口(conf/dubbo.properties:line6:"dubbo.protocol.port", 不建议修改port), 然后执行"bin/start.sh"启动java服务端；
>
> 3 修改dubbogo-examples/user-info/client/profiles/test/client.toml:line 33，写入正确的zk地址；
>
> 4 dubbogo-examples/user-info/client/下执行 sh assembly/windows/test.sh命令(linux下请执行sh assembly/linux/test.sh)，然后target/windows下即放置好了编译好的程序以及打包结果，在dubbogo-examples\user-info\client\target\windows\user_info_client-0.1.0-20160818-1346-test下执行sh bin/load.sh start命令即可客户端程序；
>
> 5 修改dubbogo-examples/user-info/server/profiles/test/server.toml:line 21，写入正确的zk地址；
>
> 6 dubbogo-examples/user-info/server/下执行 sh assembly/windows/test.sh命令(linux下请执行sh assembly/linux/test.sh)，然后target/windows下即放置好了编译好的程序以及打包结果，在dubbogo-examples\user-info\server\target\windows\user_info_server-0.1.0-xxxx下执行sh bin/load.sh start命令即可服务端程序；
>

### dubogo example2: echo ###
---

*这个程序是为了执行压力测试，整个编译部署过程可以参考user-info这个示例的相关操作步骤。*

### dubbogo example3: calculator ###
---
*这个程序主要是为了测试github.com/AlexStocks/codec/hessian/client_test.go,之所以取名calculator是因为服务端有math 8interface,无其他意义。*

测试的时候注意修改 github.com/AlexStocks/dubbogo-examples/calculator/java-server/src/com/ikurento/hessian:HessianServer.java中与url相关的设置。具体来说是，port(HessianServer.java:line 21) 和 path(HessianServer.java:line 40-42).

无论是java-server还是java-client，只需执行sh build-winidows.sh(linux: sh build-linux.sh)即可。
