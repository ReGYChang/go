- [RPC](#rpc)
  - [RPC in Go](#rpc-in-go)
  - [encoding/gob](#encodinggob)
    - [Gob En/Deconding Rules](#gob-endeconding-rules)
    - [Gob En/Deconding in Go](#gob-endeconding-in-go)

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

## encoding/gob

Gob 是 Go 中一個序列化資料結構的編解碼工具, 當一個資料結構使用 Gob 進行序列化後即可於網絡中傳輸, 其經典適用場景就是 RPC 傳輸

`net/rpc` package 默認使用 `encoding/gob` 進行編解碼, 以 `rpc.Client` 為例, 其初始化程式碼如下:

```go
func NewClient(conn io.ReadWriteCloser) *Client {
  encBuf := bufio.NewWriter(conn)
  client := &gobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf}
  return NewClientWithCodec(client)
}
```

發送端會在發送訊息之前使用 `gob.Encoder` 對資料進行編碼, 接收端在收到訊息後會通過 `gob.Decoder` 對資料進行解碼

> Summary

與 `JSON` 或 `XML` 這種基於文本描述的資料交換格式不同, `Gob` 是 bianry encoding data stream, 因此性能和傳輸效率更高, 且 `Gob` Stream 是可以自解釋的, 具備完整的表達能力

但作為針對 Go 語言的資料結構編解碼專用序列化工具, `Gob` 無法跨語言使用, 只能侷限在 Go 開發的 RPC client 與 server 間通訊

當需要與其他程式語言實現 client 或 server 通訊就需要對 `net/rpc` package 底層的編解碼工具自定義, 改用 `JSON` 或 `Protobuf` 來進行資料格式序列化

### Gob En/Deconding Rules

對 Gob 而言, 發送方和接收方的資料結構並不需要完全一致, 以官方示例為例:

![gob_rules](img/gob_rules.png)

`struct { A, B int }` 結構編碼的資料可以被後面 9 種結構類型接受解碼, 接收資料結構只要滿足與發送資料結構簽名一致, 或者為發送資料類型的子集(但不能為空), 即可正常接收並解碼

具體不同資料類型規則如下:
- `struct`, `array`, `slice` 可以被編碼, 但 `function` 和 `channel` 無法
- `int` 分為有號數和無號數, 其之間無法互相編解碼
- `bool` 被當作 `uint` 來編碼, `0` 為 `false`, `1` 為 `true`
- `float`都被當作 `float64` 類型來編碼, `float` 和 `int` 也無法互相編解碼
- `string` (包含 `string` 和 `[]byte`) 是以無號數 byte 個數 + 每個 byte 編碼形式編解碼
- `array` (包含 `slice` 和 `array`) 是按照無號數個數 + 每個 array element 編碼的形式進行編解碼
- `map` 按照無號數元素個數 + key value pair 形式進行編解碼
- `struct` 按照序列化屬性名稱 + 屬性值來進行編解碼, 若屬性值為 0 或 nil 則直接被忽略

> 最後需注意 `struct` 類型中的屬性名稱都應該以大寫字母開頭, 確保為 public

### Gob En/Deconding in Go

```go
package main

import (
  "bytes"
  "encoding/gob"
  "fmt"
  "log"
)

type P struct {
  X, Y, Z int
  Name    string
  Tags    []string
  Attr    map[string]string
}

type Q struct {
  X, Y *int32
  Name string
  Tags    []string
  Attr    map[string]string
}

func main() {
  var network bytes.Buffer
  enc := gob.NewEncoder(&network)  // init encoder gob.Encoder
  dec := gob.NewDecoder(&network)  // init decoder gob.Decoder
  // data encode（before send the data）
  err := enc.Encode(P{3, 4, 5, "regy", []string{"Java", "MongoDB", "Go"}, map[string]string{"webiste": "https://regy.dev"}})
  if err != nil {
    log.Fatal("encode error:", err)
  }
  // data decode（after receive the data）
  var q Q
  err = dec.Decode(&q)
  if err != nil {
    log.Fatal("decode error:", err)
  }
  fmt.Printf("%q: {%d,%d}, Tags: %v, Attr: %v\n", q.Name, *q.X, *q.Y, q.Tags, q.Attr)
}
```