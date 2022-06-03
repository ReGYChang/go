- [RPC](#rpc)
  - [RPC in Go](#rpc-in-go)

# RPC

RPC(Remote Procedure Call) 是一種通過網路請求從遠端 server 調用 service, 而不需要了解底層網絡細節的通訊協議

RPC 協議基於傳輸層的 TCP 或 UDP, 或是由應用層的 HTTP 構建, 允許開發者直接調用另一台 server 上的程式而無需為調用過程額外編寫通訊相關程式碼

> RPC 使得開發分散式應用程式更加容易, 現在流行的微服務通常基於 RPC 協議

相較於 HTTP 採用 browser - server (B/S) model, RPC 採用 client - server (C/S) model, 請求程式為 client, 遠端 server 提供程式是 server

執行一個 RPC 調用時, client 首先會發送一個帶有參數的請求到 server, 並等待 server response;

在 server 端 service process 保持監聽狀態, 當 client request 到達時, server 通過 parsing 請求參數計算出結果, 並向 client 發送 response, 接著繼續等待下一個 client request

client 收到 server repsonse 後可以執行相對應的業務邏輯, 也可以繼續進行其他的 RPC 調用

## RPC in Go

Go 標準庫提供 `net/rpc` package, 其實現了 RPC 協議的相關細節

`net/rpc` package 允許 RPC client 通過網絡或其他 I/O 連接調用一個 server 物件的 public methods

在 RPC server 需要將此物件註冊為可訪問的 service, 之後該物件的 public methods 就能通過 RPC 方式提供調用

一個 RPC server 可以註冊多個不同類型的物件, 但無法註冊同個類型的多個物件, 此外一個物件只有滿足以下條件的 method 才能被 RPC server 設置為可提供遠程調用:
- 必須為 public method
- 必須有兩個參數, 且參數類型必須是 package 外部可訪問的類型或是 Go build-in
- 第二個參數必須是一個 pointer
- 方法必須返回一個 `error` 類型的值

總結以上四點:
```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

其中類型 `T`, `T1` 和 `T2` 分別對應 service 物件所屬類型, 請求類型及回應類型, 默認都會使用 Go build-in `encoding/gob` package 進行編解碼

此方法 `MethodName` 第一個參數表示由 RPC client 傳入的請求參數, 第二個參數表示要返回給 RPC client 的 response, 最後返回一個 `error` 類型的值表示錯誤資訊