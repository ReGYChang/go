- [CS](#cs)
  - [CPU](#cpu)
  - [Memory](#memory)
- [Syntax](#syntax)
  - [int](#int)
  - [channel](#channel)
  - [map](#map)
  - [json](#json)
  - [const](#const)
  - [string](#string)
  - [method](#method)
  - [select](#select)
  - [slice](#slice)
  - [goroutine](#goroutine)
- [Coding](#coding)


# CS

## CPU

1. 哪種效率更高? 為何?


```go
const matrixLength = 20000

func foo() {
    matrixA := createMatrix(matrixLength)
    matrixB := createMatrix(matrixLength)

    for i := 0; i < matrixLength; i++ {
        for j := 0; j < matrixLength; j++ {
            matrixA[i][j] = matrixA[i][j] + matrixA[i][j]    
        }
}

func bar() {
    matrixA := createMatrix(matrixLength)
    matrixB := createMatrix(matrixLength)

    for i := 0; i < matrixLength; i++ {
        for i := 0; j < matrixLength; j++ {
            matrixA[i][j] = matrixA[i][j] + matrixA[j][i]    
        }
}
```

> Hint: 計算機組成原理, CPU 處理器, 暫存器, 多級緩存

## Memory

1. CS 中常見幾種 Memory 分配方案?

2. Go 記憶體分配與管理方案?

3. 為什麼很多小物件會造成 GC 壓力?

# Syntax

## int

1. 以下運算結果為何? 為什麼?


```go
func main() {
    var a uint = 1
    var b uint = 2
    fmt.Println(a-b)
}
```

Hint: 補碼, 無號數

## channel

1. 下面關於通道(channel)的描述正確的是（單選）?
    - 讀nil通道會觸發panic
    - 寫nil通道會觸發panic 
    - 讀關閉的通道會觸發panic
    - 寫關閉的通道會觸發panic

2. 下面的函數輸出為何(單選)?
    ```golang
    func SelectExam2() {
        c := make(chan int)
        select {
            case <-c:
            fmt.Println("readable")
            case c-< 1:
            fmt.Println("writeable")
        }
    }
    ```
    - 函數輸出readable
    - 函數輸出writeable
    - 函數什麼也不輸出, 正常返回
    - 函數什麼也不輸出, 陷入阻塞(blocking)

3. go 如何實現 channel? Unbuffered 跟 buffered channel 差別在哪? 為什麼是 thread-safe?

## map

1. 下面的代碼存在什麼問題？

```golang
var FruitColor map[string]string

func AddFruit(name, color string) {
    FruitColor[name] =color
}
```

## json

1. 使用標準json package操作下面的結構時, 何者的描述是正確的(單選)?
    ```golang
    type Fruit struct {
    Name string `json:&quot;name&quot;`
    Color string `json:&quot;color,omitempty&quot;`
    }
    var f = Fruit{Name:&quot;apple&quot;, Color:&quot;red&quot;}
    var s = `{&quot;Name&quot;:&quot;balana&quot;, &quot;Weight&quot;:100}`
    ```

    - 執行json.Marsha(f) 時會忽略Color字段
    - 執行json.Marsha(f)時 不會忽略Color字段
    - 執行json.Unmarsha([]byte(s), &amp;f)時 ,會忽略Color字段
    - 執行json.Unmarsha([]byte(s), &amp;f)時會出錯 ,因為Fruit類型沒有Weight字段

## const

1. 下面代碼中每個constant的值是多少?
    ```golang
    const (
        i=1<<iota
        j
        k
        l=iota
        m=1e6
    )
    ```

## string

1. 針對下面函數中的字串長度的描述正確的是(單選)?
    ```golang 
        func StringExam1() {
            var s string 
            s="台灣"
            fmt.Println(len(s))
        }
    ```
    - 字串長度表示字符個數,長度為2
    - 字串長度表示unicode編碼字節數, 長度大於2
    - 不可以針對中文字串計算長度
    - 不確定, 與運行環境有關

## method

1. 針對下列代碼的描述正確的是(單選)?
    ```golang
    type Kid struct {
        Name string
        Age int
    }

    func (k Kid) SetName(name string) {
        k.Name = name
    }

    func (k *Kid) SetAge(age int ) {
        k.Aget = age
    }
    ```
    - 編譯錯誤, 類型和類型指針不能同時作為方法的接收者
    - SetName()無法修改名字
    - SetAget()無法修改年齡
    - SetName()和SetAge()工作正常

## select

1. 針對下面的函數描述正確的是(單選)?
    ```golang
    func SelectExam5(){
        select{}
    }
    ```
    - 編譯錯誤,select語句非法
    - 運行時錯誤, 觸發panic
    - 函數陷入阻塞(blocking)
    - 函數什麼也不做直接返回

## slice

1. 下面的函數輸出為何?
    ```golang 
    func SliceRise(s []int) {
        s=append(s,0)
        for i:=range s{
        s[i]++
        }
    }

    func SlicePrint() {
        s1:=[]int{1,2}
        s2:=s1
        s2=append(s2,3)
        SliceRise(s1)
        SliceRise(s2)
        fmt.Println(s1,s2)
    }

    func main() {
        SlicePrint()
    } 
    ```

## goroutine

1. 下面的函數輸出為何?
    ```golang
    func PrintSlice() {
        s := []int{1,2,3}
        var wg sync.WaitGroup
        wg.Add(len(s))
        for _, v := range s {
            go func(){
                fmt.Println(v)
            }()
        }
        wg.Wait()
    }
    ```

2. Thread 有哪幾種 model? Goroutine 實現原理與優勢為何?

# Coding

1. sync.WaitGroup 中 Wait 函數支持 WaitTimeout 功能

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    wg := sync.WaitGroup{}
    c := make(chan struct{})
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(num int, close <-chan struct{}) {
            defer wg.Done()
            <-close
            fmt.Println(num)
        }(i, c)
    }

    if WaitTimeout(&wg, time.Second*5) {
        close(c)
        fmt.Println("timeout exit")
    }
    time.Sleep(time.Second * 10)
}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
    // implement
    // 要求 sync.WaitGroup 支持 timeout
    // 如果 timeout 到了超時時間返回 true
    // 如果 WaitGroup 自然結束返回 false
}
```

2. 設計三個函數分別打印 "cat", "dog", "fish", 要求每個函數各起一個 goroutine, 並按照順序打印各 100 次

```go

```