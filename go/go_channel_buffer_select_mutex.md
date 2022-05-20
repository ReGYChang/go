# channel

## Introduction

channel 為 go 中的一種資料類型，主要用來解決 goroutine 的同步問題以及資料共享(message passing)的問題

channel 運行在相同的記憶體位址，因此共享記憶體必須做好同步。**goroutine 通過 communication 來 shared memory，而不是使用 shared memory 來 communication** (CMT)

引用類型 channel 可用於多個 goroutine 間通訊，其內部實現了同步以確保 concurrency safe

和 map 類似，channel 也是一個對應 **make** 創建的底層資料結構的指針

```go
make(chan Type) // 等價於 make(chan Type, 0), 表示無緩衝 channel
make(chan Type, capacity) // 表示有緩衝 channel
```

channel 有兩端：
- write - chan<-
- read  - <-chan
- channel 必須同時滿足 write & read 才能傳輸, 否則 blocking
```go
package main

import (
	"fmt"
	"time"
)

// global define channel to sync data
var channel = make(chan int)

func printer(s string) {
	for _, ch := range s {
		fmt.Printf("%c", ch)
		time.Sleep(300 * time.Millisecond)
	}
}

func person1() {
	printer("hello")
	channel <- 8
}

func person2() {
	<-channel
	printer("world")
}

func main() {
	go person1()
	go person2()
	for {

	}
}

```

## Unbuffered channel

Unbuffered channel 是指接收前沒有能力保存任何值的 channel

要求發送端 goroutine 及 接收端 goroutine 同時準備好才能完成發送及接收操作, 否則 channel 會導致先執行發送或接收操作的 goroutine blocking

block: 由於某種原因資料未到達, 當前 goroutine 持續處於等待狀態直到滿足條件

sync: 在兩個或多個 goroutine 之間保持資料內容一致性