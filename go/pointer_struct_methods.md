- [Pointer](#pointer)
  - [Memory Allocation](#memory-allocation)
  - [Nil pointer & Wild pointer](#nil-pointer--wild-pointer)
  - [Pointer variable and Memory storage](#pointer-variable-and-memory-storage)
  - [Function parameter](#function-parameter)
- [Struct](#struct)
  - [Class in Go](#class-in-go)
  - [Struct Declaration](#struct-declaration)
  - [Recursive struct](#recursive-struct)
  - [Struct Instantiation](#struct-instantiation)
  - [Basic struct instance](#basic-struct-instance)
  - [Create pointer type struct](#create-pointer-type-struct)
  - [Unaddressable But Can Take Addresses](#unaddressable-but-can-take-addresses)
  - [Struct Initailization](#struct-initailization)
    - [field key-value pair](#field-key-value-pair)
    - [multiple value list](#multiple-value-list)
    - [Initalize anonymouse struct](#initalize-anonymouse-struct)
  - [Struct assignment](#struct-assignment)
  - [Struct comparision](#struct-comparision)
  - [Struct array and slice](#struct-array-and-slice)
  - [Struct as map value](#struct-as-map-value)

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

# Struct

Struct 是一種聚合的**資料類型**, 由 0 或多個任意類型的值聚合而成的實體. 每個值為 struct 的屬性

Struct 也是值類型, 可以通過 `new()` 來創建

組成 `struct` 類型的資料稱為 `field`
- `field` 擁有自己的類型和值
- `field` 必須 unique
- `field` 的類型也可以是 struct, 或該 struct 類型

## Class in Go

Go 中沒有 `class` 的概念, 也不支持 `class` 的繼承等物件導向的概念. Go 中 `struct` 與 `class` 都是複合結構體, 但 go 中 `struct` 的 `interface field(Nested interface)` 比物件導向具有更高的擴展性及靈活性

## Struct Declaration

使用關鍵字 `type` 可以將各種基本類型定義為自定義類型. **struct 是一種複合的基本類型**, 通過 `type` 定義為自定義類型後讓 struct 更易於使用

```go
type structTypeName struct {
    field1 field1Type
    field2 field2Type
    …
}
```

- struct type name: 標示自定義 struct 名稱, 同一個 package 下不能重複 namespace
- field1: 表示 struct field name, struct 中 field name msut unique
- field1Type: 表示 struct field 的具體類型

```go
type Student struct{
    id      int
    name    string
    age     int
    gender  int // 0 表示女生，1 表示男生
    addr    string
}
```

> 這裡 Student 地位等價於 int, byte, bool, string 等類型

## Recursive struct

Struct type 可以通過引用自身來定義. 這在定義 `linked list` 或 `binary tree node` 時特別有用, 此時 node 包含指向相鄰節點的鏈接 (address)

linked list
```go
type Node struct {
    Data float64
    Next *Node
}
```

doubly linked list
```go
type Node struct {
    Pre *Node
    Data float64
    Next *Node
}
```

binary tree
```go
type Tree struct {
    Left *Tree
    Data float64
    Right *Tree
}
```

## Struct Instantiation

Struct 定義只是一種記憶體佈局的描述, 只有當 struct 實體化後才會被分配記憶體空間, 因此必須定義 struct 並實體化後才能使用 struct 中的 field

實體化就是根據 struct 定義的 field 創建一份一樣格式的記憶體空間, instance & instance 間的記憶體是完全獨立的

Go 可以通過多種方式實體化 struct, 根據實際選擇不同方式

## Basic struct instance

Struct 本身也是一種類型, 可以像宣告 build-in type 一樣使用 `var` 宣告 struct type

```go
var structInstance structType
```

對 Student 進行初始化

```go
type Student struct{
    id      int
    name    string
    age     int
    gender  int // 0 表示女生，1 表示男生
    addr    string
}

func main() {
    var stu1 Student
    stu1.id = 120100
    stu1.name = "Conan"
    stu1.age = 18
    stu1.gender = 1
    fmt.Println("stu1 = ", stu1)  // stu1 =  {120100 Conan 18 1 }
}
```

>❗️未賦值的 `field` 默認為該 field type 零值, `addr = ""`

可以通過 `.` 訪問 struct member variable, 如 `stu1.name`

## Create pointer type struct

Go 中可以使用 `new` 對類型(struct, int, float, string) 進行實體化, struct 在實體化後會形成 **pointer type struct**

使用 `new` 實體化
```go
varName := new(type)
```

Go 可以像訪問一般 struct 一樣使用 `.` 訪問 pointer struct member
```go
type Student struct{
    id      int
    name    string
    age     int
    gender  int // 0 表示女生，1 表示男生
    addr    string
}

func main() {
    stu2 := new(Student)
    stu2.id = 120101
    stu2.name = "Kidd"
    stu2.age = 23
    stu2.gender = 1
    fmt.Println("stu2 = ", stu2)  // stu2 =  &{120101 Kidd 23 1 }
}
```

>❗️Go 中訪問 struct pointer member variable 可以繼續使用 `.`, 因為 go 在此設計了語法糖, 將 `stu2.name` 形式轉換為 `*stu2.name`. 不然這邊 `stu2` 表示的是儲存 Student 實體指針的的記憶體地址

## Unaddressable But Can Take Addresses

go 中對 struct 進行 `&` 操作時視為對該類型進行一次 `new` 實體化操作. 所有組合字面量都是**不可尋址**的, 但其可以被**取址**

```go
package main

import "fmt"

func main() {
	type Book struct {
		Pages int
	}
	var book = Book{} // 變數值 book 是可尋址的
	p := &book.Pages
	*p = 123
	fmt.Println(book) // {123}

	// 下面兩行 compile error，因為 Book{} 是不可尋址的，
	// 繼而 Book{}.Pages 也是不可尋址的。
	/*
	Book{}.Pages = 123
	p = &Book{}.Pages // <=> p = &(Book{}.Pages)
	*/
}
```

## Struct Initailization

Struct 在實體化時可以直接對成員變數進行初始化, 分兩種形式:
- field key-value pair
- multiple value list

### field key-value pair

field key-value pair 的初始化方式較適合選擇性填充 field 教多的 struct

key-value pair 填充是可選的, 不需初始化的 field 可以不填

struct 實體化後 field default value 是 field type 零值, 如 int 0, string "", bool false, pointer nil 等

```
varName := structTypeName {
    field1: value1,
    field2: value2,
    ...
}
```

> ❗️Note
- field name unique
- key value 之間以 `:` 分隔, key value pair 之間以 `,` 分隔

```go
stu4 := Student{
    id:     120103,
    name:   "Gin",
    age:    25,
    gender: 1,
    addr:   "unknown",
}
fmt.Println("stu4 = ", stu4) // stu4 =  {120103 Gin 25 1 unknown}
```

### multiple value list
multiple value list 的初始化方式較適合填充 field 較少的 struct

go 可以在 key-value pair 初始化的基礎上忽略 **"key"** , 使用多個 value 的列表初始化 struct field

```go
varName := structTypeName {
    field1Val,
    field2Val,
    ...
}
```

> ❗️Note
- 必須初始化 struct 所有 fields
- 每個初始值的填充順序必須與 field 在 struct 中定義的順序一樣
- key-value pair 與 value list 初始化型態不能混用

```go
stu5 := Student{
    120104,
    "Kogorou",
    38,
    1,
    "區塊鏈革命",
}
fmt.Println("stu5 = ", stu5) // stu5 =  {120104 Kogorou 38 1 區塊鏈革命}
```

### Initalize anonymouse struct

Anonymouse struct 沒有類型名稱, 不需透過 `type` 關鍵字定義就可以直接使用

```go
package main

import (
    "fmt"
)

func main() {
    var user struct{name string; age int}
    user.name = "Conan"
    user.age = 18
    fmt.Println("user = ", user)  // user =  {Conan 18}
}
```

## Struct assignment

當使用 `=` 對 struct 賦值時, 更改一個 struct field 值不會影響其他值

```go
package main

import (
    "fmt"
)

type Student struct {
    id     int
    name   string
    age    int
    gender int // 0 表示女生，1 表示男生
    addr   string
}

func main() {
    var stu1 Student
    stu1.id = 120100
    stu1.name = "Conan"
    stu1.age = 18
    stu1.gender = 1

    stu6 := stu1
    fmt.Println("stu1 = ", stu1) // stu1 =  {120100 Conan 18 1 }
    fmt.Println("stu6 = ", stu6) // stu6 =  {120100 Conan 18 1 }

    stu6.name = "柯南"
    fmt.Println("stu1 = ", stu1) // stu1 =  {120100 Conan 18 1 }
    fmt.Println("stu6 = ", stu6) // stu6 =  {120100 柯南 18 1 }
}
```

## Struct comparision

若 struct 全部成員都是可比較的, sturct 也可以比較; 若 struct 中存在不可比較的成員變數如 slice, map 等, 那麼 struct 則無法比較. 如果用 `==`, `!=` 進行判斷程式會直接抱錯, 可以用 `DeepEqual` 來進行深度比較

```go
type Student struct {
    id   int
    name string
}

func main() {
    var stu1 Student
    stu1.id = 120100
    stu1.name = "keke"

    stu6 := stu1
    stu6.name = "顆顆"

    fmt.Println(stu1.id == stu6.id && stu1.name == stu6.name) // "false"
    fmt.Println(stu1 == stu6)                                 // "false"
}
```

## Struct array and slice

假設要用 struct 儲存多個學生資料, 可以定義 struct array 來儲存並通過 loop 方式循環輸出

```go
package main

import "fmt"

type student struct {
    id   int
    name string
    score  int
}

func main() {
    // struct array
    students := [3]student{
        {101, "conan", 88},
        {102, "kidd", 78},
        {103, "lan", 98},
    }
    // print each element in struct array
    for index, stu := range students {
        fmt.Println(index, stu.name)
    }
}
```

>用 struct slice 同理

## Struct as map value

```go
package main

import "fmt"

type student struct {
    id    int
    name  string
    score int
}

func main() {
    // 結構體陣列
    students := [3]student{
        {101, "conan", 88},
        {102, "kidd", 78},
        {103, "lan", 98},
    }
    // 輸出結構體陣列每一項
    for index, stu := range students {
        fmt.Println(index, stu.name)
    }
    fmt.Println(students)
    // 計算以上學生成績總分
    sum := students[0].score
    for i, stuLen := 1, len(students); i < stuLen; i++ {
        sum += students[i].score
    }
    fmt.Println("總分是：", sum)

    // 輸出以上學生成績最高分
    maxScore := students[0].score
    for i, stuLen := 1, len(students); i < stuLen; i++ {
        if maxScore < students[i].score {
            maxScore = students[i].score
        }
    }
    fmt.Println("最高分是：", maxScore)
}

func main() {
    // define map
    m := make(map[int]student)
    m[101] = student{101, "conan", 88}
    m[102] = student{102, "kidd", 78}
    m[103] = student{103, "lan", 98}
    fmt.Println(m) // map[101:{101 conan 88} 102:{102 kidd 78} 103:{103 lan 98}]

    for k, v := range m {
        fmt.Println(k, v)
    }
}
```