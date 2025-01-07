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

# 02-HTTPServer

> 该文件目录为`gowebexample02/02-HTTPServer`

接下来，我们将学习如何在 Go 中创建一个基本的 HTTP 服务器。首先我们讨论一下我们的 HTTP 服务器应该具备哪些功能。一个基本的 HTTP 服务器有几个关键任务需要处理。

- 处理`dynamic requests`：处理来自浏览网站，登录其账户或发布图片的用户发来的请求。
- 提供`static assets`：向浏览器提供`JavaScript,CSS,images`，为用户创造动态体验。
- 接收`connections`：HTTP 服务器必须监听特定端口才能接受来自互联网的连接。

## `Process dynamic requests`：处理动态请求

`net/http`包含接受请求并动态处理它们所需的所有实用程序。可以使用`http.HandlerFunc`函数注册一个新处理程序。它的第一个参数采用要匹配的路径，第二个参数采用要执行函数。在此示例中，当有人要浏览你的网站时，他会收到一条友好的消息。

```go
http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome to my website!")
})
```

对于动态方面，`http.Request`包含有关请求及其参数的所有信息。可以使用`r.URL.Query().Get("token")`读取 GET 参数，或使用`r.FormValue("email")`读取 POST 参数【HTML表单中的字段】

## `Serving static assets`：提供静态资产

要提供 `JavaScript,CSS,images`等静态资源，可以使用内置的`http.FileServer`并将其指向一个 URL 路径。为了让文件服务器正常工作，它需要知道从哪里提供文件，我们可以像这样操作：

```go
fs := http.FileServer(http.Dir("static/"))
```

一旦我们服务器就位，我们只需要在 URL 路径中指向它，就像我们对动态请求所做的那样。需要注意的一点是：为了正确提供文件，我们需要去掉 URL 路径的一部分，通常是我们文件所在的目录的名称。

```go
http.Handle("/static/",http.StripPrefix("/static/".fs))
```

## `Accept connections`接受连接

完成我们基本 HTTP 服务器的最后一件事是，监听端口以接受来自互联网的连接。 正如你所猜想的那样，Go 也有一个内置的 HTTP 服务器，我们可以很快地启动它。一旦启动，你就可以在浏览器中查看你的 HTTP 服务器。

```go
http.ListenAndServer(":8080",nil))
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
		fmt.Fprint(w, "Welcome to my website!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
```

在`gowebexample02`文件夹下，创建`static/images`文件夹。然后在里面放入一个`Hello.jpg`

可以使用如下代码运行：

```bash
$ go run .\02-HTTPServer\HTTPServer.go
# 然后在浏览器中访问：localhost:8080/static/images/Hello
```

