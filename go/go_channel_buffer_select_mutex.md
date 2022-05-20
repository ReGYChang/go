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

特性：
- channel capacity = 0, len = 0
- 需應用在至少 2 個 goroutine, 一個 read 一個 write, 否則會造成 deadlock
- 讀寫要求同步

```go
package main

import "fmt"

func main() {
	// unbuffered channel
	ch := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("goroutine i = ", i)
			ch <- i
		}
	}()

	for i := 0; i < 5; i++ {
		num := <-ch
		// io operation -> 較耗時所以 goroutine 又搶到 cpu 並寫入 channel
		fmt.Println("main read i = ", num)
	}
}
```
```
output:

goroutine i =  0
goroutine i =  1
main read i =  0
main read i =  1
goroutine i =  2
goroutine i =  3
main read i =  2
main read i =  3
goroutine i =  4
main read i =  4
```

## Buffered channel
特性：
- channel capacity > 0
- 緩衝區可以進行資料存儲, 至容量上限則 blocking
- 具備 async communication 能力, 不需同時操作緩衝區

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// 存滿 3 個 element 之前不會 blocking
	ch := make(chan int, 3)
	fmt.Println("len= ", len(ch), "cap= ", cap(ch))

	go func() {
		for i := 0; i < 8; i++ {
			ch <- i
			fmt.Println("goroutine: i= ", i, "len= ", len(ch), "cap= ", cap(ch))
		}
	}()

	time.Sleep(time.Second * 3)

	for i := 0; i < 8; i++ {
		num := <-ch
		fmt.Println("main read: ", num)
	}
}
```

## Close channel
- 使用 close(ch) 關閉 channel
- 讀端可以判斷 channel 是否關閉

```go
if num, of := <-ch; ok == true {
    //  channel 關閉 return false, num == nil
    //  channel 未關閉 return true, num == <-ch
}
```

### Check if channel be closed
ok:
```go
package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		for i := 0; i < 8; i++ {
			ch<- i
		}
		// write side close channel
		close(ch)
	}()

	for {
		if num, ok := <-ch; ok == true {
			fmt.Println("Read the data: ", num)
		} else {
			break
		}
	}
}
```

range:

```go
package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		for i := 0; i < 8; i++ {
			ch<- i
		}
		// write side close channel
		close(ch)
	}()

	for num := range ch {
        fmt.Println("Read num: ", num)
	}
}
```

> Summary
- 資料沒發完不應該關閉
- 已經關閉的 channel 不能再向其寫資料 -> panic
- 寫端已經關閉 channel 依然可以從中讀取資料, int 默認 0
