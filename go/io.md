- [I/O](#io)
  - [Reader Interface](#reader-interface)
  - [Writer Interface](#writer-interface)
  - [Types Implement io.Reader and io.Writer](#types-implement-ioreader-and-iowriter)
  - [ReaderAt 和 WriterAt interface](#readerat-和-writerat-interface)
  - [ReaderFrom & WriterTo interface](#readerfrom--writerto-interface)
  - [Seeker interface](#seeker-interface)
  - [Closer interface](#closer-interface)
- [ioutil](#ioutil)
  - [NopCloser](#nopcloser)
  - [ReadAll](#readall)
  - [ReadDir](#readdir)
  - [ReadFile & WriteFile](#readfile--writefile)
- [fmt](#fmt)
  - [Printing](#printing)
    - [Sample](#sample)
    - [Placeholder](#placeholder)
- [Encoding/Decoding](#encodingdecoding)
  - [encoding/json](#encodingjson)
  - [JSON Encoding](#json-encoding)
    - [Data Type Mapping](#data-type-mapping)
  - [JSON Decoding](#json-decoding)
    - [Data Type Mapping](#data-type-mapping-1)
  - [Decode Unknown JSON data](#decode-unknown-json-data)
  - [Visit Decoding JSON data](#visit-decoding-json-data)
  - [Decode JSON From Stream](#decode-json-from-stream)
  - [omitempty](#omitempty)

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

## Seeker interface

`Seeker` interface 定義如下:

```go
type Seeker interface {
    Seek(offset int64, whence int) (ret int64, err error)
}
```

> 官方文件關於此 interface methods 說明:

Seek 設置下一次 Read 或 Write 的偏移量為 offset, 它的解釋取決於 whence 
- 0 表示相對於文件的起始處
- 1 表示相對於當前的偏移
- 2 表示相對於其結尾處, Seek 返回新的偏移量和一個錯誤(如果有的話)

也就是說 `Seek` 方法是用於設置偏移量的, 這樣可以從某個特定位置開始操作資料流

看起來跟 `ReaderAt/WriterAt` 有些類似, 不過 `Seeker` interface 更加靈活, 可以更好的控制讀寫資料流位置

簡單範例程式碼: 獲取倒數第二個字符(須考慮 UTF-8 編碼)

```go
reader := strings.NewReader("今天天氣真好")
reader.Seek(-6, io.SeekEnd)
r, _, _ := reader.ReadRune()
fmt.Printf("%c\n", r)
```

>💡TIPS

`whence` 值在 `io` package 中定義了相應的常數:

```go
const (
  SeekStart   = 0 // seek relative to the origin of the file
  SeekCurrent = 1 // seek relative to the current offset
  SeekEnd     = 2 // seek relative to the end
)
```

原先 `os` package 中的常數已經被標註為 Deprecated:

```go
// Deprecated: Use io.SeekStart, io.SeekCurrent, and io.SeekEnd.
const (
  SEEK_SET int = 0 // seek relative to the origin of the file
  SEEK_CUR int = 1 // seek relative to the current offset
  SEEK_END int = 2 // seek relative to the end
)
```

## Closer interface

`Closer` interface 定義如下:

```go
type Closer interface {
    Close() error
}
```

`Closer` interface 只有一個 `Close()` 方法, 用於關閉資料流

`os.File`, `compress` package, 資料庫連線, `Socket` 等需要手動關閉的資源都實現了 `Closer` interface

實際場景中經常將 `Close` 方法調用放在 `defer` 語句中

# ioutil

雖然 `io` package 提供了不少型別, 方法和函數, 但有時使用起來不是很方便, 比如讀取一個文件中所有的內容

為此, `ioutil` 中提供了一些常用的 I/O 操作函數

## NopCloser

有時候需要傳遞一個 `ioReadCloser` instance, 而目前只有一個 `io.Reader` instance, 如 `strings.Reader`, 此時就需要 `NopCloser` 來包裝轉換成 `ioReadCloser`

如在 `net/http` package 中的 `NewRequest`, 需要接收一個 `io.Reader` 參數, 但實際上 `http.Request` 的 `Body` 是 `io.ReaderCloser` 型別

若傳遞的 `io.Reader` 也實現了 `io.ReadCloser` interface 則直接轉換, 否則可以透過 `ioutil.NopCloser` 來包裝:

```go
rc, ok := body.(io.ReadCloser)
if !ok && body != nil {
    rc = ioutil.NopCloser(body)
}
```

## ReadAll

很多時候需要一次性讀取 `io.Reader` 中的資料, Go 中提供了 `ReadAll` 這個函數, 用來從 `io.Reader` 中一次性讀取所有資料

`ReadAll` 函數定義如下:

```go
func ReadAll(r io.Reader) ([]byte, error)
```

其是通過 `bytes.Buffe` 中的 `ReadFrom` 來實現讀取所有資料, 該函數成功調用後會返回 `err == nil` 而不是 `err == EOF`(無錯誤不處理)

## ReadDir

`ioutil.ReadDir` 會讀取目錄並返回排序好的文件與子目錄名([]os.FileInfo):

```go
func main() {
    dir := os.Args[1]
    listAll(dir,0)
}

func listAll(path string, curHier int){
    fileInfos, err := ioutil.ReadDir(path)
    if err != nil{fmt.Println(err); return}

    for _, info := range fileInfos{
        if info.IsDir(){
            for tmpHier := curHier; tmpHier > 0; tmpHier--{
                fmt.Printf("|\t")
            }
            fmt.Println(info.Name(),"\\")
            listAll(path + "/" + info.Name(),curHier + 1)
        }else{
            for tmpHier := curHier; tmpHier > 0; tmpHier--{
                fmt.Printf("|\t")
            }
            fmt.Println(info.Name())
        }
    }
}
```

## ReadFile & WriteFile

`ReadFile` 會讀取整個文件內容, 其實現與 `ReadAll` 類似, 不過 `ReadFile` 會先判斷文件大小, 給 `bytes.Buffer` 一個預定義的容量以避免額外分配記憶體

`ReadFile` 函數簽名如下:

```go
func ReadFile(filename string) ([]byte, error)
```

> 函數官方文件說明:

`ReadFile` 從指定文件中讀取資料並返回, 成功的調用返回 `error == nil` 而非 `error == EOF`, 因為本函數定義為讀取整個文件, 其不會將讀取返回的 `EOF` 視為應報告的錯誤(同 `ReadAll`)

`WriteFile` 函數簽名如下:

```go
func WriteFile(filename string, data []byte, perm os.FileMode) error
```

> 函數官方文件說明:

`WriteFile` 從資料寫入指定文件中, 當文件不存在時會根據 `perm` 指定的權限創建文件, 文件存在時會先清空文件內容, 對於 `perm` 參數一般可以指定為 `0666`

>💡TIPS:

`ReadFile` source code 中先獲取文件大小, 當大小 < 1e9 時才會用到文件的大小, 按 source code 注視的說法 `FileInfo` 不會很精確地獲取文件大小

# fmt

`fmt` package 實現了格式化 I/O 函數, 類似於 C 的 `printf` 和 `scanf`

`fmt` package 的官方文件中對 `Printing` 和 `Scanning` 有很詳細的說明, 這裡直接引用文件說明

以下範例中使用到的型別或變數定義:

```go
type Website struct {
    Name string
}

// define struct variable
var site = Website{Name:"studygolang"}
```

## Printing

### Sample

```go
type user struct {
    name string
}

func main() {
    u := user{"tang"}
    //Printf 格式化輸出
    fmt.Printf("% + v\n", u)     //格式化輸出結構
    fmt.Printf("%#v\n", u)       //輸出值的 Go 語言表示方法
    fmt.Printf("%T\n", u)        //輸出值的類型的 Go 語言表示
    fmt.Printf("%t\n", true)     //輸出值的 true 或 false
    fmt.Printf("%b\n", 1024)     //二進位表示
    fmt.Printf("%c\n", 11111111) //數值對應的 Unicode 編碼字符
    fmt.Printf("%d\n", 10)       //十進位表示
    fmt.Printf("%o\n", 8)        //八進位表示
    fmt.Printf("%q\n", 22)       //轉化為十六進制並附上單引號
    fmt.Printf("%x\n", 1223)     //十六進位表示, 用a-f表示
    fmt.Printf("%X\n", 1223)     //十六進位表示, 用A-F表示
    fmt.Printf("%U\n", 1233)     //Unicode表示
    fmt.Printf("%b\n", 12.34)    //無小數部分, 兩位指數的科學計數法6946802425218990p-49
    fmt.Printf("%e\n", 12.345)   //科學計數法, e表示
    fmt.Printf("%E\n", 12.34455) //科學計數法, E表示
    fmt.Printf("%f\n", 12.3456)  //有小數部分, 無指數部分
    fmt.Printf("%g\n", 12.3456)  //根據實際情況採用%e或%f輸出
    fmt.Printf("%G\n", 12.3456)  //根據實際情況採用%E或%f輸出
    fmt.Printf("%s\n", "wqdew")  //直接輸出字符串或者[]byte
    fmt.Printf("%q\n", "dedede") //雙引號括起來的字符串
    fmt.Printf("%x\n", "abczxc") //每個字節用兩字節十六進位表示, a-f表示
    fmt.Printf("%X\n", "asdzxc") //每個字節用兩字節十六進位表示, A-F表示
    fmt.Printf("%p\n", 0x123)    //0x開頭的十六進制數表示
}
```

### Placeholder

| format | meaning                                                                                                                                                    |
| ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| %%     | A% literal                                                                                                                                                 |
| %b     | A binary integer value (Radix 2), or a floating-point number with exponent 2 expressed in scientific counting                                              |
| %c     | Character type. You can input numbers according toASCII codeThe corresponding character is converted to the corresponding character                        |
| %d     | A decimal value (base 10)                                                                                                                                  |
| %e     | A floating-point number or complex number represented by scientific notation E                                                                             |
| %E     | A floating-point number or complex number represented by scientific notation E                                                                             |
| %f     | A floating point number or complex number represented by standard counting                                                                                 |
| %g     | Floating point number or complex number represented by% e or% F, either of which is output in the most compact way                                         |
| %G     | Floating point number or complex number represented by% e or% F, either of which is output in the most compact way                                         |
| %o     | An octal numeric value (base 8)                                                                                                                            |
| %p     | The address of a value in hexadecimal (base 16), prefixed with0x, the letters are lowercasea - fexpress.                                                   |
| %q     | Use go syntax and escape when necessary, string or byte slice [] byte enclosed in double quotation marks, or number enclosed in single quotation marks     |
| %s     | String. Outputs the characters in the string up to the empty characters in the string (string in characters)\0End, this\0(i.e. null character)             |
| %t     | withtrueperhapsfalseBoolean value of output                                                                                                                |
| %T     | The type of value output using go syntax                                                                                                                   |
| %U     | An integer code point represented in Unicode notation. The default value is 4 numeric characters                                                           |
| %v     | The built-in or user-defined value output in the default format, or the user-defined value output in the string () mode of its type, if this method exists |
| %x     | Integer value in hexadecimal (base is hexadecimal), numbera-fUse lowercase representation                                                                  |
| %X     | Integer value in hexadecimal (base is hexadecimal), numberA-FUse uppercase                                                                                 |

# Encoding/Decoding

## encoding/json

Go 語言 build-in package `encoding/json` 提供了一系列方法進行 json 編解碼, 下列逐一介紹這些方法

## JSON Encoding

可以通過 `encoding/json` package 的 `Marshal()` 函數將資料編碼為 JSON 文本, 此函數簽名如下:

```go
func Marshal(v interface{}) ([]byte, error)
```

傳入參數 `v` 為 interface, 意味著可以傳入任意型別的資料, 若編碼成功則返回對應的 JSON 格式文本, 失敗則透過 `error` 參數標示錯誤訊息

假設有一個 `User` 型別的 struct:

```go
type User struct { 
    Name string
    Website string
    Age  uint
    Male bool
    Skills []string
}
```

並通過以下形式對其初始化:

```go
user := User{
    "regy",
	"https://regy.dev",
	18,
	true,
	[]string{"Golang", "PHP", "C", "Java", "Python"},
}
```

隨後可以使用 `json.Marshal()` 函數將上述 `user` instance 編碼成 JSON 文本:

```go
u, err := json.Marshal(user)
```

完整程式碼如下:

```go
# src/note/json/basic.go
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string
	Website string
	Age  uint
	Male bool
	Skills []string
}

func main()  {
	user := User{
		"regy",
		"https://regy.dev",
		18,
		true,
		[]string{"Golang", "PHP", "C", "Java", "Python"},
	}

	u, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("JSON encoding failed: %v\n", err)
		return
	}

	fmt.Printf("JSON data: %s\n", u)
}
```

若編碼成功則 `err` 為 `nil`, output:

```go
JSON data: {"Name":"regy","Website":"https://regy.dev","Age":18,"Male":true,"Skills":["Golang","PHP","C","Java","Python"]}
```

底層實現邏輯為當調用 `json.Marshal(user)` 時會遍歷 `user` struct, 若發現 `user` 資料結構實現了 `json.Marshaler` interface 且包含有效的值, `Marshal()` 則會調用 `MarshalJSON()` 方法透過此資料結構生成 JSON 格式文本

### Data Type Mapping

除了 `channel`, `complex` 和函數幾種型別以外, Go 中大部分的資料型別都可以編碼為有效的 JSON 文本, 若編碼前的資料結構中出現 pointer, 則將編碼 pointer 所指向的值; 若指向零值, 則 `null` 將作為編碼後的結果

Go 中 JSON 編碼前後的資料型別映射如下:
- `bool` 編碼後依舊為 `bool`
- `float` 和 `int` 編碼後為 JSON 常規數字
- `string` 將以 `UTF-8` 編碼輸出為 `Unicode` 字符串, 特殊字符將為被轉義為 `\u003c`
- `array` 和 `slice` 會編碼為 JSON 中的 array, 但 `[]byte` 型別的值會被編碼為 `Base64` 編碼後的字符串, `slice` 型別的零值會被編碼為 `null`
- `struct` 會被編碼為 JSON Object, 且只有 struct 中大寫字母開頭的 field 才會被編碼輸出為 JSON Object string key
- 編碼一個 `map` 型別資料結構時, 其型別必須為 `map[string]T`(T 為 `encoding/json 支持的任意資料型別`)

## JSON Decoding

與 `json.Marshal()` 相對, 可以使用 `json.Unmarshal()` 函數將 JSON 解碼為 Go 中對應的資料結構

`json.Unmarshal()` 函數簽名如下:

```go
func Unmarshal(data []byte, v interface{}) error
```

解碼 JSON 資料前首先需要在 Go 中宣告一個目標型別的實體物件用於解碼後的值:

```go
var user User
```

再調用 `json.Unmarshal()` 函數將 `[]byte` 型別的 JSON 資料作為第一個參數傳入, 將 `user` 實體變數指針作為第二個參數傳入:

```go
err := json.Unmarshal(u, &user)
```

若 `u` 為有效 JSON 資料並能與 `user` struct 對應, 則 JSON 解碼後的值會一一存放到 `user` struct 對應 field 中

JSON decoding sample code:

```go
...

func main()  {
	...

	u, err := json.Marshal(user)
	...

	var user2 User
	err = json.Unmarshal(u, &user2)
	if err != nil {
		fmt.Printf("JSON decoding failed: %v\n", err)
	}
		return

	fmt.Printf("JSON decoding result: %#v\n", user2)
}
```

解碼成功後的 `user2` 資料如下:

```go
JSON decoding result: main.User{Name:"regy", Website:"https://regy.dev", Age:0x12, Male:true, Skills:[]string{"Golang", "PHP", "C", "Java", "Python"}}
```

### Data Type Mapping

實際上, json.Unmarshal() 函數會根據一個約定的順序查找目標結構中的字段, 如果找到一個即發生匹配

假設某個 JSON 對像有一個名為 `Foo` 的索引 (不區分大小寫), 要將 `Foo` 所對應的值填充到目標結構體的目標字段上, `json.Unmarshal()` 將會遵循如下順序進行查找匹配:

- 一個包含 Foo 標籤的字段 (不區分大小寫)
- 一個名為 Foo 或者除了首字母其他字母不區分大小寫的名為 Foo 的字段 (這些字段在類型聲明中必須都是以大寫字母開頭、可被外部訪問的公開字段)

後面兩個比較好理解, 第一個我們在微服務架構教程中通過 protoc 生成的原型文件裡面經常可以看到:

```go
type User struct {
    Id  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
    Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
    Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email"`
    Password string `protobuf:"bytes,4,opt,name=password,proto3" json:"password"`
}
```

這裡的 Name 被打上 `json:"name"` 標籤

當 JSON 資料結構和 Go 語言裡邊的目標型別的結構對不上時, 會發生什麼呢?

```go
u2 := []byte(`{"name": "regy", "website": "https://regy.dev", "alias": "GOGOGO"}`)
var user3 User
err = json.Unmarshal(u2, &user3)
if err != nil {
fmt.Printf("JSON decoding failed: %v\n", err)
	return
}
fmt.Printf("JSON decoding result: %#v\n", user3)
```

output:

```go
JSON decoding result: main.User{Name:"regy", Website:"https://regy.dev", Age:0x0, Male:false, Skills:[]string(nil)}
```

可以看到, 如果 JSON 中的字段在 Go 語言對應目標型別中不存在, `json.Unmarshal()` 函數在解碼過程中會丟棄該字段

上述程式碼中由於 `Alias` 字段並沒有在 `User` 型別中定義, 所以會被忽略, 只有 `Name` 和 `Website` 這兩個字段的值才會被填充到 `user3` 中

這個特性讓我們可以從同一段 JSON 數據中篩選指定的值填充到多個不同的 Go 語言型別中

對於 JSON 中沒有而 `User` 中定義的字段, 會以對應數據類型的默認值填充, 比如上述 `Age`, `Male`, `Skills` 字段均是如此

以上是在 JSON 結構已知情況下的解碼, 如果 JSON 結構是動態的或未知的, 又該怎麼處理呢?

## Decode Unknown JSON data

在 Go 中可以通過 `interface{}` 來表示任意資料型別, 同樣適用於未知結構的 JSON 資料解碼: 只需將這段 JSON 解碼輸出到一個 `interface{}` 型別的值即可

在實際解碼過程中, JSON 結構中的資料元素將做以下型別轉換:

- `boolean` 將會轉換為 Go 語言的 `bool` 型別
- `數值`會被轉換為 Go 語言的 `float64` 型別
- `字符串`轉換後還是 `string` 型別
- `JSON Array` 會轉換為 `[]interface{}` 型別
- `JSON Object` 會轉換為 `map[string]interface{}` 型別
- `null` 值會轉換為 `nil`

在 Go `encoding/json` 中允許使用 `map[string]interface{}` 和 `[]interface{}` 型別值來分別存放未知結構的 `JSON Object` 或 `Array`

這次將解碼結果映射到 `interface{}` 物件:

```go
u3 := []byte(`{"name": "regy", "website": "https://regy.dev", "age": 18, "male": true, "skills": ["Golang", "PHP"]}`)
var user4 interface{}
err = json.Unmarshal(u3, &user4)
if err != nil {
    fmt.Printf("JSON decoding failed: %v\n", err)
    return
}
fmt.Printf("JSON decoding result: %#v\n", user4)
```

上述程式碼中 `user4` 被定義為一個 `interface{}`, `json.Unmarshal()` 函數將一個 JSON Object 物件 `u3` 解碼輸出到 `user4` 中, 最終 `user4` 將會是一個 key-value pair 的 `map[string]interface{}` 結構:

```go
map[string]interface {}{"age":18, "male":true, "name":"regy", "skills":[]interface {}{"Golang", "PHP"}, "website":"https://regy.dev"}
```

`u3` 為一個 JSON Object, 內部屬性也會遵循上述型別轉換規則一一轉換

## Visit Decoding JSON data

取得解碼後的資料結構前須先判斷目標結構是否為預期的資料型別, 可通過 `for loop` 一一獲取解碼後的目標資料:

```go
user5, ok := user4.(map[string]interface{})
if ok {
    for k, v := range user5 {
        switch v2 := v.(type) {
        case string:
            fmt.Println(k, "is string", v2)
        case int:
            fmt.Println(k, "is int", v2)
        case bool:
            fmt.Println(k, "is bool", v2)
        case []interface{}:
            fmt.Println(k, "is an array:")
            for i, iv := range v2 {
                fmt.Println(i, iv)
            }
        default:
            fmt.Println(k, "is another type not handle yet")
        }
    }
}
```

## Decode JSON From Stream

此外 `encoding/json` 還提供了 `Decoder` 和 `Encoder` 兩個型別, 用於支持 JSON 資料的 stream read/write, 並提供 `NewDecoder()` 和 `NewEncoder()` 兩個函數用於實現:

```go
func NewDecoder(r io.Reader) *Decoder 
func NewEncoder(w io.Writer) *Encoder
```

以下演示從 `Stdin` input stream 中讀取 JSON 資料並將其解碼, 最後再寫入 `Stdout`:

```go
# src/note/json/stream.go
package main

import (
    "encoding/json"
    "log"
    "os"
)

func main() {
    dec := json.NewDecoder(os.Stdin)
    enc := json.NewEncoder(os.Stdout)
    for {
        var v map[string]interface{}
        if err := dec.Decode(&v); err != nil {
            log.Println(err)
            return
        }
        if err := enc.Encode(&v); err != nil {
            log.Println(err)
        }
    }
}
```

須先輸入 JSON 結構資料供 `Stdin` input stream 讀取, 再通過 `json.NewDecoder` 返回的 decoder 對其進行解碼, 最後通過 `json.NewEncoder` 返回的編碼器將資料編碼後寫入 `Stdout` output stream 打印

使用 `Decoder` 和 `Encoder` 對 data stream 進行處理應用更廣泛, 比如讀寫 `HTTP`, `WebSocket` 或檔案等, Go 標準庫 `net/rpc/jsonrpc` 就是應用了 `Decoder` 和 `Encoder` 的實際例子:

```go
// NewServerCodec returns a new rpc.ServerCodec using JSON-RPC on conn.
func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
    return &serverCodec{
        dec:     json.NewDecoder(conn),
        enc:     json.NewEncoder(conn),
        c:       conn,
        pending: make(map[uint64]*json.RawMessage),
    }
}
```

## omitempty

在定義 json struct 時經常會使用到 `omitempty`, 使用上需要特別注意

定義每個 field 對應的 json format, 如定義一個 `Dog` struct:

```go
type Dog struct {
	Breed string
	WeightKg int
}
```

將其初始化並序列化為 JSON 格式:

```go
func main() {
	d := Dog{
		Breed:    "dalmation",
		WeightKg: 45,
	}
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
```

output:

```go
{"Breed":"dalmation","WeightKg":45}
```

若有其中一個 struct member 沒有初始化則結果可能與預期不符:

```go
func main() {
	d := Dog{
        Breed:    "pug",
	}
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
```

output:

```go
{"Breed":"pug","WeightKg":0}
```

預期希望 `dog` 的 `weitght` 為未知, 而不是 0, 為了將未知的 field 忽略, 則使用 `omitempty` tag 來實現:

```go
type Dog struct {
	Breed    string
	// The first comma below is to separate the name tag from the omitempty tag 
	WeightKg int `json:",omitempty"`
}
```

output:

```go
{"Breed":"pug"}
```

現在 `WeightKg` 則被設置為默認零值, `int` 為 0, `string` 為 "", `pointer` 為 nil

當使用 struct 互相嵌套時, 則 `omitempty` 則可能出現預期外結果:

```go
type dimension struct {
	Height int
	Width int
	}

type Dog struct {
	Breed    string
	WeightKg int
	Size dimension `json:",omitempty"`
}

func main() {
	d := Dog{
		Breed: "pug",
	}
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
```

output:

```go
{"Breed":"pug","WeightKg":0,"Size":{"Height":0,"Width":0}}
```

這裡雖然使用 `omitempty` tag 但 `dimension` 還是沒有被忽略, 因為不知道 struct `dimension` 空值為何, Go 只能判斷如 `int`, `string`, `pointer` 這種基本型別的空值

為了使自定義 struct 未初始化時能被忽略, 可以使用 struct pointer:

```go
type Dog struct {
	Breed    string
	WeightKg int
	// Now `Size` is a pointer to a `dimension` instance
	Size *dimension `json:",omitempty"`
}
```

output:

```go
{"Breed":"pug","WeightKg":0}
```

Go 有能力判斷 pointer type 空值為 nil, 所以直接賦值

```go
type Dog struct {
	Age *int `json:",omitempty"`
}

func main() {
	age := 0
	d := Dog{
		Age: &age,
	}

	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
```

output:

```go
{"Age":0}
```

因為 `Age` 為 pointer 型別, 而初始化時也傳入非 nil pointer, 所以不會被忽略而顯示 0

```go
type Restaurant struct {
	NumberOfCustomers int `json:",omitempty"`
}

func main() {
	d := Restaurant{
		NumberOfCustomers: 0,
	}
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}
```

output:

```go
{}
```

這裡 `NumberOfCustomers` 被忽略顯時不是預期的結果, 因為 Go 將 0 當成零值, 跟沒賦值結果相同

解決方法之一即是使用 `int pointer` 並傳入非 nil pointer:

```go
type Restaurant struct {
	NumberOfCustomers *int `json:",omitempty"`
}

func main() {
	d1 := Restaurant{}
	b, _ := json.Marshal(d1)
	fmt.Println(string(b))
	//Prints: {}
	
	n := 0
	d2 := Restaurant{
		NumberOfCustomers: &n,
	}
	b, _ = json.Marshal(d2)
	fmt.Println(string(b))
	//Prints: {"NumberOfCustomers":0}
}
```