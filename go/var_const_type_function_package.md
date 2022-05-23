- [Variable](#variable)
  - [Anonymouse Variable](#anonymouse-variable)
- [Constant](#constant)
  - [Literal Constants(unnamed constants)](#literal-constantsunnamed-constants)
  - [iota](#iota)
- [Type](#type)
- [Function](#function)

# Variable

Go 採用靜態類型, 宣告變數時需指定變數的類型

變數宣告
```go
// 未賦值默認為 0
var a int 
// 宣告時賦值
var a int = 1
// 宣告時賦值
var a = 1
```

var a = 1, 由於 1 是 int 類型, 所以賦值時 a 自動被確認為 int 類型故可以省略類型名稱, 另一種更簡單表述:

```go
a := 1
msg := "Hello World"
```

這裡 `:=` 直接代替了變數的定義及賦值

## Anonymouse Variable

_ 為匿名變數, 會丟棄對應資料不處理. 通常配合函數返回值處理

```go
_, b := 3, 2
```

匿名變數不佔用 namespace, 且不會分配記憶體空間

# Constant

與變數定義相似, 將 `var` 換成 `const`, 且常數在定義時**必須賦值**

常數用來儲存不會發生變化的資料, 例如圓周率, 身分證號碼等, 在整個 runtime 不允許變動的值

```go
package main

import "fmt"

func main(){
    const pi float64 = 3.14159
    // pi = 4.56  compile error
    fmt.Println(pi)

    // 不是使用 :=
    const e = 2.7182
    fmt.Println("e =", e)
}
```

多個常數宣告
```go
const (
    pi = 3.14159
    e = 2.7182
)
```

`const` 同時宣告多個常數時若省略賦值則表示和上一行值相同

```go
const (
    n1 = 99
    n2  // n2 = 99
    n3  // n3 = 99
)
```

## Literal Constants(unnamed constants)

指程式中 hard coding 的 const

```go
23  // 整數類型常數
3.14159  // 浮點數類型常數
3.2+12i  // 復數類型常數
true  // 布林值類型常數
"foo"  // 字符串類型常數
```

## iota

`iota` 是 go 常數計數器, 只能在常數的表達式中使用. 用於生成一組相似規則初始化的常數, 但不用每行都寫一遍初始化表達式

>❗️ 在一個 `const` 宣告語句中, `iota` 將會被設為 0, 並在每一個有常數聲明的行加一

`iota` 可以理解為 `const` 的行索引, 能簡化定義, 在定義枚舉時很有效

```go
package main

import "fmt"

func main(){
    const (
        a = iota  // 0
        b  // 1
        c  // 2
        d  // 3
    )
    fmt.Println(a, b, c, d)
}
```

`iota` 遇到 `const` 會重置為 0

```go
package main

import "fmt"

func main(){
    const (
        a = iota
        b // 1
        c // 2
        d // 3
    )
    fmt.Println(a, b, c, d)
    // iota 遇到 const 重置為 0
    const e = iota  // 0
    fmt.Println(e)
}
```

使用 `_` 跳過某些值

```go
package main

import "fmt"

func main(){
    const (
        a = iota  // 0
        _
        c  // 2
        d  // 3
    )
    fmt.Println(a, c, d)
}
```

`iota` 宣告中插入值

```go
package main

import "fmt"

func main(){
    const (
        a = iota  // 0
        b = 100  // 100
        c = iota  // 2
        d  // 3
    )
    fmt.Println(a, b, c, d)
}
```

常數寫在同一行 iota 值相同, 下一行值 +1

```go
package main

import "fmt"

func main() {
    const(
        a = iota  // 0
        b, c = iota, iota  // 1, 1
        d, e  // 2, 2
        f, g, h = iota, iota, iota  // 3, 3, 3
        i, j, k  // 4, 4, 4
    )
    fmt.Println(a)
    fmt.Println(b, c)
    fmt.Println(d, e)
    fmt.Println(f, g, h)
    fmt.Println(i, j, k)
}
```

為常數賦初始值, 換行後 iota 根據行 +1, 不是根據值 +1

```go
package main

import "fmt"

func main(){
    const (
        a = 6  // 6
        b, c = iota, iota  // 1 1
        d, e  // 2 2
        f, g, h = iota, iota, iota  // 3 3 3
        i, j, k  // 4 4 4
    )
    fmt.Println(a)
    fmt.Println(b, c)
    fmt.Println(d, e)
    fmt.Println(f, g, h)
    fmt.Println(i, j, k)
}
```

若一行中賦值初始值不同, 則下一行的值與上一行的相同

```go
package main

import "fmt"

func main(){
    const (
        a, b = 1, 6  // 1 6
        c, d  // 1 6
        e, f, g = 2, 8, 10  // 2 8 10
        h, i, j  // 2 8 10
    )
    fmt.Println(a, b)
    fmt.Println(c, d)
    fmt.Println(e, f, g)
    fmt.Println(h, i, j)
}
```

若一行中既有賦初始值, 又有 `iota`, 則下一行中對應初始值的位置的值不變, 對應 `iota` 位置的值 +1

```go
package main

import "fmt"

func main(){
    const (
        a, b, c = 3, iota, iota  // 3 0 0
        d, e, f  // 3 1 1
        g, h, i = iota, 16, iota  // 2 16 2
        j, k, l  // 3 16 3
    )
    fmt.Println(a, b, c)
    fmt.Println(d, e, f)
    fmt.Println(g, h, i)
    fmt.Println(j, k, l)
}
```

定義數量級

```go
package main

import "fmt"

func main(){
    const (
        _ = iota
        KB = 1 << (10 * iota)  // 1024
        MB = 1 << (10 * iota)
        GB = 1 << (10 * iota)
        TB = 1 << (10 * iota)
        PB = 1 << (10 * iota)
    )
    fmt.Println(KB, MB, GB, TB, PB)
}
```

# Type

- int
  - 按空間大小分為: `int8`, `int16`, `int32`, `int64`
  - 對應無號整數: `uint8`, `uint16`, `uint32`, `uint64`
- float
  - `float32`, `float64`
- complex
  - 默認類型是 `complex128` 64 bits 實數 + 64 bits 虛數
  - 另一種 `complex64` 32 bits 實數 + 32 bits 虛數
    ```go
    var c1 complex
    c1 = 1 + 2i

    var c2 complex64
    c2 = 2 + 3i

    var c3 complex128
    c3 = 3 + 4i

    fmt.Println(c1)
    fmt.Println(c2)
    fmt.Println(c3)
    ```
- boolean
  - `true`, `false`
  - boolean default `fault`
  - `true`, `false` 均小寫
  - 不允許將 int 轉 bool
  - 無法進行數值運算, 或與其他類型轉換
- string
  - "", 會識別轉譯字符
  - ``, 不會識別轉譯字符, 以原生形式輸出
  - go 使用 UTF-8 encoding, 每個字符佔 1 byte
- alias
  - `byte` 是 `uint8` alias, 表示 ASCII 中一個字符 (1 byte)
  - `rune` 是 `int32` alias, 表示一個 UTF-8 字符, 處理中文日文或其他特殊字符 (4 bytes)

# Function

function 宣告使用關鍵字 `func`, 可有多個參數及多個返回值. package main 中的 func main() 約定為可執行程式的入口

```go
func funcName(param1 Type1, param2 Type2, ...) (return1 Type3, ...) {
    // body
}
```

```go
func add(num1 int, num2 int) int {
	return num1 + num2
}

func div(num1 int, num2 int) (int, int) {
	return num1 / num2, num1 % num2
}
func main() {
	quo, rem := div(100, 17)
	fmt.Println(quo, rem)     // 5 15
	fmt.Println(add(100, 17)) // 117
}
```