# 01-HelloWorld

> 该文件目录为`gowebexample02/01-HelloWorld`

Go 是一门内置了 web 服务器的编程语言。来自标准库的 `net/http` 包，其包含有关 HTTP 协议的所有功能。这包括一个 HTTP 客户端和一个 HTTP 服务器。

## `Registering a Request Handler`：注册请求处理程序。

首先，创建一个处理程序，它接收来自浏览器、HTTP 客户端或 API 请求的所有传入 HTTP 连接。 Go 中的处理程序是一个具有以下签名的函数：

```go
func (w http.ResponseWriter, r *http.Request)
```

该函数接收两个参数：

1. `http.ResponseWriter`：可以在其中编写`text/html`响应。
2. `http.Request`：包含有关此 HTTP 请求的所有信息【包括 URL 或 header 字段】

然后向默认 HTTP 服务器注册请求处理程序：

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
})
```

## `Listen for HTTP Connections` ：侦听 HTTP 连接

`request handler`，本身无法接收来自外部的任何 HTTP 连接。HTTP 服务器必须监听端口才能将连接传递给 `request handler`，由于在大多数情况下端口 80 是 HTTP 流量的默认端口，一般会被占用，所以选用`:8080`端口，因此此服务器也将监听该端口。以下代码将启动 Go 的默认 HTTP 服务器，并在端口 8080 上侦听连接。 您可以使用浏览器导航到 `http://localhost/:8080`，并查看服务器处理您的请求。

```go
http.ListenAndServer(":80",nil)
```

示例代码如下：

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested:%s\n", r.URL.Path)
	})
	http.ListenAndServe(":8080", nil)
}
```

可以使用如下代码运行：

```bash
$ go run .\01-HelloWorld\helloWorld.go

# 当代码运行起来后，点开浏览器，访问 localhost:8080/
# 在浏览器中，会出现如下字段：
Hello, you've requested:/

# 最后，ctrl + c 停止代码运行，即关闭服务器。
```

