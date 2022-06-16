- [Function](#function)
- [Higher-order Function](#higher-order-function)
  - [Decorator Pattern](#decorator-pattern)
  - [Recursion](#recursion)
    - [Fibonacci](#fibonacci)
    - [Summary](#summary)
- [Anonymous Function](#anonymous-function)
  - [Closure](#closure)
  - [Why Closure](#why-closure)
  - [Closure Best Practice](#closure-best-practice)
    - [Summary](#summary-1)
  - [Closure Use Cases](#closure-use-cases)
- [Defer](#defer)
  - [Defer Methods](#defer-methods)
  - [Arguments Evaluation](#arguments-evaluation)
  - [Defer Stack](#defer-stack)
  - [Defer Use Case](#defer-use-case)

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

# Higher-order Function

## Decorator Pattern

 > 高階函數係指接收其他函數做為參數傳入，或把其他函數做為結果返回的函數。可以透過實現高階函數來實現 Go 裝飾器模式

 裝飾器模式 (Decorator) 應用場景是為了某個已存在的功能模組 ( 類別或者函數) 添加一些”裝飾”功能，而又不會侵入或修改原有的功能模組。Java 中可以通過註解優雅地實現裝飾器模式，不過在 Golang 中沒有提供註解之類的語法糖，在函數式編程中可以透過高階函數來實現裝飾器模式



```go
func multiply(a, b int) int {
	return a * b
}

func useTimesM(myfunc func(a, b int) int, a, b int) int {
	startTime := time.Now()
	x := myfunc(a, b)
	fmt.Println(time.Since(startTime))
	return x
}

func main() {
	a := 2
	b := 8
	c := multiply(a, b)
	d := useTimesM(multiply, a, b)
	fmt.Printf("%d x %d = %d\n", a, b, c)
	fmt.Printf("%d x %d = %d\n", a, b, d)
}
```

## Recursion
> 遞迴函數指在函數內部調用函數自身的函數

遞迴函數須具備的條件：

- 一個問題可被拆分成多個子問題
- 原問題與子問題除了數據規模不同，解題思路都是一樣的
- 不能無限制調用，需有退出遞迴狀態條件

### Fibonacci

```go
func fab(num int) int {
	if num == 1 || num == 2 {
		return 1
	} else {
		return fab(num-1) + fab(num-2)
	}
}
```

### Summary

 回調函數：Callback，就是將一個函數 func2 作為 func1 的一個參數，func2 為回調函數，func1 為高階函數

# Anonymous Function

> 擁有函數名的函數隻能在包級語法塊中被聲明，通過函數字面量（function literal），我們可繞過這一限制，在任何表達式中表示一個函數值。函數字面量的語法和函數聲明相似，區别在於func關鍵字後沒有函數名。函數值字面量是一種表達式，它的值被成爲匿名函數（anonymous function）。函數字面量允許我們在使用時函數時，再定義它。通過這種技巧，我們可以改寫之前對strings.Map的調用：

```go
strings.Map(func(r rune) rune { return r + 1 }, "HAL-9000")
```

更爲重要的是，通過這種方式定義的函數可以訪問完整的詞法環境（lexical environment），這意味着在函數中定義的內部函數可以引用該函數的變量

## Closure

```go
package main

import (
	"fmt"
	"time"
)

func useTimes(myfunc func(int) int, arg int) {
	startTime := time.Now()
	myfunc(arg)
	fmt.Println(time.Since(startTime))
}

func squres() func() int {
	var x int
	fmt.Println("x = ", x, &x)

	return func() int {
		x++
		return x * x
	}
}

func main() {
	useTimes(func(arg int) (sum int) {
		for i := 0; i < arg; i++ {
			sum += i
		}
		return sum
	}, 100000)

	f := squres()
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
}
```

> 閉包就是一個**函數**及與**其相關的引用環境**組合的一個實體。閉包只是形式和表現上像函數，實際上不是函數。函數是一些可執行的程式碼，這些程式碼在函數被定義後就確定了，不會在 runtime 發生變化，所以一個函數只有一個 instance


> 閉包在 runtime 可以有多個 instances, 不同的引用環境和相同的函數組合可以產生不同的 instance。閉包在某些程式語言又被稱為 Lambda 表達式。函數本身不存儲任何資訊，只有當與引用環境結合後形成的閉包才具有”記憶性”。函數式編譯器靜態的概念，而閉包是 runtime 動態的概念


> **物件是附有行為的數據，閉包是附有數據的行為**

- 返回的是匿名函數，但匿名函數引用到函數外的n，因此這個匿名函數和n形成一個實體，構成閉包
- 閉包是類別，函數是操作，n是字段。函數使用到n即構成閉包
- 閉包關鍵在於分析出返回的函數所引用的變量

## Why Closure

- 強化模組化
    - 便於以簡單方式開發較小模組，提高開發速度和程式的復用性
    - 譬如計算 array 所有數字和、所有數字積或打印所有 element?
- 抽象
    - 閉包是數據與行為的組合，具有較佳的抽象能力
- 簡化程式碼
    - 函數是一階值(First-class value)，即函數可以作為另一個函數的返回值或參數，還可以作為一個變數的值
    - 函數可以嵌套定義，即在一個函數內部可以定義另一個函數
    - 允許定義匿名函數
    - 可以捕獲引用環境，並把引用環境和函數程式碼組成一個可調用的實體

## Closure Best Practice

> 所謂閉包是指有權訪問另一個函數作用域中的變量的函數，就是在一個函數內部創建另一個函數。Golang 中所有的匿名函數都是閉包

- 編寫一個函數 makeSuffix(suffix string) 可以接收一個文件後綴名(ex .jpg) 並返回一個閉包
- 調用閉包，可以傳入一個文件名，如果該文件名未指定後綴(ex .jpg)則返回文件名 .jpg，若有 .jpg 後綴則返回原文件名
- 要求使用閉包完成

```go
func makeSuffix(suffix string) func (string) string {
	return func(name string) string {
		if !strings.HasSuffix(name,suffix){
			return name + suffix
		}
		return name
	}
}
```

- strings.HasSuffix：該函數可以判斷某個字符串是否有指定後綴

```go
f2 := makeSuffix(".jpg")
fmt.Println("after file name processor: ",f2("winter"))
fmt.Println("after file name processor: ",f2("bird.jpg"))
```

### Summary

- 將匿名函數作為另一個函數的參數，回調函數
- 將匿名函數作為另一個函數的返回值，可以形成閉包結構
- 返回的匿名函數和 makeSuffix (suffix string) 的 suffix 變量組合成一個閉包，因為返回的函數引用到 suffix 這個變量
- 使用閉包好處：
    - 傳統做法需要每次都傳入後綴名，比如 .jpg
    - 閉包可以保留上次引用的某個值，所以傳入一次可以反覆使用
- **閉包不關心這些捕獲的變數或常數是否已經超出了作用域，只要閉包還在使用它，這些變數就還會存在**
- 實現工廠模式生成器

## Closure Use Cases

- 工廠模式 - 生成器
- 單例模式
- 限流模式

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