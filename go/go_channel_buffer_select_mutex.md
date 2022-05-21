- [channel](#channel)
  - [Introduction](#introduction)
  - [Unbuffered channel](#unbuffered-channel)
  - [Buffered channel](#buffered-channel)
  - [Close channel](#close-channel)
    - [Check if channel be closed](#check-if-channel-be-closed)
  - [One-way channel](#one-way-channel)
    - [One-way channel features](#one-way-channel-features)
    - [One-way channel as parm in the function](#one-way-channel-as-parm-in-the-function)
  - [Timer](#timer)
    - [三種定時器方法](#三種定時器方法)
    - [Ticker](#ticker)
- [Producer and consumer](#producer-and-consumer)
- [Select](#select)
  - [select implement fibonacci sequence](#select-implement-fibonacci-sequence)
  - [Timeout](#timeout)

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
- **讀寫要求同步**

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
- 具備 **async communication** 能力, 不需同時操作緩衝區

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

## One-way channel

```go
// default channel 是雙向的 -> make(chan type)
var ch chan int
ch = make(chan int)

// single-side write channel
// can not read operation
var sendCh chan<- int
sendCh = make(chan<- int)

// single-side read channel
// can not write operation
var recvCh <-chan int
recvCh = make(<-chan int)

// convert
// 雙向 channel 可以隱式轉換為任意一種單向 channel
sendCh = ch

// 單向 channel 不能轉換為雙向 channel
ch = sendCh / recvCh -> error

// call by reference
```

### One-way channel features

```go
func main() {
	// 雙向 channel
	ch := make(chan int)

	var sendCh chan<- int = ch
	sendCh<- 789
    // error
	//num := <-sendCh

	var recvCh <-chan int = ch
	num := <-recvCh
	fmt.Println(num)
    // error
    //recvCh<- 789

	// 反向賦值 - error
	//var ch2 chan int = sendCh
}
```

### One-way channel as parm in the function
```go
func send(out chan<- int) {
	out <- 89
	close(out)
}

func recv(in <-chan int) {
	n := <-in
	fmt.Print("Read int: ", n)
}

func main() {
	ch := make(chan int)
	go func() {
		// 雙向 channel 轉換為 write channel
		send(ch)
	}()

	recv(ch)
}
```

## Timer

time.Timer 是一個定時器, 代表未來的一個單一事件, 可以設定 timer 要等待多長時間

```go
type Timer struct {
    C <-chan Time
    r runtimeTimer
}
```

其提供一個 channel, 在設定時間到達前沒有資料寫入 timer.C 會一直 blocking 直到設定時間到, 系統會自動向 timer.C 這個 channel 中寫入當前時間, blocking 即被解除

### 三種定時器方法
```go
func main() {

	// 1. sleep
	time.Sleep(time.Second)

	// 2. Timer.C
	fmt.Println("OS Time now: ", time.Now())

	// create timer
	myTimer := time.NewTimer(time.Second * 2)

	// reset timer
	myTimer.Reset(time.Second * 10)

	// stop timer
    // <-myTimer.C 會 blocking
	//myTimer.Stop()

	nowTime := <-myTimer.C
	fmt.Println("OS Time now: ", nowTime)

	// 3. time.After()
	nowTime2 := <-time.After(time.Second * 2)
	fmt.Println("OS Time now: ", nowTime2)

}
```

### Ticker

週期定時器
```go
func main() {

	fmt.Println("now: ", time.Now())
	myTicker := time.NewTicker(time.Second)

	go func() {
		for {
			nowTime := <-myTicker.C
			fmt.Println("nowTime: ", nowTime)
		}
	}()

	for {

	}

    // now:  2022-05-21 15:38:28.100423 +0800 CST m=+0.000159209
    // nowTime:  2022-05-21 15:38:29.102015 +0800 CST m=+1.001732334
    // nowTime:  2022-05-21 15:38:30.102456 +0800 CST m=+2.002153334
    // nowTime:  2022-05-21 15:38:31.102266 +0800 CST m=+3.001944751
    // nowTime:  2022-05-21 15:38:32.102818 +0800 CST m=+4.002477001
    // nowTime:  2022-05-21 15:38:33.102815 +0800 CST m=+5.002454251
    // nowTime:  2022-05-21 15:38:34.102833 +0800 CST m=+6.002453251

}
```

# Producer and consumer
- Producer: Write Data
- Consumer: Read Data
- Buffer 
  - Decouple (降低 producer & consumer 耦合)
  - Increase concurrency (生產者消費者數量不對等時可以保持正常通訊)
  - Cache (生產者消費者處理速度不一致時可以暫存資料)

```go
func producer(out chan<- int) {
	for i := 0; i < 10; i++ {
		fmt.Println("Producer write the data: ", i*i)
		out <- i * i
	}
	close(out)
}

func consumer(in <-chan int) {
	for num := range in {
		fmt.Println("Consumer read the data: ", num)
		time.Sleep(time.Second)
	}
}

func main() {
	ch := make(chan int, 5)

	go producer(ch)
	consumer(ch)
}
```


# Select
- go 提供關鍵字 select
- 通過 select 可以監聽 channel 上資料流動
- select 用法與 switch 類似, 每個選擇條件由於 case 描述
- 較 switch 有較多限制, 其中最大限制為每個 case 中必須是一個 IO 操作

```go
func main() {

	ch := make(chan int)

	select {
	case <-ch:
		// 若 ch 成功讀到資料則執行該 case 處理
	case ch<- 1:
		// 若成功向 ch 寫入資料則進行該 case 處理
	}
    default:
        // 若上面都沒有成功, 則進入 default 處理
}
```

- select 會按順序從頭到尾評估每一個發送及接收語句
- 若其中的任意一條可以繼續執行(未被 blocking), 那就從那些可執行語句中任意選擇一條處理
- 若沒有任意一條可以執行(即所有 channel blocking)有以下兩種情況
  - 若給出 default 就會執行 default 語句, 同時程式執行會從 select 語句後的語句中恢復
  - 若沒有 default 即 blocking select, 直到至少有一條可以執行
- 為了避免 busy polling, 一般不使用 default

```go
func main() {

	// data communication channel
	dataCh := make(chan int)

	// check if exit channel
	quitCh := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			dataCh<- i
			time.Sleep(time.Second)
		}
		close(dataCh)
		quitCh<- true
		runtime.Goexit()
	}()

	for {
		select {
		case num := <-dataCh:
			fmt.Println("data is: ", num)
		case <-quitCh:
			// break: 跳出 select code block
			// terminate process
			return
		}
		fmt.Println("----------------------")
	}

}
```

>❗️注意事項
- 監聽的 case 沒有滿足條件則 blocking
- 監聽的 case 有多個滿足監聽條件則任選一個執行
- 可使用 default 處理所有 case 都不滿足監聽條件的狀況(通常不用-> busy polling)
- select 自身不帶有循環機制, 需要借助外層 for loop
- break 只能跳出 select, 類似於 switch 用法

## select implement fibonacci sequence
```go
func fibonacci(ch <-chan int, quit <-chan bool) {
	for {
		select {
		case num := <-ch:
			fmt.Println("Read data: ", num)
		case <-quit:
			//return
			runtime.Goexit()
		}
	}
}

func main() {

	dataCh := make(chan int)
	quitCh := make(chan bool)

	go fibonacci(dataCh, quitCh)

	x, y := 1, 1
	for i := 0; i < 20; i++ {
		dataCh <- x
		x, y = y, x+y
	}
	quitCh <- true
}
```

## Timeout

> 有時會出現 goroutine blocking 的情況, 可以利用 select 來處理 timeout

```go
func main() {

	c := make(chan int)
	q := make(chan bool)

	go func() {
		for {
			select {
			case v := <-c:
				fmt.Println(v)
			case <-time.After(5 * time.Second):
				fmt.Println("Timeout")
				q <- true
				//runtime.Goexit()
				//return
				goto label
			}
		}
	label:
	}()

    // trigger select case 1
	for i := 0; i < 2; i++ {
		c <- i
		time.Sleep(2 * time.Second)
	}

	<-q
}
```