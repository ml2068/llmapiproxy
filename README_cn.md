# llmapiproxy 

通过 VPS 代理，从本地轻松调用 LLM API

HTTP 反向代理服务器 

## 概述

这是一个用 Go 编写的简单 HTTP 反向代理服务器。它可以将接收到的 HTTP 请求转发到指定的目标服务器，并添加授权头（api-key）到请求中。

## 特点 

*   **环境变量支持**：服务器可以使用 github.com/joho/godotenv 包从 .env 文件加载环境变量。 

*   **后台运行模式**：服务器可以以守护进程模式运行，这允许它在后台运行并且不会在终端关闭时终止。 

*   **HTTP 反向代理**：服务器可以将接收到的 HTTP 请求转发到指定的目标服务器，并添加授权头（api-key）到请求中。 

*   **响应头日志记录**：服务器可以记录目标服务器的响应头。 

## 使用方法 

要运行程序，请设置环境变量 PORT、TARGET 和 APIKEY，然后启动程序。

检查守护进程和 api.log 文件。

通过杀死进程（kill PID）停止程序。

`go run apiproxy.go -daemon`

`ps -ef | grep go`

`kill 732924`

## 配置 

服务器可以通过环境变量进行配置。以下环境变量受支持：

* PORT：监听的端口号。 

* TARGET：目标网址(API endpoint ) LLM 服务器的 URL。 

* API_KEY：用于 LLM 授权的 API 密钥。 

## 许可证 

[Apache-2.0 许可证](./LICENSE)
