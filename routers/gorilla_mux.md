- [gorilla/mux](#gorillamux)
- [mux.Router](#muxrouter)
  - [Custom Handler](#custom-handler)
- [Matching Routes](#matching-routes)
  - [Request Parameters](#request-parameters)
  - [HTTP Methods](#http-methods)
  - [Path Prefix](#path-prefix)
  - [Matching Domain](#matching-domain)
  - [URL Schemes](#url-schemes)
  - [Header Values](#header-values)
  - [Query Values](#query-values)
  - [Custom Matcher Function](#custom-matcher-function)
  - [Subrouting](#subrouting)
  - [Registered URLs](#registered-urls)
- [Middleware](#middleware)
- [Static Files](#static-files)
- [Testing Handlers](#testing-handlers)

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

# mux.Router

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

# Matching Routes

## Request Parameters

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

## HTTP Methods

`gorilla/mux` 支持通過 `Methods` 方法來限定 request methods:

```go
r := mux.NewRouter()
r.HandleFunc("/hello/{name:[a-z]+}", sayHelloWorld).Methods("GET", "POST")
r.Handle("/zh/hello/{name}", &HelloWorldHandler{}).Methods("GET")
log.Fatal(http.ListenAndServe(":8080", r))
```

使用 `curl` 測試, 對 `http://localhost:8080/zh/hello/golang` 發起 `POST` 請求結果為空, 表示不支持此 request method

## Path Prefix

`gorilla/mux` routing 也支持 matching path prefix:

```go
r.PathPrefix("/hello").HandlerFunc(sayHelloWorld)
```

> 通常 matching path prefix 不會單獨使用, 會與 subrouters 結合使用, 從而實現對 router 分組

## Matching Domain

`gorilla/mux` 還支持 matching domain, 只需在原來的 routing rules 追加上 `Host` 方法調用並指定 domain 即可:

```go
r.HandleFunc("/hello/{name:[a-z]+}", sayHelloWorld).Methods("GET").Host("goweb.test")
```

如此一來只有當 request URL domain 為 `goweb.test` 時才會匹配到對應 routing rules

## URL Schemes

`gorilla/mux` router 支持通過 `Schemes` 方法設定 matching URL shemes:

```go
r.Handle("/zh/hello/{name}", &HelloWorldHandler{}).Methods("GET").Host("zh.goweb.test").Schemes("https")
```

如此一來只有 `HTTPS` request 才能訪問對應 routing rules, 對於 `HTTP` request 則會返回 `404` status code

## Header Values

可以在 `gorilla/mux` 路由定義中通過 `Headers` 方法設置 request header matching

下面範例中 request header 必須包含 `X-Requested-With` 且值為 `XMLHttpRequest` 才可以訪問指定路由 `/request/header`:

```go
r.HandleFunc("/request/header", func(w http.ResponseWriter, r *http.Request) {
    header := "X-Requested-With"
    fmt.Fprintf(w, "Including request header[%s=%s]", header, r.Header[header])
}).Headers("X-Requested-With", "XMLHttpRequest")
```

這樣做的意義在於限制 client 只能通過 Ajax request 訪問此路由

## Query Values

除了 request header 之外, 還可以通過 `Queries` 方法限定 query values

下面範例中 query values 必須包含 `token` 且值為 `test` 才能訪問指定路由 `/query/string`:

```go
r.HandleFunc("/query/string", func(w http.ResponseWriter, r *http.Request) {
    query := "token"
    fmt.Fprintf(w, "Including query value[%s=%s]", query, r.FormValue(query))
}).Queries("token", "test")
```

## Custom Matcher Function

`gorilla/mux` router 支持通過 `MatcherFunc` 方法自定義 routing matching rules, 在該方法中可以獲取到 request instance `request`, 這樣就可以取得所有的 user request info 並對其進行判斷, 符合預期的 requests 才能匹配並訪問路由

下面範例限定只有來自 `https://regy.dev` domain 的 requests 才可以匹配到 `/custom/matcher` 路由:

```go
r.HandleFunc("/custom/matcher", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Requests from specific domain: %s", r.Referer())
}).MatcherFunc(func(request *http.Request, match *mux.RouteMatch) bool {
    return request.Referer() == "https://regy.dev"
})
```

## Subrouting

`gorilla/mux` 支持路由分組和命名, 以及根據命名路由生成對應 URL

`gorilla/mux` 基於 subrouter 來實現路由分組功能, 下面範例以文章 CRUD 為例, 將文章相關 routing rules 劃分到 routing prefix 為 `/posts` 的 subrouter 中:

```go
func listPosts(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "List posts")
}

func createPost(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Create post")
}

func updatePost(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Update post")
}

func deletePost(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Delete post")
}

func showPost(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Show post")
}

...

// group routes（base on subrouter + path prefix）
postRouter := r.PathPrefix("/posts").Subrouter()
postRouter.HandleFunc("/", listPosts).Methods("GET")
postRouter.HandleFunc("/create", createPost).Methods("POST")
postRouter.HandleFunc("/update", updatePost).Methods("PUT")
postRouter.HandleFunc("/delete", deletePost).Methods("DELETE")
postRouter.HandleFunc("/show", showPost).Methods("GET")
```

如此一來 `/posts` prefix 會印用到後面所有基於 `postRouter` subrouter 的 routing rules 上, 且針對不同操作還限制了對應的 request methods, 測試上述路由訪問:

```go
curl http://localhost:8080/posts/posts/
curl http://localhost:8080/posts/posts/create -X POST
curl http://localhost:8080/posts/posts/update -X PUT
curl http://localhost:8080/posts/posts/delete -X DELETE
curl http://localhost:8080/posts/posts/show -X GET
```

若上述 router 是後台管理路由, 還可以結合 subrouter 近一步劃分:

```go
postRouter := r.PathPrefix("/posts").Host("admin.goweb.test").Subrouter()
```

如此一來只有 domain 為 `admin.goweb.test` 時才可以訪問對應路由, 提高了安全性

## Registered URLs

`gorilla/mux` 支援路由命名, 通過 `Name` 方法在 routing rules 中指定:

```go
postRouter := r.PathPrefix("/posts").Subrouter()
postRouter.HandleFunc("/", listPosts).Methods("GET").Name("posts.index")
postRouter.HandleFunc("/create", createPost).Methods("POST").Name("posts.create")
postRouter.HandleFunc("/update", updatePost).Methods("PUT").Name("posts.update")
postRouter.HandleFunc("/delete", deletePost).Methods("DELETE").Name("posts.delete")
postRouter.HandleFunc("/show/{id:[0-9]+}", showPost).Methods("GET").Name("posts.show")
```

可以像下面這樣根據上述路由命名生成與之對應的 URL:

```go
// print matching routes URL
indexUrl, _ := r.Get("posts.index").URL()
log.Println("posts list link: ", indexUrl)

createUrl, _ := r.Get("posts.create").URL()
log.Println("create post link: ", createUrl)

showUrl, _ := r.Get("posts.show").URL("id", "1")
log.Println("show post list: ", showUrl)
```

# Middleware

![middleware](img/middleware.png)

Middleware 典型的應用場景包含 authentication, logging, request header operations 和 `ResponseWriter Hijack`

一個經典的 Mux routing middleware 通常通過一個 clousure 來定義, 可以在 clousure function 中處理傳入的 request 和 response instances 或增加額外的業務邏輯, 再調用傳入的 handler 繼續後續 request handle

比如可以這樣定義一個 logging middleware:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        log.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}
```

要將上述 mux logging middleware apply 到所有 routes 上可以通過 `Use` 方法:

```go
r := mux.NewRouter()
r.Use(loggingMiddleware)
```

也可以將其 apply 到 subrouter 以限制其 scope:

```go
postRouter := r.PathPrefix("/posts").Subrouter()
postRouter.Use(loggingMiddleware)
```

下面範例實現 mux 版本的 token check middleware:

```go
func checkToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        token := r.FormValue("token")
        if token == "regy.dev" {
            log.Printf("Token check success: %s\n", r.RequestURI)
            // Call the next handler, which can be another middleware in the chain, or the final handler.
            next.ServeHTTP(w, r)
        } else {
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
    })
}

// apply middleware to subrouter
postRouter := r.PathPrefix("/posts").Subrouter()
postRouter.Use(checkToken)
```

這樣只有傳遞了正確的 `token` 參數才可以正常訪問路由

# Static Files

HTTP Server 除了處理動態資源外, 也有處理靜態資源的能力, 如 HTML, CSS, javascript, 圖片等

處理靜態資源需要借助 `PathPrefix()` 方法指定靜態資源所在 path prefix, 然後在 request handler 中通過 `http.FileServer` 直接返回文件內容本身作為響應:

```go
func main()  {
    r := mux.NewRouter()
    r.Use(loggingMiddleware)

    // parsing server start parameters dir as static resources web root path
    // default .
    var dir string
    flag.StringVar(&dir, "dir", ".", "static resources path")
    flag.Parse()

    // handle http://localhost:8000/static/<filename> static routes
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))    
    
    // other routes
    ...
    
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

當請求 `http://localhost:8080/status/app.js` 文件時, 會到 `static` path 下尋找 `app.js`, 若找不到則會返回 404, 否則返回 file 作為響應

# Testing Handlers

對應用來說, health check 無非是檢查應用本身是否可用, 及其依賴的核心服務是否可用, 這些核心服務通常包括 DB, Cache 等

下面範例為最簡化版本的 health check api, 只檢查了應用本身是否可用, 判斷方式是其是否 正常訪問並 return 200 status code

```go
// server.go

package main

import (
    "github.com/gorilla/mux"
    "io"
    "log"
    "net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    // Furthermore, we can check DB, Cache status via PING command and return the status within response 
    io.WriteString(w, `{"alive": true}`)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/health", HealthCheckHandler)
    log.Fatal(http.ListenAndServe("localhost:8080", r))
}
```

可以通過 `curl -v http://localhost:8080/health` 測試 api `/health` 是否可用:

![health_check_api](../network/img/health_check_api.png)

除了通過 `curl` 對 HTTP api 進行測試, 也可以編寫測試程式碼對 HTTP api 進行測試, 這裡使用 Go `httptest` package 來編寫

`httptest` 可用於模擬 web server, 來測試 `net/http` package 發送的 HTTP request 和捕獲 HTTP response 的方法, 要編寫一個 HTTP 測試步驟如下:
- create HTTP Multiplexer
- apply testing handler methods to multiplexer and test
- base on `net/http` methods create a `Request` instance to simulate client requests(include request URL and parameters)
- base on `net/http` methods create a `ResponseRecorder` instance and pass to multiplexer `ServeHTTP` method to get responses
- Finally get response status code and instance from `ResponseRecorder` and judge if test case pass

再來按照上述流程編寫 HTTP test, HTTP test 和 Unit Test 約定規則一樣, 因此需要在 `server.go` 同層目錄下創建一個測試文件 `server_test.go` 並編寫測試程式碼:

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthCheckHandler(t *testing.T) {
    // 初始化 router 並添加被測試的 handler methods
    mux := http.NewServeMux()
    mux.HandleFunc("/health", HealthCheckHandler)

    // create 一個 request instance 模擬 client requests, 其中包含 request methods, URL, parameters 等
    req, err := http.NewRequest("GET", "/health", nil)
    if err != nil {
        t.Fatal(err)
    }

    // 創建 ResponseRecorder 捕捉響應
    rr := httptest.NewRecorder()

    // 傳入測試請求和響應實體並執行請求
    mux.ServeHTTP(rr, req)

    // 檢查 status code（通過 ResponseRecorder 獲得）是否為 200, 若不是則測試不通過
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // 检查响应实体是否符合预期结果，如果不是，则测试不通过
    expected := `{"alive": true}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
```

![health_check_test](../network/img/health_check_test.png)