- [Defer](#defer)
  - [Defer Methods](#defer-methods)
  - [Arguments Evaluation](#arguments-evaluation)
  - [Defer Stack](#defer-stack)
  - [Defer Use Case](#defer-use-case)

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