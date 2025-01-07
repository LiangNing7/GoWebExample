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
# 然后在浏览器中访问：localhost:8080/static/images/Hello.jpg
```

# Routing

> 该文件目录为`gowebexample02/03-Routing`

Go 的 `net/http` 包为 HTTP 协议提供了许多功能。 它做得不太好的一件事是复杂的请求路由，比如将请求 URL 分段为单个参数。 幸运的是，有一个非常流行的包可以做到这一点，它以 Go 社区中良好的代码质量而闻名。 在此示例中，您将看到如何使用 `gorilla/mux` 包创建具有命名参数、GET/POST 处理程序和域限制的路由。

## `Installing the gorilla/mux`：安装`gorilla/mux`软件包

`gorilla/mux` 是一个适用于 Go 默认 HTTP 路由器的包。它附带了许多功能，可提高编写 Web 应用程序时的生产力。 它还符合 Go 的默认请求处理程序签名 `func (w http.ResponseWriter, r *http.Request)`，因此该包可以与其他 HTTP 库（如中间件或现有应用程序）混合和匹配。使用 `go get` 命令从 GitHub 安装包，如下所示：

```bash
$ go get -u github.com/gorilla/mux
```

## `Creating a new Router`：创建新`Router`

首先创建一个新的请求路由器。路由器是 Web 应用程序的主路由器，稍后将作为参数传递给服务器。 它将接收所有 HTTP 连接，并将其传递给您将在其上注册的请求处理程序。 您可以像这样创建一个新路由器：

```go
r := mux.NewRouter()
```

## `Registering a Request Handler`：注册请求处理程序

一旦你有了新的路由器，你可以像往常一样注册请求处理程序。 唯一的区别是，你不再调用 `http.HandleFunc(...)`，而是像这样在你的路由器上调用 HandleFunc：`r.HandleFunc(...)`。

## `URL Parameters`：URL 参数

`gorilla/mux` 路由器最大的优势在于能够从请求 URL 中提取片段。 例如，这是您应用程序中的一个 URL：

```go
/books/go-programming-blueprint/page/10
```

此 URL 具有两个动态片段：

1. 图书标题：`go-programming-blueprint`
2. 页面：`10`

要让请求处理程序匹配上面提到的 URL，请用 URL 模式中的占位符替换动态段，如下所示：

```go
r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
	// get the book
    // navigate to the page
})
```

最后一步是从这些段中获取数据。 该包带有一个函数 `mux.Vars(r)`，它将 `http.Request` 作为参数并返回一个段映射。

```go
func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    vars["title"] // the book title slug
    vars["page"] // the page
}
```

## `Setting the HTTP server's router`：设置 HTTP 服务器的路由器

有没有想过 `http.ListenAndServe(":80", nil)` 中的 `nil` 是什么意思？它是 HTTP 服务器主路由器的参数。 默认情况下，它是 `nil`，这意味着使用 `net/http` 包的默认路由器。要使用自己的路由器，请用路由器 `r` 的变量替换 `nil`。

```go
http.ListenAndServer(":80",r)
```

示例代码如下：

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})
	http.ListenAndServe(":8080", r)
}
```

可以使用如下代码运行：

```bash
$ go run .\03-Routing\Routing.go
# 运行代码后，打开浏览器，访问：localhost:8080:/books/go-programming-blueprint/page/10
# 在网页中显示：You've requested the book: go-programming-blueprint on page 10
```

