- [Create Gee Web Framework](#create-gee-web-framework)
  - [main.go](#maingo)
  - [gee.go](#geego)
- [Context Design](#context-design)
  - [context.go](#contextgo)
  - [router.go](#routergo)
- [gee.go](#geego-1)

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

# Context Design

對於 Web Server 來說, 無非是根據請求 `*http.Request` 來構建響應 `http.ResponseWriter`

但這兩個物件提供的 interface 粒度太細, 如若要構建一個完整的響應需要考慮 `Header` 和 `Body`, 而 `Header` 包含了 `StatusCode`, `ContentType` 等幾乎每次請求都需要設置的資訊, 若不有效封裝則需要大量復寫繁瑣的程式碼且容易出錯

以返回 JSON 資料作比較範例:

封裝前

```go
obj = map[string]interface{}{
    "name": "geektutu",
    "password": "1234",
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
encoder := json.NewEncoder(w)
if err := encoder.Encode(obj); err != nil {
    http.Error(w, err.Error(), 500)
}
```

封裝後

```go
c.JSON(http.StatusOK, gee.H{
    "username": c.PostForm("username"),
    "password": c.PostForm("password"),
})
```

針對使用場景封裝 `*http.Request` 和 `http.ResponseWriter` 的方法並簡化相關方法調用, 只是設計 `Context` 的原因之一

對於 framework 來說還需要支持額外的功能, 如將來解析動態路由 `/hello/:name` 時參數 `:name` 值擺放的位置? 或是 framework 需要支持 middleware, 而 middleware 所產生的訊息要放在哪?

Context 隨著每個請求的產生而建立, 隨請求結束而銷毀, 和當前請求強相關的資訊都應由 Context 乘載

因此設計 `Context` struct, 擴展性和複雜性保留在內部, 對外則簡化了接口, 路由的 handle function 和要實現的 middleware, 參數都統一使用 `Context` instance

這裡做了以下改進:
- 將 router 獨立, 方便後續新增功能
- 設計 `Context`, 封裝 Request 和 Response, 提供對 JSON, HTML 等返回型別的支持

main.go
```go
func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
```

- Handler 的參數變成 `gee.Context`, 提供了查詢 `Query/PostForm` 參數的功能
- `gee.Context` 封裝了 HTML/String/JSON 函數, 能夠快速構建 HTTP response

## context.go

```go
import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
```

- 幫 `map[string]interface{}` 起了一個別名 `gee.H`, 在構建 JSON 資料時更簡潔
- `Context` 目前只包含了 `http.ResponseWriter` 和 `*http.Request`, 另外還提供了 `Method` 和 `Path` 這兩個常用屬性的直接訪問
- 提供了訪問 `Query` 和 `PostForm` 參數的方法
- 提供了快速構建 String/Data/JSON/HTML 響應的方法

## router.go

將和 router 相關的方法和 struct 拆出來, 方便之後對 router 功能的增強, 如提供動態路由的支持

router 的 `handle` 方法做了一個細微的調整, 即 `handle` 的參數變成了 `Context`

```go
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
```

# gee.go

```go
// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
}

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
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
	c := newContext(w, req)
	engine.router.handle(c)
}
```

在調用 `router.handle` 之前構建了 `Context` instance, 目前 `Context` 僅封裝了原本的兩個參數
