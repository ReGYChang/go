- [I/O](#io)
  - [Reader Interface](#reader-interface)
  - [Writer Interface](#writer-interface)
  - [Types Implement io.Reader and io.Writer](#types-implement-ioreader-and-iowriter)
  - [ReaderAt 和 WriterAt interface](#readerat-和-writerat-interface)
  - [ReaderFrom & WriterTo interface](#readerfrom--writerto-interface)

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

> 官方文檔對於該 interface 方法說明:

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

>💡TIP:｀

從 interface 命名可以觀察到, 在 Go 中 interface 的命名約定是以 `er` 結尾, 這裡並非強制要求, 標準庫中有些 interface 也不是以 `er` 結尾

## ReaderAt 和 WriterAt interface

`ReaderAt` interface 定義如下:

```go
type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}
```

> 官方文件中關於該 interface 方法說明如下:

`ReadAt` 從 basic input source off set 開始, 將 `len(p)` bytes 讀取到 `p` 中, 並返回讀取的 bytes 數 `n(0 <= n <= len(p))` 及遇到的 error

當 `ReadAt` 返回的 `n < len(p)` 時, 其也會在調用過程中返回一個 non-nil error 來解釋為何沒有 return 更多的 bytes, 在這點上 `ReadAt` 相較 `Read` 更嚴謹

即使 `ReadAt` 返回的 `n < len(p)`, 它也會在調用過程中使用 `p` 的全部作為暫存空間; 若可讀取的數據不到 `len(p)` 字節, `ReadAt` 就會阻塞, 直到所有數據都可用或一個錯誤發生, 在這一點上 `ReadAt` 不同於 `Read`

若 `n = len(p)` 個字節從輸入源的結尾處由 `ReadAt` 返回, `Read` 可能返回 `err == EOF` 或者 `err == nil`

若 `ReadAt` 攜帶一個偏移量從輸入源讀取, `ReadAt` 應當既不影響偏移量也不被它所影響。

可對相同的輸入源並行執行 `ReadAt` 調用

由上可見 `ReaderAt` interface 可以從指定偏移量處開始讀取資料

簡單程式碼示範如下:

```go
reader := strings.NewReader("regy.dev")
p := make([]byte, 6)
n, err := reader.ReadAt(p, 2)
if err != nil {
    panic(err)
}
fmt.Printf("%s, %d\n", p, n)
```

output:

```go
gy.dev, 6
```

`WriterAt` interface 定義如下:

```go
type WriterAt interface {
    WriteAt(p []byte, off int64) (n int, err error)
}
```

官方文件對於 `WriterAt` interface 的說明:

`WriteAt` 從 `p` 中將 `len(p)` 個字節寫入到偏移量 off 處的基本數據流中, 它返回從 `p` 中被寫入的字節數 `n (0 <= n <= len(p))` 以及任何遇到的引起寫入提前停止的錯誤, 若 `WriteAt` 返回的 `n < len(p)`, 它就必須返回一個 non-nil 的錯誤。

若 `WriteAt` 攜帶一個偏移量寫入到目標中, `WriteAt` 應當既不影響偏移量也不被它所影響

若被寫區域沒有重疊, 可對相同的目標並行執行 `WriteAt` 調用

可以通過此 interface 將資料寫入到資料流的特定偏移量之後

通過程式碼範例演示 `WriteAt` 方法的使用(`os.File` 實現了 `WriterAt interface`):

```go
file, err := os.Create("writeAt.txt")
if err != nil {
    panic(err)
}
defer file.Close()
file.WriteString("regy.dev")
n, err := file.WriteAt([]byte("iro.meow"), 5)
if err != nil {
    panic(err)
}
fmt.Println(n)
```

## ReaderFrom & WriterTo interface

`ReaderFrom` 定義如下:

```go
type ReaderFrom interface {
    ReadFrom(r Reader) (n int64, err error)
}
```

> 官方文件關於此 interface methods 說明:

`ReadFrom` 從 `r` 中讀取數據, 直到 `EOF` 或發生錯誤

其返回值 `n` 為讀取的字節數, 除 `io.EOF` 之外, 在讀取過程中遇到的任何錯誤也將被返回

如果 `ReaderFrom` 可用, `Copy` 函數就會使用它

>❗️NOET:

`ReaderFrom` 方法不會返回 `err == EOF`

下面範例實現將文件中的資料全部讀取並顯示在標準輸出:

```go
file, err := os.Open("writeAt.txt")
if err != nil {
    panic(err)
}
defer file.Close()
writer := bufio.NewWriter(os.Stdout)
writer.ReadFrom(file)
writer.Flush()
```

當然也可以通過 `ioutil` package 的 `ReadFile` 函數獲取文件全部內容, `ioutil.ReadFile` 其實也是通過 `ReadFrom` 方法實現(使用 `bytes.Buffer`, 其實現 `ReaderFrom` interface)

若不使用 `ReadFrom` interfaace 而是用 `io.Reader` interface, 有兩個思路:
- 先獲取文件大小(File 的 Stat 方法), 再定義一個相同大小的 `[]byte`, 通過 `Read` 一次性讀取
- 定義一個小的 `[]byte`, 不斷調用 `Read` 直到遇到 `EOF` 並將所有讀取到的 `[]byte` 串接

`WriteTo` 定義如下:

```go
type WriterTo interface {
    WriteTo(w Writer) (n int64, err error)
}
```

> 官方文件對於該 interface 說明如下:

`WriteTo` 將數據寫入 `w` 中, 直到沒有數據可寫或發生錯誤, 其返回值 `n` 為寫入的字節數, 在寫入過程中遇到的任何錯誤也將被返回

如果 `WriterTo` 可用, `Copy` 函數就會使用它

以下程式碼示範將一段文本輸出到標準輸出:

```go
reader := bytes.NewReader([]byte("regy.dev"))
reader.WriteTo(os.Stdout)
```

> 如果需要一次性從某個地方讀或寫到某個地方, 可以考慮使用 `io.ReaderFrom` 和 `io.WriterTo`