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

