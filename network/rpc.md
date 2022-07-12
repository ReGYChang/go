- [RPC](#rpc)
  - [RPC in Go](#rpc-in-go)
  - [encoding/gob](#encodinggob)
    - [Gob En/Deconding Rules](#gob-endeconding-rules)
    - [Gob En/Deconding in Go](#gob-endeconding-in-go)
- [gRPC](#grpc)
  - [Protobuf](#protobuf)
    - [Protobuf Setup](#protobuf-setup)
    - [Protobuf Demo](#protobuf-demo)

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

# gRPC

`gRPC` 是一個高性能且通用的 open source RPC framework, 其是由 Google 基於 mobile application development 開發並基於 `HTTP/2` 而設計, 支援 `ProtoBuf` 序列化標準及眾多開發語言

`HTTP/2` 在 `HTTP/1.1` 的基礎上做了大量優化, `HTTP/1.1` 雖然引入了 `Keep-Alive` 機制復用 TCP connection, 但還是有許多問題:
- 使用 `Keep-Alive` request 是 serializable(非 pipeline 時), 而 pipeline 時則有 `head-of-line blocking(HOL)` issue
- 每次都需要傳輸不必要的 Header
- 無法雙向通訊

> HTTP/1.1 允許多個 request 復用 connection, 可以同時將 request 全部發送出去, 不需要等一個 return 再發第二個, 以提升 concurrency, 而 server 需要將 response 按照 pipeline 中發送的順序進行順序返回, 如果前面的 request blocking, 後面的 request/response 則會被動等待

HTTP/2 則解決了這些問題, 並引入了新的機制:
- 在 client/server 建立了 header table, 每次只發送 index 以縮減 header 體積
- 建立 virtual channel, 將 data 拆分成多個 stream, 每個 steram 有自己的 id 及 piority, 且 stream 可以雙向傳輸, 每個 stream 可以拆成多個 frames, 可以將 request 切成多個 stream 發送, 每個 stream 獨立返回, 以避開 HTTP/1.1 serializable 或 `HOL` 問題

基於 HTTP/2 data stream 機制, gRPC client/server 可以實現批次優化, client 可以累積 request 一次性發送給 server, server 也可以批次返回結果, 以實現 stream rpc

與很多 RPC 系統類似, `gRPC` 也是基於以下理念: 定義一個 service, 指定其為能夠被遠端調用的方法(包含參數和返回型別), 在 server 實現 intereface, 並運行一個 gRPC server 來處理 client 調用, 在 client 擁有一個`存根`就像 server 一樣的方法

`gRPC` 默認使用 `protocol buffers`, 其為 Google 開源的一套成熟的結構化資料序列化標準

## Protobuf

`Protocol buffers` 是一種與程式語言, 平台無耦合的資料交換格式, 用於序列化結構化資料, 較 XML, JSON 而言, `Protobuf` 序列化後的 data stream 更小, 傳輸速度更高, 且操作更簡單

> Protocol buffers are a language-neutral, platform-neutral extensible mechanism for serializing structured data.

`protoc` 主要用於 compile `protobuf(.proto)` 檔案和 runtime, 其為 C++ 編寫, 以超高的壓縮率著稱, release 下載地址如下: 

[https://github.com/protocolbuffers/protobuf/releases
](https://github.com/protocolbuffers/protobuf/releases
)

需注意一點, `Protocol buffers` 可以獨立使用, 不一定要與 gRPC 綁定使用, 但若使用 gRPC 則一定要使用 `Protocol buffers`

使用 protobuf 作為序列化傳輸方案有以下幾個優點:
- 節省網路傳輸量, 傳輸速度更快且檔案大小更小
- 降低 CPU effort, parsing JSON 本身為 CPU intensive, 而 protobuf 本身為 binary format, 更接近底層資料表徵, 因此能有效降低 CPU effort
- 能夠根據不同程式語言 compile 出不同的檔案
- 可以邊寫註解, 型別顯式明確

透過 protobuf 定義好傳輸的資料欄位(message) 和呼叫的方法(service) 後, gRPC 即可在不同程式語言上運行

對於 JSON 等文本形式的序列化協定來說, protubuf 能達到幾十倍空間及性能的提升, 比如傳輸整數 123, 文本類協定需要 3 個 bytes(ascii 31 32 33) 來傳輸, 而 binary 類只需要一個 byte (01111011) 即可表達

同時 protobuf 會維護 `.proto` 檔案, 如此一來在 parsing 檔案生成 stub 程式時可以對 function name 進行編號, 傳輸時只需傳輸編號而不用傳 function name, 如此一來可以省下大量 bytes 傳輸量, 其他更多精巧的壓縮方式如 `TLV`, 可以參考 [proto encoding](https://link.zhihu.com/?target=https%3A//developers.google.com/protocol-buffers/docs/encoding)



### Protobuf Setup

Uncompress:

```shell
$ unzip protoc-3.14.0-linux-x86_64.zip -d protoc-3.14.0-linux-x86_64
```

Add env:

```shell
$ sudo vim /etc/profile 
```

Export env:

```shell
export PATH=$PATH:/home/regy/17x/protoc-3.14.0-linux-x86_64/bin
```

Active:

```shell
$ source /etc/profile
```

Check if install success:

```shell
$ protoc --version
libprotoc 3.14.0
```

除了安裝 `protoc` 之外還需要安裝各個程式語言對應的 compile plugin, 以下為 Go compile plugin:

```shell
go get google.golang.org/protobuf/cmd/protoc-gen-go
```

### Protobuf Demo

create `.proto` file: hello_world.proto

```protobuf
syntax = "proto3";  // 定義要使用的 protocol buffer 版本

package calculator;  // for name space
option go_package = "proto/calculator";  // generated code 的 full Go import path

message CalculatorRequest {
  int64 a = 1;
  int64 b = 2;
}

message CalculatorResponse {
  int64 result = 1;
}

service CalculatorService {
  rpc Sum(CalculatorRequest) returns (CalculatorResponse) {};
}
```

`protoc` compile:

```shell
# Syntax: protoc [OPTION] PROTO_FILES
$ protoc --proto_path=IMPORT_PATH  --go_out=OUT_DIR  --go_opt=paths=source_relative path/to/file.proto
```

- `-proto_path or -I`: 指定 `import` path, 可以指定多個參數, compile 時按順序查找, 不指定時默認查找當前目錄
  - `.proto` 也可以引入其他 `.proto` 檔案, 用於指定被 import 檔案位置
- `-go_out`: 指定輸出檔案路徑

```shell
$ protoc --go_out=. hello_word.proto
```

compile 結束後會產生一個 `hello_world.pb.go` 檔案, 即 compile 完成

> Compile Progress

- parsing `.proto` 檔案, compile 成 `protobuf` 原生資料結構並保存於 memory
- 將 `protobuf` 相關資料結構傳遞給對應程式語言的 compile plugin, 由 plugin 負責將接收到的 `protobuf` 原生資料結構渲染輸出為特定語言 template

後續提到的 `gRPC Plugins`, `gRPC-Gateway` 也是 `protoc` compile plugin, 將 `.proto` 檔案 compile 成對應組件所需要的原始檔