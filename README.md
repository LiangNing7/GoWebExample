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

# 03-Routing

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

## `Features of the gorilla/mux Router`： `gorilla/mux`路由器的特性

### `Methods`：方法

将请求处理程序限制为特定的 HTTP 方法。

- ```go
  r.HandleFunc("/books/{title}", CreateBook).Methods("POST")
  r.HandleFunc("/books/{title}", ReadBook).Methods("GET")
  r.HandleFunc("/books/{title}", UpdateBook).Methods("PUT")
  r.HandleFunc("/books/{title}", DeleteBook).Methods("DELETE")
  ```

### `Hostnames & Subdomains`：主机名和子域

将请求处理程序限制到特定主机名或子域名。

- ```go
  r.HandleFunc("/books/{title}", BookHandler).Host("www.mybookstore.com")
  ```

### `Schemes`：方案

将请求处理程序限制为 `http/https`

- ```go
  r.HandleFunc("/secure", SecureHandler).Schemes("https")
  r.HandleFunc("/insecure", InsecureHandler).Schemes("http")
  ```

### `Path prefixes & Subrouters`：路径前缀和子路由器

将请求处理程序限制为特定路径前缀。

- ```go
  bookrouter := r.PathPrefix("/books").Subrouter()
  bookrouter.HandleFunc("/", AllBooks)
  bookrouter.HandleFunc("/{title}", GetBook)
  ```

# 04-MySQLDatabase

> 该文件目录为`gowebexample02/04-MySQLDatabase`

在某个时间点，您希望您的 Web 应用程序存储和检索来自数据库的数据。 当您处理动态内容、为用户提供输入数据的表单或存储登录和密码凭据以对用户进行身份验证时，几乎总是这种情况。为此，我们有数据库。

数据库的种类和形式多种多样。整个网络上普遍使用的一种数据库是 MySQL 数据库。它已经存在了很长时间，并且已经无数次证明了自己的地位和稳定性。

在这个示例中，我们将深入了解 Go 中数据库访问的基础知识，创建数据库表，存储数据并再次检索数据。

## `Installing the go-sql-driver/mysql`：安装`go-sql-driver/mysql`包

Go 编程语言附带一个名为 `database/sql` 的便捷包，用于查询各种 SQL 数据库。这非常有用，因为它将所有常见的 SQL 特性抽象为一个供你使用的 API。Go 不包含的是数据库驱动程序。在 Go 中，数据库驱动程序是一个包，它实现了特定数据库（在本例中为 MySQL）的底层详细信息。正如你可能已经猜到的，这对于保持向前兼容性非常有用。由于在创建所有 Go 包时，作者无法预见到未来出现的每个数据库，并且支持所有可能的数据库将是一项巨大的维护工作。

要安装 MySQL 数据库驱动程序，请转到您选择的终端并运行：

```bash
go get -u github.com/go-sql-driver/mysql
```

## `Connecting to a MySQL database`：连接到`MySQL`数据库

在安装所有必需的软件包后，我们需要检查的第一件事是，我们是否可以成功连接到我们的 MySQL 数据库。要检查我们是否可以连接到我们的数据库，导入database/sqlandgo-sql-driver/mysql包并打开一个连接，如下所示：

```go
import "database/sql"
import _ "go-sql-driver/mysql"


// Configure the database connection (always check errors)
db, err := sql.Open("mysql", "username:password@(127.0.0.1:3306)/dbname?parseTime=true") // 这里的 username，password，dbname，请输入自己MySQL对应的username与password与dbname
//在安装完MySQL后，需要登录到mysql，然后创建对应的数据库，例如 create database goweb;


// Initialize the first connection to the database, to see if everything works correctly.
// Make sure to check the error.
err := db.Ping()
```

## `Creating our first database table`：创建我们的第一个数据库表

我们数据库中的每个数据条目都存储在一个特定的表中。数据库表由列和行组成。列为每个数据条目提供一个标签并指定其类型。行是插入的数据值。在我们的第一个示例中，我们希望创建一个这样的表：

| id   | username | password | created_at          |
| ---- | -------- | -------- | ------------------- |
| 1    | johndoe  | secret   | 2019-08-10 12:30:00 |

翻译成 SQL，创建表的命令将如下所示：

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME,
    PRIMARY KEY (id)
);
```

现在我们有了 SQL 命令，我们可以使用 `database/sql` 包在 MySQL 数据库中创建表：

```go
query := `
    CREATE TABLE users (
        id INT AUTO_INCREMENT,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME,
        PRIMARY KEY (id)
    );`

