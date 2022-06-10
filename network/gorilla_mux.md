- [gorilla/mux](#gorillamux)
  - [mux.Router](#muxrouter)
  - [Routing parameters](#routing-parameters)
  - [Custom Handler](#custom-handler)

# gorilla/mux

Go 官方 standard libray `net/http` 通過自帶的 `DefaultServeMux` 提供的 routing handler 雖然簡單, 但存在許多不足:
- 不支持參數設定, 如 `/user/:uid` 這種泛型匹配
- 對 restful api 支持不友善, 無法限制訪問路由的方法
- 對於擁有很多 routing rules 的應用, 編寫大量路由規則非常繁瑣

`gorilla/mux` 提供了更強大的 routing handler, 與 `http.ServeMux` 實現原理相同, `gorilla/mux` 提供的 router implementation type `mux.Router` 也會匹配 user requests 與系統註冊的 routing rules, 並將請求轉發

`mux.Router` 主要有以下特性:
- 實現 `http.Handler` interface, 所以與 `http.ServeMux` 完全兼容
- 可以基於 URL host, path, prefix, scheme, request header, request parameters, request method 進行 routing
- URL host, path, query string 支持可選正則匹配
- 支持構建或反轉已註冊的 URL host, 以便維護對資源的引用
- 支持路由嵌套, 以便不同路由可以共享通用條件, 比如 host, path prefix 等

安裝 `gorilla/mux`:

```go
go get -u github.com/gorilla/mux
```

## mux.Router

```go
package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

func sayHelloWorld(w http.ResponseWriter, r *http.Request)  {
    w.WriteHeader(http.StatusOK)  // set HTTP status code 200
    fmt.Fprintf(w, "Hello, World!")  // send response to client
}

func main()  {
    r := mux.NewRouter()
    r.HandleFunc("/hello", sayHelloWorld)
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

在 `main` 函數中第一行顯式初始化了 `mux.Router` 作為 router, 並在這個 router 中註冊 routing rules, 最後將此 router 傳入 `http.ListenAndServe` 函數

在 broswer 訪問 `http://localhost:8080/hello` 即可渲染出以下結果:

```go
Hello, World!
```

## Routing parameters

如果想在 router 中設置 routing parameters, 例如 `/hello/world`, `/hello/regy`, 可以通過以下方法實現:

```go
r.HandleFunc("/hello/{name}", sayHelloWorld)
```

或是可以通過正則表達式限制參數字符串:

```go
r.HandleFunc("/hello/{name:[a-z]+}", sayHelloWorld)
```

以上 routing parameters 僅支持小寫字母, 不支持其他字符

相應地在 closure 處理函數中, 需要這樣解析 routing parameters:

```go
func sayHelloWorld(w http.ResponseWriter, r *http.Request)  {
    params := mux.Vars(r)
    w.WriteHeader(http.StatusOK)  // set HTTP status code to 200
    fmt.Fprintf(w, "Hello, %s!", params["name"])  // send response to client
}
```

可以通過 `http://localhost:8080/hello/regy` 這種方式請求路由:

```go
Hello, Regy!
```

>❗️若請求參數中包含中文則會直接返回 404, 表示路由匹配失敗

## Custom Handler

與 `http.ServeMux` 相同, 在 `mux.Router` 也可以將請求轉發給自定義 handler 類型:

```go
package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

func sayHelloWorld(w http.ResponseWriter, r *http.Request)  {
    params := mux.Vars(r)
    w.WriteHeader(http.StatusOK)  // set HTTP status code to 200
    fmt.Fprintf(w, "Hello, %s!", params["name"])  // send response to client
}

type HelloWorldHandler struct {}

func (handler *HelloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
    params := mux.Vars(r)
    w.WriteHeader(http.StatusOK)  // set HTTP status code to 200
    fmt.Fprintf(w, "你好, %s!", params["name"])  // send response to client
}

func main()  {
    r := mux.NewRouter()
    r.HandleFunc("/hello/{name:[a-z]+}", sayHelloWorld)
    r.Handle("/zh/hello/{name}", &HelloWorldHandler{})
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

這裡的 `HelloWorldHandler` 也要實現 `Handler` interface 定義的 `ServeHTTP` 方法, 調用方式與之前相同, 只需要通過 `r.Handle` 方法傳入 Handler instance 即可