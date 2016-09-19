# dubbogo #
a golang micro-service framework compatible with alibaba dubbo. just using jsonrpc 2.0 protocol over http now.

## 说明 ##
---
> dubbogo 目前版本(0.1.1)支持的codec 是jsonrpc 2.0，transport protocol是http。
> 只要你的java程序支持jsonrpc 2.0 over http，那么dubbogo程序就能调用它。
> dubbogo自己的server端也已经实现，即dubbogo既能调用java service也能调用dubbogo实现的service。
> 由于dubbogo还没有上传到github，使用的时候请先下载https://github.com/AlexStocks/dubbogo，然后放在路径$/gopath}/github.com/AlexStocks/下面。

## dubbogo examples ##
---
*dubbogo examples是基于dubbogo的实现的代码示例，目前提供echo和user-info两个例子*

> dubbogo-examples借鉴java的编译思路，提供了区别于一般的go程序的而类似于java的独特的编译脚本系统。

### dubogo example1: user-info ###
---
*从这个程序可以看出dubbogo程序能够调用dubbo的服务*

> 1 部署zookeeper服务；
>
> 2 请部署 https://github.com/QianmiOpen/dubbo-rpc-jsonrpc 服务端，如果你不想编译，可以使用我编译好的 dubbogo-examples/user-info/java-server/dubbo_jsonrpc_example.bz2，注意修改zk地址；
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