// Executes the SQL query in our database. Check err to ensure there was no error.
_, err := db.Exec(query)
```

## `Inserting our first user`：插入我们的第一个用户

如果您熟悉 SQL，那么向我们的表中插入新数据就像创建我们的表一样容易。需要注意的一件事是：默认情况下，Go 使用预处理语句将动态数据插入到我们的 SQL 查询中，这是一种将用户提供的数据安全地传递到我们的数据库而不存在任何损坏风险的方法。在 Web 编程的早期，程序员直接使用查询将数据传递到数据库，这导致了大规模的漏洞，并可能破坏整个 Web 应用程序。请不要这样做。正确操作很容易。

要将我们的第一个用户插入到我们的数据库表中，我们创建一个类似于以下内容的 SQL 查询。 正如你所看到的，我们省略了 id 列，因为 MySQL 会自动设置它。问号告诉 SQL 驱动程序，它们是实际数据的占位符。 这是你可以看到我们讨论过的准备好的语句的地方。

```sql
INSERT INTO users (username, password,created_at) VALUES (?, ?, ?)
```

我们现在可以在 Go 中使用此 SQL 查询，并将新行插入到表中：

```go
import "time"

username := "johndoe"
password := "secret"
createdAt := time.Now()

// Inserts our data into the users table and returns with the result and a possible error.
// The result contains information about the last inserted id (which was auto-generated for us) and the count of rows this query affected.
result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
```

要获取为用户新创建的 `id`，只需像这样获取即可：

```go
userID, err := result.LastInsertId()
```

## `Querying our users table`：查询我们的 `users`表

现在我们有一个用户在我们的表中，我们想要查询它并获取它的所有信息。在 Go 中，我们有两种可能性来查询我们的表。 有 `db.Query` 可以查询多行，以便我们进行迭代，并且在只想要查询特定行的情况下有 `db.QueryRow`。

查询特定行基本上与我们之前介绍过的其他 SQL 命令类似。我们的 SQL 命令通过其 ID 查询单个用户，如下所示：

```sql
SELECT id, username, password, created_at FROM users WHERE id = ?
```

在 Go 中，我们首先声明一些变量来存储我们的数据，然后查询单个数据库行，如下所示：

```go
var (
    id        int
    username  string
    password  string
    createdAt time.Time
)

// Query the database and scan the values into out variables. Don't forget to check for errors.
query := `SELECT id, username, password, created_at FROM users WHERE id = ?`
err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt)
```

## `Querying all users`：查询所有用户

在前面的部分中，我们介绍了如何查询单个用户行。许多应用程序都有用例，您希望查询所有现有用户。 这与上面的示例类似，但涉及更多编码。

我们可以使用上面示例中的 SQL 命令并去掉 WHERE 子句。这样，我们就可以查询所有现有用户。

```sql
SELECT id, username, password, created_at FROM users
```

在 Go 中，我们首先声明一些变量来存储我们的数据，然后查询单个数据库行，如下所示：

```go
type user struct {
    id        int
    username  string
    password  string
    createdAt time.Time
}

rows, err := db.Query(`SELECT id, username, password, created_at FROM users`) // check err
defer rows.Close()

var users []user
for rows.Next() {
    var u user
    err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt) // check err
    users = append(users, u)
}
err := rows.Err() // check err
```

用户切片现在可能包含类似这样的内容：

```json
users {
    user {
        id:        1,
        username:  "johndoe",
        password:  "secret",
        createdAt: time.Time{wall: 0x0, ext: 63701044325, loc: (*time.Location)(nil)},
    },
    user {
        id:        2,
        username:  "alice",
        password:  "bob",
        createdAt: time.Time{wall: 0x0, ext: 63701044622, loc: (*time.Location)(nil)},
    },
}
```

## `Deleting a user from our table`：从我们的表中删除用户

最后，从我们的表中删除用户与上面部分的 `.Exec` 一样简单：

```go
_, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1) // check err
```

示例代码如下：

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 要保证有 goweb 数据库，登录到mysql后，使用 create database goweb; 进行创建。
	db, err := sql.Open("mysql", "root:liangningMySql123@(127.0.0.1:3306)/goweb?parseTime=true")
	if err != nil {
		log.Fatal("1", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("2", err)
	}

	{ // Create a new table
		query := `
			CREATE TABLE users(
			    id INT AUTO_INCREMENT,
			    username TEXT NOT NULL,
			    password TEXT NOT NULL,
			    created_at DATETIME,
			    PRIMARY KEY (id)
			);`
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}

	{ // Insert a new user
		username := "johndoe"
		password := "123456"
		createdAt := time.Now()

		result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		id, err := result.LastInsertId()
		fmt.Println(id)
	}

	{ // Query a single user
		var (
			id        int
			username  string
			password  string
			createdAt time.Time
		)

		query := "SELECT id,username,password,created_at FROM users WHERE id = ?"
		if err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, username, password, createdAt)
	}

	{ // Query all users
		type user struct {
			id         int
			username   string
			password   string
			created_at time.Time
		}

		rows, err := db.Query(`SELECT id,username,password,created_at FROM users`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var users []user
		for rows.Next() {
			var u user
			err := rows.Scan(&u.id, &u.username, &u.password, &u.created_at)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v\n", users)
	}

	{
		_, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1)
		if err != nil {
			log.Fatal(err)
		}
	}
}
```

可以使用如下代码运行：

