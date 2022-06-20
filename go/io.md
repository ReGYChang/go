- [I/O](#io)
  - [Reader Interface](#reader-interface)
  - [Writer Interface](#writer-interface)
  - [Types Implement io.Reader and io.Writer](#types-implement-ioreader-and-iowriter)
  - [ReaderAt 和 WriterAt interface](#readerat-和-writerat-interface)

# I/O

`io` package 為 I/O 功能提供了基本的接口, 由於這些接口封裝的 I/O 操作由不同的低級操作實現, 因此在另外聲明之前不應該假設其併發執行是安全的

`io` package 中最重要的兩個 interface: `Reader` 和 `Writer`, 只要實現這兩個 interface 就可以使用 `io` package 的功能

## Reader Interface

Reader interface 定義如下:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

官方文檔關於此 intercace methods 的說明:

`Read` 將 len(p) 個 bytes 讀取到 p 中, 並 return 讀取的 bytes 數 n(0 <= n <= len(p)) 及 error

即便 `Read` return 的 n < len(p), 它也會在調用過程中佔用 len(p) 個 bytes 作為暫存空間, 並返回可用資料, 而不是等待更多資料

當 `Read` 成功讀取 n > 0 個 bytes 後遇到一個錯誤或 `EOF(end-of-life)`, 其會 return 讀取的 bytes 數, 並可能會同時在本次調用中 return 一個 non-nil error, 或在下一次調用 return 這個 error(且 n 為 0)

一般情況下 `Reader` 會 return 一個非 0 bytes 數 n, 若 n = len(p) 個 bytes 從 input source 結尾處由 `Read` 返回, `Read` 可能返回 `err == EOF` 或 `err == nil`, 且之後的 `Read()` 都應該返回 `(n:0, err:EOF)`

調用者在考慮錯誤之前應先處理返回的資料, 這樣做可以正確地處理在讀取一些 bytes 後產生的 I/O 錯誤, 同時允許 EOF 出現

`Reader` interface 只包含一個 `Read` 方法, 只要實現了 `Read` 方法的物件都滿足 `io.Reader` interface

下面來看一下此 interface 的用法:

```go
func ReadFrom(reader io.Reader, num int) ([]byte, error) {
    p := make([]byte, num)
    n, err := reader.Read(p)
    if n > 0 {
        return p[:n], nil
    }
    return p, err
}
```

`ReadFrom` 函數將 `io.Reader` 作為參數, 即 `ReadFrom` 可以從任意地方讀取資料, 只要來源實現 `io.Reader` interface, 如可以從標準輸入, 文件, 字符串等儲取資料:

```go
// 從標準輸入讀取
data, err = ReadFrom(os.Stdin, 11)

// 從普通文件讀取, 其中 file 是 os.File 的 instance
data, err = ReadFrom(file, 9)

// 從字符串讀取
data, err = ReadFrom(strings.NewReader("from string"), 12)
```

>💡TIP:

`io.EOF` 變數的定義為 `var EOF = errors.New("EOF")`, 其為 `error` 型別, 根據 reader interface 說明, 在 n > 0 且資料被讀取完的情況下, 返回的 `error` 可能是 `EOF` 也有可能是 nil

## Writer Interface

`Writer` interface 定義如下:

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

官方文檔對於該 interface 方法說明:

`Write` 將 len(p) 個 bytes 從 p 中寫入到基本資料流中, 返回從 p 中被寫入的 bytes 數量 n(0 <= n <= len(p)) 及任何遇到引起寫入提前結束的 error

若 `Write` 返回的 n < len(p), 它就必須返回一個 non-nil error

與 `Reader` 相同, 所有實現 `Write` 方法的型別都實現了 `Writer` interface

這裡通過標準庫的例子來了解 `Writer` 用法:

在 `fmt` package 有一組函數 `Fprint/Fprintf/Frpintln`, 它們接收一個 `io.Writer` 型別參數, 它們接收一個 `io.Writer` 型別的參數(第一個參數), 也就是其將資料格式化輸出到 `io.Writer` 中

以 `fmt.Fprintln` 為例, 並同時看一下 `fmt.Println` 函數 source code:

```go
func Println(a ...interface{}) (n int, err error) {
    return Fprintln(os.Stdout, a...)
}
```

顯然 `fmt.Println` 會將內容輸出到標準中

## Types Implement io.Reader and io.Writer

通過上述例子可以發現 `os.File` 同時實現了這兩個 interface, 還可以看到 `os.Stdin/Stdout` 這樣的程式碼, 其在 `os` package 中聲明如下:

```go
var (
    Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
    Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
    Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)
```

也就是說 `Stdin/Stdout/Stderr` 只是三個特殊的文件型別的標示(os.File 的 instance), 自然也實現了 `io.Reader` 和 `io.Writer`

通過查看標準庫文件列出實現了 `io.Reader` 和 `io.Writer` interface 的型別:
- `os.File` 同時實現了 `io.Reader` 和 `io.Writer`
- `strings.Reader` 實現了 `io.Reader`
- `bufio.Reader/Writer` 分别實現了 `io.Reader` 和 `io.Writer`
- `bytes.Buffer` 同時實現了 `io.Reader` 和 `io.Writer`
- `bytes.Reader` 實現了 `io.Reader`
- `compress/gzip.Reader/Writer` 分别實現了 `io.Reader` 和 `io.Writer`
- `crypto/cipher.StreamReader/StreamWriter` 分别實現了 `io.Reader` 和 `io.Writer`
- `crypto/tls.Conn` 同时實現了 `io.Reader` 和 `io.Writer`
- `encoding/csv.Reader/Writer` 分别實現了 `io.Reader` 和 `io.Writer`
- `mime/multipart.Part` 實現了 `io.Reader`
- `net/conn` 分别實現了 `io.Reader` 和 `io.Writer`(Conn interface 定義了 Read/Write)

除此之外 `io` package 本身也有這兩個 interface 的實現型別, 如:
- Implement `Reader`: `LimitedReader`, `PipeReader`, `SectionReader`
- Implement `Writer`: `PipeWriter`

以上型別較常使用的有: `os.File`, `strings.Reader`, `bufio.Reader/Writer`, `bytes.Buffer`, `bytes.Reader`

>💡TIP:

從 interface 命名可以觀察到, 在 Go 中 interface 的命名約定是以 `er` 結尾, 這裡並非強制要求, 標準庫中有些 interface 也不是以 `er` 結尾

## ReaderAt 和 WriterAt interface

`ReaderAt` interface 定義如下:

```go
type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}
```

