- [Struct](#struct)
  - [Class in Go](#class-in-go)
  - [Struct Declaration](#struct-declaration)
  - [Recursive Struct](#recursive-struct)
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
  - [Struct slice as map value](#struct-slice-as-map-value)
  - [Struct as function parameter](#struct-as-function-parameter)
  - [Composition Instead of Inheritance](#composition-instead-of-inheritance)
- [Method](#method)

# Struct

Struct 是一種聚合的**資料類型**, 由 0 或多個任意類型的值聚合而成的實體. 每個值為 struct 的屬性

Struct 也是值類型, 可以通過 `new()` 來創建

組成 `struct` 類型的資料稱為 `field`
- `field` 擁有自己的類型和值
- `field` 必須 unique
- `field` 的類型也可以是 struct, 或該 struct 類型

## Class in Go

Go 中沒有 `class` 的概念, 也不支持 `class` 的繼承等物件導向的概念.

Go 中 `struct` 與 `class` 都是複合結構體, 但 go 中 `struct` 的 `interface field(Nested interface)` 比物件導向具有更高的擴展性及靈活性

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

## Recursive Struct

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

go 中對 struct 進行 `&` 操作時視為對該類型進行一次 `new` 實體化操作.

所有組合字面量都是**不可尋址**的, 但其可以被**取址**

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

若 struct 全部成員都是可比較的, sturct 也可以比較; 若 struct 中存在不可比較的成員變數如 slice, map 等, 那麼 struct 則無法比較.

如果用 `==`, `!=` 進行判斷程式會直接抱錯, 可以用 `DeepEqual` 來進行深度比較

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

## Struct slice as map value

struct slice (本質是 slice)

```go
package main

import "fmt"

type student struct {
    id    int
    name  string
    score int
}

func main() {
    m := make(map[int][]student)
    
    m[101] = append(m[101], student{1, "conan", 88}, student{2, "kidd", 78})
    m[102] = append(m[101], student{1, "lan", 98}, student{2, "blame", 66})

    // 101 [{1 conan 88} {2 kidd 78}]
    // 102 [{1 conan 88} {2 kidd 78} {1 lan 98} {2 blame 66}]
    for k, v := range m {
        fmt.Println(k, v)
    }

    for k, v := range m {
        for i, data := range v {
            fmt.Println(k, i, data)
        }
    }
}
```

## Struct as function parameter

Struct 作為參數傳遞使用 pass by value (parameter & argument 在不同儲存區域)

```go
package main

import "fmt"

type student struct {
    id    int
    name  string
    score int
}

func foo(stu student) {
    stu.name = "lan"
}

func main() {
    stu := student{101, "conan", 88}
    fmt.Println(stu)  // {101 conan 88}
    foo(stu)
    fmt.Println(stu)  // {101 conan 88}
}
```

go 函數參數傳值時是以 pass by value 的方式進行, 所以函數內部無法修改傳遞給函數的原始資料結構, 修改的指示 copy 後的副本; 若傳遞的原始資料結構很大, 完整複製的開銷不小.

**故若條件允許, 應當給需要 struct instance 作為參數的函數傳遞 struct pointer**

>❗️NOTE

- struct slice 作為函數參數是 pass by address
- struct array 作為函數參數是 pass by value

## Composition Instead of Inheritance

go 中實現繼承採用組合的語法, 稱其為匿名組合

```go
type Base struct {
    name string
}

func (base *Base) Set(myname string) {
    base.name = myname
}

func (base *Base) Get() string {
    return base.name
}

type Derived struct {
    Base
    age int 
}

func (derived *Derived) Get() (nm string, ag int) {
    return derived.name, derived.age
}


func main() {
    b := &Derived{}

    b.Set("sina")
    fmt.Println(b.Get())
}
```

Base type 定義了 `Get()`, `Set()` 兩個 methods, 而 Derived type 繼承了 Base type 並 `overwrite` `Set()` method.

當 Derived object 調用 `Set()` 會 loading Base type 對應的 method; 而調用 `Get()` method 會加載 `overwrite` 後的 `Set()`

當組合類型和被組合類型包含同名的成員時會如何？

```go
type Base struct {
    name string
    age int
}

func (base *Base) Set(myname string, myage int) {
    base.name = myname
    base.age = myage
}

type Derived struct {
    Base
    name string
}

func main() {
    b := &Derived{}

    b.Set("sina", 30)
    fmt.Println("b.name =",b.name, "\tb.Base.name =", b.Base.name)
    fmt.Println("b.age =",b.age, "\tb.Base.age =", b.Base.age)
}

//b.name =        b.Base.name = sina
//b.age = 30      b.Base.age = 30

```

# Method

go 中同時有 function 和 method.

Method 就是一個包含 receiver 的函數, receiver 可以是 build-in type 或 struct type 的一個值或 pointer.

所有定義類型的 method 屬於該類型的方法集

定義一個新類型 Interger, 它和 int 一樣, 只是為它 build-in int type 增加了新的 method Less()

```go
type Integer int 

func (a Integer) Less(b Integer) bool {
    return a < b 
}

func main() {
    var a Integer = 1 

    if a.Less(2) {
        fmt.Println("less then 2")
    }   
}
```

> 可以看出 go 在自定義類型的物件中沒有 C++/Java 那種隱藏的 this pointer, 而是在定義成員方法時顯式宣告了其所屬的物件

method 語法如下：

```go
func (r ReceiverType) funcName(parameters) (results)
```

當調用 method 時, 會將 receiver 作為函數的第一個參數：

```go
funcName(r, parameters);
```

所以 receiver 是值類型還是指針類型要看 method 的作用. 若要修改物件的值, 就需要傳遞物件的 pointer.

Pointer 作為 receiver 會對實體物件內容產生操作, 而普通類型作為 receiver 僅僅是以副本進行操作, 不會對 argument object 產生操作

```go
func (a *Ingeger) Add(b Integer) {
    *a += b
}

func main() {
    var a Integer = 1 
    a.Add(3)
    fmt.Println("a =", a)     //  a = 4
}
```

若 Add method 不使用 pointer, 則 `a` 返回結果不變, **因為 go 函數參數也是基於 pass by value**

Go 沒有構造函數的概念, 通常使用一個全域函數完成：

```go
func NewRect(x, y, width, height float64) *Rect {
    return &Rect{x, y, width, height}
}   

func main() {
    rect1 := NewRect(1,2,10,20)
    fmt.Println(rect1.width)
}
```