```bash
$ go run .\04-MySQLDatabase\mySQLDatabase.go
1
1 johndoe 123456 2025-01-07 07:14:06 +0000 UTC
[]main.user{main.user{id:1, username:"johndoe", password:"123456", created_at:time.Time{wall:0x0, ext:63871830846, loc:(*time.Location)(nil)}}}
```

最后在`goweb`数据库下，有`users`的空表。

# 05-Templates

> 该文件目录为`gowebexample02/05-Templates`

Go 的 `html/template` 包为 HTML 模板提供了一种丰富的模板语言。 它主要用于 Web 应用程序，以结构化的方式在客户端浏览器中显示数据。 Go 模板语言的一大好处是自动转义数据。 无需担心 XSS 攻击，因为 Go 会解析 HTML 模板并在将其显示到浏览器之前转义所有输入。

## `First Template`：第一个模板

用 Go 编写模板非常简单。此示例显示了一个待办事项列表，该列表以 HTML 中的无序列表 (`ul`) 形式编写。 在呈现模板时，传入的数据可以是 Go 的任何类型的数据结构。它可以是一个简单的字符串或数字，甚至可以是嵌套数据结构，如下面的示例所示。要访问模板中的数据，最顶层的变量通过 `{{.}}` 访问。大括号内的点称为管道，是数据的根元素。

```go
data := ToDoPageData{
    PageTitle: "My TODO list",
    Todos: []Todo{
        {Title: "Task 1", Done: false},
        {Title: "Task 2", Done: true},
        {Title: "Task 3", Done: true},
    }
}
```

```html
<h1>{{.PageTitle}}</h1>
<ul>
    {{range .Todos}}
        {{if .Done}}
            <li class="done">{{.Title}}</li>
        {{else}}
            <li>{{.Title}}</li>
        {{end}}
    {{end}}
</ul>
```

## `Control Structures`：控制结构

模板语言包含一组丰富的控制结构来呈现您的 HTML。在这里，您将获得最常用结构的概述。 要获取所有可能结构的详细列表，请访问：`text/template`

| Control Structure                | Definition                                   |
| -------------------------------- | -------------------------------------------- |
| `{{/* a comment */}}`            | 定义注释                                     |
| `{{.}}`                          | 渲染根元素                                   |
| `{{.Title}}`                     | 在嵌套元素中呈现`Title`字段                  |
| `{{if .Done}} {{else}} {{end}}`  | 定义 if 语句                                 |
| `{{range .Todos}} {{.}} {{end}}` | 遍历所有`Todos`并使用 `{{.}}` 渲染每个`Todo` |
| `{{block "content" .}} {{end}}`  | 定义一个名为`content`的块                    |

## `Parsing Templates from Files`：从文件中解析模板

模板可以从字符串或磁盘上的文件进行解析。 由于通常情况下模板是从磁盘解析的，因此本示例演示如何执行此操作。 在此示例中，Go 程序所在目录中有一个名为 `layout.html` 的模板文件。

```go
tmpl,err := template.ParseFiles("layout.html")
// or
tmpl := template.Must(template.ParseFiles("layout.html"))
```

## `Execute a Template in a Request Handler`：在请求处理程序中执行模板

一旦从磁盘解析模板，就可以在请求处理程序中使用它。 `Execute` 函数接受一个 `io.Writer` 用于写出模板和一个 `interface{}` 用于将数据传递到模板中。 当该函数在 `http.ResponseWriter` 上调用时，`Content-Type` 标头会自动设置在 HTTP 响应中，为 `Content-Type: text/html; charset=utf-8`。

```go
func(w http.ResponseWriter, r *http.Request) {
    tmpl.Execute(w, "data goes here")
}
```

示例代码如下：

```go
package main

import (
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	tmpl := template.Must(template.ParseFiles("templates/layout.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":8080", nil)
}
```

创建`templates`文件夹，在其中创建`layout.html`，在其中写入：

```html
<h1>{{.PageTitle}}</h1>
<ul>
    {{range .Todos}}
        {{if .Done}}
            <li class="done">{{.Title}}</li>
        {{else}}
            <li>{{.Title}}</li>
        {{end}}
    {{end}}
</ul>
```

可以使用如下代码运行：

```bash
$ go run .\05-Templates\templats.go
# 打开浏览器，访问 localhost:8080/ 会出现如下界面
```

![image-20250107154043910](https://gitee.com/liangningi/typora_picture/raw/master/Go/202501071541218.png)

# 06-AssetsAndFiles

> 该文件目录为`gowebexample02/06-AssetsAndFiles`

此示例将展示如何从特定目录提供静态文件，如 CSS、JavaScript 或 images。

示例代码如下：

```go
package main

import "net/http"

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
```

在`gowebexample02`文件夹中创建`assets/css`文件夹，并在器中创建`styles.css`文件

其文件中写入：

```css
body {
    background-color: black;
}
```

可以使用如下代码运行：

```bash
$ go run .\06-AssetsAndFiles\assetsAndFiles.go

# 使用浏览器访问 localhost:8080/static/css/styles.css
```

