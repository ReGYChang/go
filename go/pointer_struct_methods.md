- [Pointer](#pointer)
  - [Memory Allocation](#memory-allocation)
  - [Nil pointer & Wild pointer](#nil-pointer--wild-pointer)
  - [Pointer variable and Memory storage](#pointer-variable-and-memory-storage)
  - [Function parameter](#function-parameter)

# Pointer

Pointer 指代表某個記憶體地址的值, 這個記憶體地址往往是記憶體中一個變數的值的起始位置

Pointer 有幾個特性：
- default: nil
- `&` 取變數記憶體地址, `*` 通過 pointer 訪問目標物件
- 不支持指針運算, 不支持 `->` 運算符, 直接用 `.` 訪問目標物件屬性或方法

```go
package main

import "fmt"

func main(){ 
    var x int = 99
    var p *int = &x

    fmt.Println(p)
}
```

當運行到 `var x int = 99` 時記憶體中會分配一塊空間, 我們為它命名為 `x`. 他在記憶體中的有一個地址, 如 `0xc00000a0c8`. 當我們想訪問這個空間時, 可以使用**記憶體地址**或 `x` 去訪問

運行到 `var p *int = &x` 時, 我們定義一個**指針變數** `p`, `p` 儲存了變數 `x` 的記憶體地址

> 結論是指針就是記憶體地址, 指針變數就是儲存記憶體位置的變數

接著修改 `x` 內容
```go
package main

import "fmt"

func main() {
    var x int = 99
    var p *int = &x

    fmt.Println(p) // 0xc0000182d8

    x = 100

    fmt.Println("x: ", x) // 100
    fmt.Println("*p: ", *p) // 100
    
    *p = 999

    fmt.Println("x: ", x) // 999
    fmt.Println("*p: ", *p) // 999
}
```

- `x` & `*p` 的值是一樣的
- `*p` 稱為 `間接引用`
- `*p = 999` 通過 `x` 變數的記憶體地址來操作 `x` 對應的記憶體空間
- 不管是 `x` or `*p` 操作的都是同一塊記憶體空間

## Memory Allocation

![memory_allocation](img/memory_allocation.png)

其中 `.data` 存的是初始化後的資料

程式碼存在 stake, 一般 `make()` or `new()` 出來的存在 heap

Stack 用來給函式提供記憶體空間, 在 stack 上存取記憶體

函式調用時會在 call stack 上產生一個 stack frame, 每個獨立的 stack frame 一般包括：

- local variable
- parameter
- call function context

**其中 parameter 與 local variable 存儲地位相同**

當程式運行時首先執行 `main()`, 即產生一個 stack frame

當運行到 `var x int = 99` 時就會在 stack frame 中分配一塊記憶體位置

同理 `var p *int = &x`

## Nil pointer & Wild pointer

Nil pointer: 未被初始化的指針

```go
var p *int
```

這時若想取其值操作 `*p` 則會報錯

Wild pointer: 被一塊無效的記憶體地址空間初始化

```go
var p *int = 0xc00000a0c8
```

## Pointer variable and Memory storage

表達式 `new(T)` 會創建一個 `T` 類型的匿名變數, 為 `T` 類型的新值分配並清空一塊記憶體空間, 然後將這塊記憶體空間地址返回, 而這個結果就是指向這個新的 `T` 類型值的指針值, 返回的指針類型為 `*T`

`new()` 創建的記憶體空間位於 heap 上, 空間的默認值是資料類型的默認值. 如 `p := new(int)` 則 `*p` 為0
```go
package main

import "fmt"

func main(){
    p := new(int)
    fmt.Println(p)
    fmt.Println(*p)
}
```

只需使用 `new()` 函式創建而無需擔心記憶體的生命週期或怎樣將其空間關閉, go GC 會自動做記憶體管理

```go
package main

import "fmt"

func main(){
    p := new(int)
    
    *p = 1000
    
    fmt.Println(p)
    fmt.Println(*p)

    var x int = 10
    var y int = 20
    x = y
}
```

`var y int = 20` 中的 `y` 代表的是記憶體位置, 稱為`左值`; 而 `x = y` 中的 `y` 代表的是**記憶體空間中的值**, 稱作 `右值`

`x = y` 表示的是把 `y` 對應的記憶體空間值寫到 `x` 記憶體空間中

`=` 左邊的變數代表變數指向的**記憶體位置**, 相當於**寫操作**

`=` 右邊的變數代表變數**記憶體空間存的值**, 相當於**讀操作**

```go
p := new(int)

*p = 1000

fmt.Println(*p)
```

`*p = 1000` 意思是將 1000 寫入 *p 指向的記憶體空間中

而 `fmt.Println(*p)` 則是把 `*p` 的記憶體空間中的值打印出來

若在 funcion 中創建變數：

```go
func foo() {
    p := new(int)

    *p = 1000
}
```

當運行 `foo()` 會產生一個 stack frame, 運行結束會釋放 stack frame

而 `p` 是 `new()` 創建的, 存在在 `heap` 上, 當 stack frame 釋放時 `p` 並沒有消失, `p` 指向的記憶體位置也沒有消失, 於是可以基於這個特性實現 `call by reference`

## Function parameter

函式傳遞參數有兩種：

- pass by reference: 將記憶體位置作為函式參數傳遞
- pass by value: 將 argument(實參)的值拷貝一份給 parameter(行參)

> 無論是 pass by reference or value, 都是 `argument` 將自己的值拷貝給 `parameter`, 只不過值有可能是記憶體地址或是值

pass by value

```go
package main

import "fmt"

func swap(x, y int){
    x, y = y, x
    fmt.Println("swap  x: ", x, "y: ", y)
}

func main(){
    x, y := 10, 20
    swap(x, y)
    fmt.Println("main  x: ", x, "y: ", y)
}
```

output

```go
swap  x:  20 y:  10
main  x:  10 y:  20
```

- 運行 `main()` 時系統在 stack 產生一個 stack frame, 上面有 `x` 及 `y` 兩個變數
- 運行 `swap()` 時系統在 stack 產生一個 stack frame, 上面有 `x` 及 `y` 兩個變數
- 運行 `x, y = y, x` 後交換 `swap()` 產生的 stack frame 中的 `xy` 值, 不影響 `main()` 中的 `xy` 值
- `swap()` 運行完成後釋放 stack frame, 其上的 `x,y` 也隨之消失
- 當運行 `fmt.Println("main x: ", x, "y: ", y)` 時 `xy` 值依舊保持不變

pass by reference

```go
package main

import "fmt"

func swap2(a, b *int){
    *a, *b = *b, *a
}

func main(){
    x, y := 10, 20
    swap(x, y)
    fmt.Println("main  x: ", x, "y: ", y)
}
```

output

```go
main  x:  20 y:  10
```

- 運行 `main()` 並創建 stack frame, 上面有 `xy` 兩個變數
- 運行 `swap()` 並創建 stack frame, 上面有 `ab` 兩個變數
- `ab` 中存儲的值是 `xy` 的記憶體地址
- 當運行 `*a, *b = *b, *a` 時, 左邊的 `*a` 表示 `x` 的記憶體地址, 左邊的 `*b` 表示的是 `y` 的記憶體位置; 右邊的 `*b` 表示 `y` 的值, 右邊的 `*a` 表示的是 `x` 的值
- `main()` 中的 `xy` 被交換
- `swap()` 釋放 stack frame 也不影響被交換結果