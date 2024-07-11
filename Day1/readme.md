# 设计一个框架
* 大部分时候，我们需要实现一个 Web 应用，第一反应是应该使用哪个框架。不同的框架设计理念和提供的功能有很大的差别。
* 那为什么不直接使用标准库，而必须使用框架呢？在设计一个框架之前，我们需要回答框架核心为我们解决了什么问题。只有理解了这一点，才能想明白我们需要在框架中实现什么功能。

## 先看看标准库net/http如何处理一个请求
```go
func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/count", counter)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}
```

* net/http提供了基础的Web功能，即监听端口，映射静态路由，解析HTTP报文。
* 一些Web开发中简单的需求并不支持，需要手工实现。
    * 动态路由：例如hello/:name，hello/*这类的规则。
    * 鉴权：没有分组/统一鉴权的能力，需要在每个路由映射的handler中实现。
    * 模板：没有统一简化的HTML机制。

来源：https://geektutu.com/post/gee.html

# 关于VSCode安装Go tools失败的问题
* 设置代理：
```go
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.io,direct
```
* 设置完成后重启VS Code，按照提示安装即可。

# 第一部分纯享代码
```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal("ListenAndServe: ", http.ListenAndServe(":999", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

```
## 解释
我们设置了2个路由，/和/hello，分别绑定 indexHandler 和 helloHandler ， 根据不同的HTTP请求会调用不同的处理函数。访问/，响应是URL.Path = /，而/hello的响应则是请求头(header)中的键值对信息。
<hr>
main 函数的最后一行，是用来启动 Web 服务的，第一个参数是地址，:999表示在 999 端口监听。而第二个参数则代表处理所有的HTTP请求的实例，nil 代表使用标准库中的实例处理。第二个参数，则是我们基于net/http标准库实现Web框架的入口。

## 测试
```shell
    curl http://localhost:999/
    curl http://localhost:999/hello
```

# 第二部分：实现http.Handler接口
* 第二个参数的类型是什么呢？通过查看net/http的源码可以发现，Handler是一个接口，需要实现方法 ServeHTTP ，也就是说，只要传入任何实现了 ServerHTTP 接口的实例，所有的HTTP请求，就都交给了该实例处理了。

## 第二部分纯享代码
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	// "os"
)

type Engine struct {
}

func main() {
	enginge := new(Engine)
	log.Fatal(http.ListenAndServe(":8080", enginge))
}

func (enginge *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
```
## 解释
* 我们定义了一个空的结构体Engine，实现了方法ServeHTTP。这个方法有2个参数，第二个参数是 Request ，该对象包含了该HTTP请求的所有的信息，比如请求地址、Header和Body等信息；第一个参数是 ResponseWriter ，利用 ResponseWriter 可以构造针对该请求的响应。

* 在 main 函数中，我们给 ListenAndServe 方法的第二个参数传入了刚才创建的engine实例。至此，我们走出了实现Web框架的第一步，即，将所有的HTTP请求转向了我们自己的处理逻辑。还记得吗，在实现Engine之前，我们调用 http.HandleFunc 实现了路由和Handler的映射，也就是只能针对具体的路由写处理逻辑。比如/hello。但是在实现Engine之后，我们拦截了所有的HTTP请求，拥有了统一的控制入口。在这里我们可以自由定义路由映射的规则，也可以统一添加一些处理逻辑，例如日志、异常处理等。

* 代码的运行结果与之前的是一致的。

## 测试
```shell
    curl http://localhost:999/
    curl http://localhost:999/hello
```

# Gee框架的雏形
```go
gee/
  |--gee.go
  |--go.mod
main.go
go.mod
```
## go.mod
```go
module gee-web

go 1.22.5
```

# 第三部分：搭建出整个框架的雏形
## 第三部分纯享代码
### gee/gee.go
```go
package gee

import (
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (enginge *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	enginge.router[key] = handler
}

// Get
func (enginge *Engine) GET(pattern string, handler HandlerFunc) {
	enginge.addRoute("GET", pattern, handler)
}

// Post
func (enginge *Engine) Post(pattern string, handler HandlerFunc) {
	enginge.addRoute("POST", pattern, handler)
}

// Run
func (enginge *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, enginge)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 NOT FOUND"))
	}
}
```
### gee/go.mod
```go
module gee

go 1.22.5

```

### ./go.mod
```go
module gee-web

go 1.22.5

require gee v0.0.0

replace gee => ./gee
```
### ./main.go
```go
package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	r.Run(":999")
}

```

## 详解
### main.go
* 使用New()创建 gee 的实例，使用 GET()方法添加路由，最后使用Run()启动Web服务。这里的路由，只是静态路由，不支持/hello/:name这样的动态路由，动态路由我们将在下一次实现。
### gee.go
* 首先定义了类型HandlerFunc，这是提供给框架用户的，用来定义路由映射的处理方法。我们在Engine中，添加了一张路由映射表router，key 由请求方法和静态路由地址构成，例如GET-/、GET-/hello、POST-/hello，这样针对相同的路由，如果请求方法不同,可以映射不同的处理方法(Handler)，value 是用户映射的处理方法。

* 当用户调用(*Engine).GET()方法时，会将路由和处理方法注册到映射表 router 中，(*Engine).Run()方法，是 ListenAndServe 的包装。

* Engine实现的 ServeHTTP 方法的作用就是，解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND 。

## 测试
```shell
    curl http://localhost:999/
    curl http://localhost:999/hello
```

# 总结
* 整个Gee框架的原型已经出来了。实现了路由映射表，提供了用户注册静态路由的方法，包装了启动服务的函数。当然，到目前为止，我们还没有实现比net/http标准库更强大的能力，不用担心，很快就可以将动态路由、中间件等功能添加上去了。

![a585907ff8c372aaab7e0311af9f658821372460](http://qny.expressisland.cn/gdou24/a585907ff8c372aaab7e0311af9f658821372460.png)