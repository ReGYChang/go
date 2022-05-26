- [Defer](#defer)
  - [Defer Methods](#defer-methods)
  - [Arguments Evaluation](#arguments-evaluation)
  - [Defer Stack](#defer-stack)
  - [Defer Use Case](#defer-use-case)
- [Error](#error)
  - [Error Type](#error-type)
  - [Get More Infomation from Error](#get-more-infomation-from-error)
    - [Type Assertion for Struct - Struct Field](#type-assertion-for-struct---struct-field)
    - [Type Assertion for Struct - Call Methods](#type-assertion-for-struct---call-methods)
    - [Compare Error](#compare-error)

# Defer

`defer` 語句會將後面跟隨的語句進行延遲處理

在 `defer` 所屬的函數即將返回時，將延遲處理的語句按 `defer` 的反序 (LIFO) 進行執行，意即先被 defer 的語句最後被執行，最後被 defer 的語句最先被執行

```go
package main

import (  
    "fmt"
)

func finished() {  
    fmt.Println("Finished finding largest")
}

func largest(nums []int) {  
    defer finished()
    fmt.Println("Started finding largest")
    max := nums[0]
    for _, v := range nums {
        if v > max {
            max = v
        }
    }
    fmt.Println("Largest number in", nums, "is", max)
}

func main() {  
    nums := []int{78, 109, 2, 563, 300}
    largest(nums)
}
```

上述程式找出一個給定 slice 的最大值, `largest` 函數接收一個 int 類型的 slice 作為參數並打印該 slice 中的最大值

`largest()` 函數將要返回之前會先調用 `finished()` 函數

output:

```go
Started finding largest  
Largest number in [78 109 2 563 300] is 563  
Finished finding largest
```

## Defer Methods

`defer` 不僅限於函數調用, 也可以用於方法調用

```go
package main

import (  
    "fmt"
)


type person struct {  
    firstName string
    lastName string
}

func (p person) fullName() {  
    fmt.Printf("%s %s",p.firstName,p.lastName)
}

func main() {  
    p := person {
        firstName: "John",
        lastName: "Smith",
    }
    defer p.fullName()
    fmt.Printf("Welcome ")  
}
```

output:

```go
Welcome John Smith
```

## Arguments Evaluation

Go 中並非在調用 `defer` 函數時才確定 arguments value, 而是當執行 `defer` 語句時就會對 `defer` 函數的 arguments 進行求值

```go
package main

import (  
    "fmt"
)

func printA(a int) {  
    fmt.Println("value of a in deferred function", a)
}
func main() {  
    a := 5
    defer printA(a)
    a = 10
    fmt.Println("value of a before deferred function call", a)

}
```

`a` 初始值為 5, 執行到 `defer` 語句時由於 `a` 等於 5, 因此 `defer` 函數 `printA` 的 argument 也是 5

output:

```go
value of a before deferred function call 10  
value of a in deferred function 5
```

## Defer Stack

當一個函數內多次調用 `defer`, Go 會把 `defer` 調用放入一個 stack, 按照 LIFO 順序執行

```go
package main

import (  
    "fmt"
)

func main() {  
    name := "Naveen"
    fmt.Printf("Orignal String: %s\n", string(name))
    fmt.Printf("Reversed String: ")
    for _, v := range []rune(name) {
        defer fmt.Printf("%c", v)
    }
}
```

`for range` 會遍歷一個 string, 並調用 `defer fmt.Printf("%c, v")`

output:

```go
Orignal String: Naveen  
Reversed String: neevaN
```

## Defer Use Case

當一個函數應該在與當前 code flow 無關的環境下調用時, 可以使用 `defer`

```go
package main

import (  
    "fmt"
    "sync"
)

type rect struct {  
    length int
    width  int
}

func (r rect) area(wg *sync.WaitGroup) {  
    if r.length < 0 {
        fmt.Printf("rect %v's length should be greater than zero\n", r)
        wg.Done()
        return
    }
    if r.width < 0 {
        fmt.Printf("rect %v's width should be greater than zero\n", r)
        wg.Done()
        return
    }
    area := r.length * r.width
    fmt.Printf("rect %v's area %d\n", r, area)
    wg.Done()
}

func main() {  
    var wg sync.WaitGroup
    r1 := rect{-67, 89}
    r2 := rect{5, -67}
    r3 := rect{8, 9}
    rects := []rect{r1, r2, r3}
    for _, v := range rects {
        wg.Add(1)
        go v.area(&wg)
    }
    wg.Wait()
    fmt.Println("All go routines finished executing")
}
```

上述程式創建了 `rect` struct, 並創建了 `rect` 的方法 `area` 來計算矩形的面積

`main` 函數創建了 3 個 `rect` 類型變數: r1, r2, r3

將這 3 個變數添加到 `rects` slice 中並使用 `for range` 遍歷, 把 `area` 方法作為一個併發的 goroutine 進行調用

`WaitGroup` 作為參數傳遞給 `area` 方法後, 透過調用 `wg.Done` 通知 `main` 函數 goroutine 完成並返回

`wg.Done()` 應該在 `area` 將要返回之前調用, 且與 code flow path 無關, 因此只需調用一次 `defer` 來有效替換掉 `wg.Done()` 的多次調用

```go
package main

import (  
    "fmt"
    "sync"
)

type rect struct {  
    length int
    width  int
}

func (r rect) area(wg *sync.WaitGroup) {  
    defer wg.Done()
    if r.length < 0 {
        fmt.Printf("rect %v's length should be greater than zero\n", r)
        return
    }
    if r.width < 0 {
        fmt.Printf("rect %v's width should be greater than zero\n", r)
        return
    }
    area := r.length * r.width
    fmt.Printf("rect %v's area %d\n", r, area)
}

func main() {  
    var wg sync.WaitGroup
    r1 := rect{-67, 89}
    r2 := rect{5, -67}
    r3 := rect{8, 9}
    rects := []rect{r1, r2, r3}
    for _, v := range rects {
        wg.Add(1)
        go v.area(&wg)
    }
    wg.Wait()
    fmt.Println("All go routines finished executing")
}
```

output:

```go
rect {8 9}'s area 72  
rect {-67 89}'s length should be greater than zero  
rect {5 -67}'s width should be greater than zero  
All go routines finished executing
```

使用 `defer` 另一個好處是, 假設使用 `if` 就又給 `area` 方法添加了一條 return path, 需要確保這條增加的 return path 調用了 `wg.Done()`, 而使用 `defer` 調用 `wg.Done()` 則無需再為新的 return path 添加 `wg.Done()`

# Error

`error` 表示程式中出現了異常情況：比如說試圖打開一個文件, 而文件系統裡卻並沒有這個文件, 這種異常情況會使用 `error` 類型表示

如同其他 build-in type(int, float64), 錯誤值可以儲存在變數中或作為函數返回值等

試圖打開一個不存在的文件:

```go
package main

import (  
    "fmt"
    "os"
)

func main() {  
    f, err := os.Open("/test.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(f.Name(), "opened successfully")
}
```

調用 `os.Open()` 試圖打開路徑為 `/test.txt` 的文件

```go
func Open(name string) (file *File, err error)
```

若成功打開文件, `Open` 函數會返回一個 file handler 和一個零值(`nil`)的 error; 若打開文件時發生錯誤則會返回一個不為 `nil` 的 error

一般函數或方法返回 error 會作為最後一個值返回, 處理 error 時會將返回的 error 與 nil 比較, nil 值表示沒有 error 發生

output:

```go
open /test.txt: No such file or directory
```

## Error Type

`error` 是一個 interface type:

```go
type error interface {  
    Error() string
}
```

所有實現 `Error()` 方法的類型都可以當作是一個 error 類型, `Error()` 方法給出了錯誤的描述

`fmt.Println()` 在打印錯誤時會在內部調用 `Error() string` 方法來得到錯誤描述

## Get More Infomation from Error

有幾種方法可以透過 error 來獲取更多資訊

舉例前例查找錯誤的文件路徑, 直接解析錯誤的字符串:

```go
open /test.txt: No such file or directory
```

解析這條錯誤訊息雖然獲取了發生錯誤的文件路徑, 但這種方式很不優雅, 隨著語言版本的更新, 這條錯誤的描述隨時都有可能變化導致程式異常

Go library 提供了各種提取錯誤相關資訊的方法

### Type Assertion for Struct - Struct Field

仔細觀察 [Open()](https://pkg.go.dev/os#OpenFile) 文件, 其返回的錯誤類型是 `*PathError`

`PathError` 是 struct type, 在 library implementation 如下:

```go
type PathError struct {  
    Op   string
    Path string
    Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }
```

`*PathError` 通過宣告 `Error() string` 方法而實現 `error` interface

`Error() string` 將文件操作, 路徑及實際錯誤拼接並返回字符串, 於是得到了錯誤訊息:

```go
open /test.txt: No such file or directory
```

struct `PathError` 的 `Path` field 就有導致錯誤的文件路徑, 若修改錯誤處理邏輯:

```go
package main

import (  
    "fmt"
    "os"
)

func main() {  
    f, err := os.Open("/test.txt")
    if err, ok := err.(*os.PathError); ok {
        fmt.Println("File at path", err.Path, "failed to open")
        return
    }
    fmt.Println(f.Name(), "opened successfully")
}
```

這裡使用了 `Type Assertion` 來獲取 `error` interface underlying value, 接著使用 `err.Path` 來打印該路徑

output:

```go
File at path /test.txt failed to open
```

### Type Assertion for Struct - Call Methods

第二種獲得更多錯誤訊息的方法就是對 underlying type 進行斷言, 並通過調用該 struct type methods

standard library 中 `DNSError` struct type 定義如下:

```go
type DNSError struct {  
    ...
}

func (e *DNSError) Error() string {  
    ...
}
func (e *DNSError) Timeout() bool {  
    ... 
}
func (e *DNSError) Temporary() bool {  
    ... 
}
```

`DNSError` struct 還有 `Timeout() bool` 及 `Temporary() bool` 兩個方法, 它們返回一個 bool value, 指出該錯誤是由超時引起, 且是臨時性錯誤

透過斷言 `*DNSError` 類型並調用這些方法來確定該錯誤是臨時性錯誤還是由超時導致的:

```go
package main

import (  
    "fmt"
    "net"
)

func main() {  
    addr, err := net.LookupHost("golangbot123.com")
    if err, ok := err.(*net.DNSError); ok {
        if err.Timeout() {
            fmt.Println("operation timed out")
        } else if err.Temporary() {
            fmt.Println("temporary error")
        } else {
            fmt.Println("generic error: ", err)
        }
        return
    }
    fmt.Println(addr)
}
```

上述程式中試圖獲取 `golangbot123.com` (無效 domain name) 的 ip

通過 `*net.DNSError` 的類型斷言獲取到錯誤的 underlying value, 並分別檢查該錯誤是由超時引起還是一個臨時性的錯誤

本例中錯誤既不是臨時性錯誤, 也不是由超時引起的, 因此程式輸出為:

```go
generic error:  lookup golangbot123.com: no such host
```

### Compare Error

第三種獲取更多錯誤訊息的方式是與 `error` 類型的變數直接比較

`filepath` 中的 `Glob` 用於返回滿足 glob 模式的所有文件名, 如果模式寫的不對函數會返回一個錯誤 `ErrBadPattern`

`filepath` package 對 `ErrBadPattern` 定義如下:

```go
var ErrBadPattern = errors.New("syntax error in pattern")
```

`errors.New()` 用於創建一個新的錯誤

當模式不正確時, `Glob` 函數會返回 `ErrBadPattern`

```go
package main

import (  
    "fmt"
    "path/filepath"
)

func main() {  
    files, error := filepath.Glob("[")
    if error != nil && error == filepath.ErrBadPattern {
        fmt.Println(error)
        return
    }
    fmt.Println("matched files", files)
}
```

上述程式查詢了模式為 `[` 的文件, 為錯誤模式, 並檢查了該錯誤是否為 `nil`

將 `error` 直接與 filepath.ErrBadPattern 比較; 若條件滿足該錯誤就是由模式錯誤所導致

```go
syntax error in pattern
```