- [Create Gee Web Framework](#create-gee-web-framework)
  - [main.go](#maingo)
  - [gee.go](#geego)

# Create Gee Web Framework

這裡模仿 `gin` framework 設計並編寫 `gee` framwork 的 design 及 API

首先創建新文件夾, 初始化 module 並引用 local path packages:

```shell
go mod init gee
go mod edit -replace gee=/Users/regy/Github/omni-chat/gee
```

## main.go

```go
package main

import (
	"fmt"
	"net/http"

	"gee"
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

	r.Run(":9999")
}
```

參考 gin framework 的 design, 使用 `New()` 創建 gee instance, 並使用 `GET()` 方法新增路由, 最後使用 `Run()` 啟動 web server

這裡的路由是靜態路由, 還不支持類似 `/hello/:name` 這樣的動態路由

## gee.go

```go
package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
```

- 定義了 `HandlerFunc` 型別提供給 framework user 使用, 其用來定義 route mapping handle function
- 在 `Engine` struct 中添加一個 route mapping map routers, key 由 request method 和 static routing address 組成, 如 `GET-/`, `GET-/hello`, 針對相同的路由若請求方法不同, 可以映射到不同的 handler 處理
- 當 user 調用 `(*Engine).GET()` 方法時, 會將路由和 handler registe 到 engine 的 routes mapping map
- `(*Engine).Run()` 方法是 ListenAndServe warpped 方法
- Engine 實現 `ServeHTTP` 方法的作用為解析 request path 並查找 routes mapping map, 若找到執行對應 handler, 否則 return `404 NOT FOUND`

