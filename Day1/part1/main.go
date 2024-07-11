package main

import (
	"fmt"
	"log"
	"net/http"
	// "os"
)

// 我们设置了2个路由，/和/hello，分别绑定 indexHandler 和 helloHandler ，
// 根据不同的HTTP请求会调用不同的处理函数。访问/，响应是URL.Path = /，而/hello的响应则是请求头(header)中的键值对信息。
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	// http.ListenAndServe(":8080", nil)
	log.Fatal("ListenAndServe: ", http.ListenAndServe(":999", nil))
}

/*
indexHandler 处理主页请求

	该函数旨在响应浏览器请求主页时的行为。它通过读取请求的URL路径，
	并将该路径以格式化字符串的形式回显到浏览器，帮助开发者调试或确认请求的URL是否符合预期。

参数:

	w http.ResponseWriter - 用于向客户端发送响应的接口
	r *http.Request - 包含客户端请求信息的结构体指针

返回值:无
*/
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

/*
helloHandler 是一个HTTP处理函数，用于响应HTTP请求。
它的目的是将请求头的信息写入响应体。

参数:
  w http.ResponseWriter: 用于向客户端发送HTTP响应的接口。
  r *http.Request: 代表客户端发起的HTTP请求的对象。

该函数不返回任何值。
*/

func helloHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
