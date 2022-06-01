- [Goroutine](#goroutine)
	- [Go Concurrency Implementation](#go-concurrency-implementation)
		- [Go Thread Implementation - GMP](#go-thread-implementation---gmp)
	- [create goroutine](#create-goroutine)
	- [goroutine features](#goroutine-features)
	- [runtime package](#runtime-package)
		- [Gosched](#gosched)
		- [Goexit](#goexit)
		- [GOMAXPROCS](#gomaxprocs)
		- [Other](#other)
- [Channel](#channel)
	- [Unbuffered channel](#unbuffered-channel)
	- [Buffered channel](#buffered-channel)
		- [WaitGroup](#waitgroup)
		- [Worker Pool](#worker-pool)
	- [Close channel](#close-channel)
		- [Check if channel be closed](#check-if-channel-be-closed)
	- [One-way channel](#one-way-channel)
		- [One-way channel features](#one-way-channel-features)
		- [One-way channel as parm in the function](#one-way-channel-as-parm-in-the-function)
	- [Timer](#timer)
		- [Ticker](#ticker)
- [Producer and consumer](#producer-and-consumer)
- [Select](#select)
	- [select implement fibonacci sequence](#select-implement-fibonacci-sequence)
	- [Timeout](#timeout)
- [Lock](#lock)
	- [Deadlock](#deadlock)
	- [Mutex](#mutex)
	- [RWMutex](#rwmutex)
	- [sync.Cond](#synccond)

# Goroutine

goroutine 主要被設計用來處理 concurrency, loading 較
thread 更輕, go 在語言層面實現了 goroutine 之間的 communication

Goroutine 是 Go 中最基本的執行單位, 每個 Go 程式至少有一個 gourine: main goroutine 在程式啟動時會自動創建

> Thread

Thread 又被稱為 ｀, 是程式 execution stream 最小單位

一個標準的 thread 由 thread ID, 當前 programming counter, 暫存器集合, stack frame 所組成

Thread 是 process 的一個實體, 是被系統獨立調度和分派的基本單位; 擁有自己獨立的 stack 及與 process 共享的 heap, thread 切換一般由 OS 調度

> Coroutine

與 thread 類似, 擁有 shared heap 及 stake frame, 切換一般在程式碼中顯式控制, 其避免了 context switch 的資源消耗並兼顧了 multithreading 的優點, 簡化了 concurrency 的複雜度

> Goroutine

Coroutine 是一種 cooperative tasks 的機制, 最簡單意義上 coroutine 不是併發, 而 goroutine 則支持併發, 且可以運行在一個或多個 thread 上

## Go Concurrency Implementation

> Thread

不論語言層面使用哪種併發模型, 在 OS 層一定是以 thread 的型態存在

而 OS 根據資源訪問權限的不同, 體系架構分為 user space 和 kernel space; 
- kernel space 主要操作訪問 CPU 資源, I/O 資源, 記憶體資源等硬體資源, 為上層應用程式提供最基本的基礎資源
- user space 主要負責上層應用程式的活動空間, 其無法直接訪問資源, 必須通過 `system call`, `library function` 或 `shell` 來調用 kernel space 提供的資源

程式語言中所謂的 "thread", 往往是 user mode thread, 和 OS 本身的 kernel mode thread(KSE)不同

Thread Model 實現可以分為以下幾種方式:

- User Level Thread Model: 多個 user level thread 對應一個 kernel level thread, thread 的創建, 終止, 切換或同步必須自身完成, 可以做快速的 context switch, 缺點是無法有效利用 multi-cores CPU
- Kernel Level Thread Model: 這種模型可以直接調用 OS kernel thread, 所有 thread 的創建, 終止, 切換或同步都由 kernel 完成; 一個 user level thread 對應一個 system thread, 可以利用 multi-cores 機制, 但 context switch 需要消耗額外資源(C++)
- Two-Level Model: 介於 user level 和 kernel level 間的 thread model, 和 kernel level thread mode 類似, 一個 process 中可以對應多個 kernel level thread, 但 process 中的 thread 不與 kernel thread 一一對應; 這種 thread model 會先創建多個 kernel level thread, 然後用自己的 user level thread 去對應創建的多個 kernel level thread, 自己的 user level thread 需要本身程式做調度, kernel level thread 交給 OS 調度

Go 的 thread model 就是一種特殊的 two-level thread model(GMP Model), M 個 user level thread 對應 N 個 kernel level thread, 缺點是增加了調度器的複雜度

> Concurrency Model

Go 實現了兩種併發形式: 第一種是普遍認知的 multithreading shared memory, 即 Java 或 C++ 中的 multithreading programming; 另一種是 Go 中特有的 `CSP`(communicating sequential processes) 併發模型

CSP 併發模型是在 1970 年左右提出的概念, 屬於較新的概念, 不同於傳統 multithreading 通過 shared memory 來通訊, CSP 講究的是 **"以通訊的方式共享記憶體"**

`DO NOT COMMUNICATE BY SHARING MEMORY; INSTEAD, SHARE MEMORY BY COMMUNICATING.
`

普通的 thread concurrency model 像是 Java, C++ 或 Python, thread 間通訊都是通過 shared memory 進行, 非常典型的方式是在訪問 shared data(array, map or object) 時通過 lock 來訪問, 因此衍生出一種方便操作的資料結構 `Thread safety` , 如 Java 中 `java.util.concurrent` 中的資料結構

Go 中也實現了傳統的 multithreading concurrency model

Go 的 GSP concurrency model 是通過 `goroutine` 和 `channel` 來實現的
- `goroutine` 是 Go 中併發的執行單位
- `channel` 是 Go 中各個併發 struct(goroutine)之前的通訊機制, 有點類似 linux pipeline

### Go Thread Implementation - GMP

- `G` 指 `Goroutine`, 本質上也是一種 ｀Lightweight Process(LWP), 包含 stack, channel 等
- `M` 指 `Machine`, 一個 `M` 直接關聯一個 kernel thread, 由 OS 管理
- `P` 指 `Processor`, 代表 `M` 所提供的 context 環境, 也是處理 user level 程式碼邏輯的處理器, 主要負責銜接 `M` 和 `G` 的調度 context, 將等待執行的 `G` 和 `M` 對接

`P` 的數量是由環境變數中的 `GOMAXPROCS` 所決定, 通常來說與 cores 數對應, 例如在 4 cores sever 上會啟動 4 個 thread, `G` 會有很多個, 每個 `P` 會將 `Goroutine` 從一個就緒的 queue 中做 pop 操作, 為了減小 lock 競爭, 通常每個 `P` 會負責一個 queue

Go 中每有一個 go 語句被執行, runqueue 就在其尾段加入一個 goroutine, 一旦 context 運行 goroutine 直到調度點, 它會從 runqueue 中彈出 goroutine, 設置 stack frame 和 pointer 開始運行 goroutine

> 去除 `P`(Processor)

實際上無法去除 context, 讓 `Goroutine` 的 `runqueues` 直接掛到 `M` 上

Context 的目的是當遇到 kernel thread blocing 時, 可以直接放開其他 thread

比如讓系統調用 `sysall`, 一個 thread 肯定無法同時執行程式碼和 system call 而被 block, 此時 thread `M` 需要放棄當前的 context 環境 `P`, 以便讓其他的 `Goroutine` 被調度執行




## create goroutine

```go
package main

import (
    "fmt"
    "time"
)

func foo() {
    i := 0
    for true {
        i++
        fmt.Println("new goroutine: i = ", i)
        time.Sleep(time.Second)
    }
}

func main() {
    // create goroutine, start a new task
    go foo()

    i := 0
    for true {
        i++
        fmt.Println("main goroutine: i = ", i)
        time.Sleep(time.Second)
    }
}
```

output
```go
main goroutine: i =  1
new goroutine: i =  1
new goroutine: i =  2
main goroutine: i =  2
main goroutine: i =  3
new goroutine: i =  3
...
```
## goroutine features

main goroutine 退出後, 其他子 goroutine 也會自動退出

```go
package main

import (
    "fmt"
    "time"
)

func foo() {
    i := 0
    for true {
        i++
        fmt.Println("new goroutine: i = ", i)
        time.Sleep(time.Second)
    }
}

func main() {
    // create goroutine, start a new task
    go foo()

    time.Sleep(time.Second * 3)

    fmt.Println("main goroutine exit")
}
```

output

```go
new goroutine: i =  1
new goroutine: i =  2
new goroutine: i =  3
main goroutine exit
```

## runtime package

### Gosched

`runtime.Gosched()` 用於讓出當前 goroutine 佔用的 CPU 時間, 並讓出當前 goroutine 執行權限, 由 scheduler 安排其他等待的任務運行, 並在下次再次獲得 CPU 時間時從讓出 CPU 的位置恢復執行

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // create goroutine
    go func(s string) {
        for i := 0; i < 2; i++ {
            fmt.Println(s)
        }
    }("world")

    for i := 0; i < 2; i++ {
        runtime.Gosched()
        fmt.Println("hello")
    }
    time.Sleep(time.Second * 3)
}
```

output

```go
world
world
hello
hello
```

若無 `runtime.Gosched()` 結果如下
```go
hello
hello
world
world
```

>❗️ `runtime.Gosched()` 僅讓出一次 CPU 時間

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // create goroutine
    go func(s string) {
        for i := 0; i < 2; i++ {
            fmt.Println(s)
            time.Sleep(time.Second)
        }
    }("world")

    for i := 0; i < 2; i++ {
        runtime.Gosched()
        fmt.Println("hello")
    }
}
```

output -> main goroutine 退出後其他 goroutine 也自動結束

```go
world
hello
hello
```

### Goexit

調用 `runtime.Goexit()` 將立即中止 goroutine 執行, scheduler 確保所有 defer 調用執行

與 `return` 不同在於, `return` 是返回當前函式調用給調用者

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    go func() {
        defer fmt.Println("A.defer")
        func() {
            defer fmt.Println("B.defer")
            runtime.Goexit() // terminate this goroutine
            fmt.Println("B") // no exec
        }()
        fmt.Println("A") // no exec
    }() 

    time.Sleep(time.Second * 3)
}
```

output

```
B.defer
A.defer
```

### GOMAXPROCS

調用 `runtime.GOMAXPROCS()` 用來設置 concurrency 計算的 CPU cores 最大值並返回上一次設定值

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    runtime.GOMAXPROCS(1)  // single core

    for true {
        go fmt.Print(0)
        fmt.Print(1)
    }
}
```

output

```go
111111 ... 1000000 ... 0111 ...
```

執行 `runtime.GOMAXPROCS(1)` 時最多同時只能由一個 goroutine 被執行, 故會打印很多1. 一段時間後 go sheduler 會將其休眠並喚醒另一個 goroutine, 就開始打印出0. 打印時 goroutine 是被調度到 os thread 上

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    runtime.GOMAXPROCS(2)

    for true {
        go fmt.Print(0)
        fmt.Print(1)
    }
}
```

output

```go
111111111111111000000000000000111111111111111110000000000000000011111111100000...
```

### Other

```go
func GOROOT() string
```

GOROOT 返回 go root path. 若存在 GOROOT 環境變數則返回; 否則返回創建 go 時的 root path

```go
func Version() string
```

返回 go version

```go
func NumCPU() int
```

返回 server 邏輯 CPU cores

```go
func GC()
```

調用觸發 GC 執行



# Channel

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

### WaitGroup

`WaitGroup` 用於實現 Worker Pools, 等待一批 goroutine 執行結束, 結束前程式控制會一直 blocking, 直到這些 goroutine 全部執行完畢

```go
package main

import (  
    "fmt"
    "sync"
    "time"
)

func process(i int, wg *sync.WaitGroup) {  
    fmt.Println("started Goroutine ", i)
    time.Sleep(2 * time.Second)
    fmt.Printf("Goroutine %d ended\n", i)
    wg.Done()
}

func main() {  
    no := 3
    var wg sync.WaitGroup
    for i := 0; i < no; i++ {
        wg.Add(1)
        go process(i, &wg)
    }
    wg.Wait()
    fmt.Println("All go routines finished executing")
}
```

[WaitGroup](https://golang.org/pkg/sync/#WaitGroup) 是一個 struct 類型, 使用計數器來工作

當調用 `WaitGroup` 的 `Add` 傳遞一個 `int` 時, `WaitGroup` 的計數器會加上 `Add` 的 parm 

要減少計數器可以調用 `WaitGroup` 的 `Done()` 方法, `Wait()` 方法會 blocking 調用它的 goroutine, 直到計數器變為 0 後才會停止 blocking

> 傳遞 `&wg` 非常重要, 若沒有傳遞 `&wg` 那麼每個 goroutine 都會得到一個 `WaitGroup` 的值拷貝, 因而當他們執行結束時 `main` 函數會持續 blocking

### Worker Pool

Buffered channel 重要應用之一就是實現 [Worker Pool](https://en.wikipedia.org/wiki/Thread_pool)

一般來說 `Worker Pool` 就是一組等待任務分配的 threads, 一旦完成了所分配的任務, 這些 threads 可以繼續等待任務的分配

>Worker pool 的任務是計算所輸入每個數字的每一位和

Worker Pool 核心功能如下:
- 創建一個 goroutine pool, 監聽一個等待工作分配的輸入型 buffered channel
- 將工作添加到該 buffered channel
- 工作完成後再將結果寫入一個輸出型 buffered channel
- 從輸出型 buffered channel 讀取並打印結果

首先是創建一個 struct, 表示工作及結果

```go
type Job struct {  
    id       int
    randomno int
}
type Result struct {  
    job         Job
    sumofdigits int
}
```

所有 `Job` struct 變數都會有 `id` 及 `randomno` 兩個 field, `randomno` 用於計算每位數之和

而 `Result` struct 有一個 `job` field, 表示所對應的工作, 還有一個 `sumofdigits` field, 表示計算的結果(每位數字和)

再來分別創建用於接收工作和寫入結果的 buffered channel

```go
var jobs = make(chan Job, 10)  
var results = make(chan Result, 10)
```

Worker goroutine 會監聽 buffered channel `jobs` 中更新的工作,  一旦 worker goroutine 完成工作其結果會寫入 buffered channel `results`

如下, `digits` 函數的任務實際上就是計算整數的每一位之和, 最後返回該結果

為了模擬出 `digits` 在計算過程中花費了一段時間, 在函數中添加了兩秒的睡眠時間

```go
func digits(number int) int {  
    sum := 0
    no := number
    for no != 0 {
        digit := no % 10
        sum += digit
        no /= 10
    }
    time.Sleep(2 * time.Second)
    return sum
}
```

再寫一個創建工作 goroutine 的函數:

```go
func worker(wg *sync.WaitGroup) {  
    for job := range jobs {
        output := Result{job, digits(job.randomno)}
        results <- output
    }
    wg.Done()
}
```

上述函數創建了一個 worker 來讀取 `jobs` 的資料, 根據當前的 `job` 和 `degits` 函數的返回值創建一個 `Result` struct 變數並將結果寫進 `results` buffered channel

`worker` 函數接收了一個 `WaitGroup` 類型的 `wg` 作爲參數, 當所有 `jobs` 完成時調用 `Done()` 方法

`createWorkerPool`  函數創建了一個 goroutine 的 worker pool

```go
func createWorkerPool(noOfWorkers int) {  
    var wg sync.WaitGroup
    for i := 0; i < noOfWorkers; i++ {
        wg.Add(1)
        go worker(&wg)
    }
    wg.Wait()
    close(results)
}
```

上述函數的參數是需要創建的 worker goroutine 的數量, 在創建 goroutine 之前調用了 `wg.Add(1)` 方法, 於是 `WaitGroup` 計數器遞增

接著創建 worker goroutine 並向 `worker` 函數傳遞 `&wg`

創建完需要的 worker goroutine 後函數調用 `wg.Wait()` 等待所有的 goroutine 執行完畢, 當所有 goroutine 執行完畢後函數會關閉 `results` channel

再寫一個函數把工作分配給 worker:

```go
func allocate(noOfJobs int) {  
    for i := 0; i < noOfJobs; i++ {
        randomno := rand.Intn(999)
        job := Job{i, randomno}
        jobs <- job
    }
    close(jobs)
}
```

上述 `allocate` 函數接收創建的工作數量作為輸入參數, 生成了最大值為 999 的偽隨機數, 並使用該隨機數創建了 `Job` struct 變數

這個函數把 for loop 的計數器 `i` 作為 id, 最後把創建的 struct 變數寫入 `jobs` buffered channel

當寫入所有的 `job` 時關閉 `jobs` buffered channel

再來是創建一個函數來讀取 `results` channel和打印輸出

```go
func result(done chan bool) {  
    for result := range results {
        fmt.Printf("Job id %d, input random no %d , sum of digits %d\n", result.job.id, result.job.randomno, result.sumofdigits)
    }
    done <- true
}
```

`result` 函數讀取 `results` channel 並打印出 `job` 的 `id`, 輸入的隨機數, 該隨機數的每位數之和

`result` 函數也接受 `done` channel 作為參數, 當打印完所有結果時 `done` 會被寫入 true

最後在 `main()` 函數中調用所有函數

```go
func main() {  
    startTime := time.Now()
    noOfJobs := 100

    go allocate(noOfJobs)

    done := make(chan bool)

    go result(done)

    noOfWorkers := 10

    createWorkerPool(noOfWorkers)

    <-done

    endTime := time.Now()
    diff := endTime.Sub(startTime)

    fmt.Println("total time taken ", diff.Seconds(), "seconds")
}
```

- 通過 `endTime` 和 `startTime` 差值顯示程式運行時間
- 將 `noOfJobs` 設為 100, 並調用 `allocate` 向 `jobs` channel 新增工作
- 創建 `done` channel 並傳遞給 `result` goroutine, 其會開始打印結果並在結束時發出通知
- 通過調用 `createWorkerPool` 函數會創建一個有 10 個 goroutine 的 worker pool, `main` 函數會監聽 `done` channel 通知等待所有結果打印結束

overall:

```go
package main

import (  
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Job struct {  
    id       int
    randomno int
}
type Result struct {  
    job         Job
    sumofdigits int
}

var jobs = make(chan Job, 10)  
var results = make(chan Result, 10)

func digits(number int) int {  
    sum := 0
    no := number
    for no != 0 {
        digit := no % 10
        sum += digit
        no /= 10
    }
    time.Sleep(2 * time.Second)
    return sum
}
func worker(wg *sync.WaitGroup) {  
    for job := range jobs {
        output := Result{job, digits(job.randomno)}
        results <- output
    }
    wg.Done()
}
func createWorkerPool(noOfWorkers int) {  
    var wg sync.WaitGroup
    for i := 0; i < noOfWorkers; i++ {
        wg.Add(1)
        go worker(&wg)
    }
    wg.Wait()
    close(results)
}
func allocate(noOfJobs int) {  
    for i := 0; i < noOfJobs; i++ {
        randomno := rand.Intn(999)
        job := Job{i, randomno}
        jobs <- job
    }
    close(jobs)
}
func result(done chan bool) {  
    for result := range results {
        fmt.Printf("Job id %d, input random no %d , sum of digits %d\n", result.job.id, result.job.randomno, result.sumofdigits)
    }
    done <- true
}
func main() {  
    startTime := time.Now()
    noOfJobs := 100
    go allocate(noOfJobs)
    done := make(chan bool)
    go result(done)
    noOfWorkers := 10
    createWorkerPool(noOfWorkers)
    <-done
    endTime := time.Now()
    diff := endTime.Sub(startTime)
    fmt.Println("total time taken ", diff.Seconds(), "seconds")
}
```

output:

```go
Job id 1, input random no 636, sum of digits 15  
Job id 0, input random no 878, sum of digits 23  
Job id 9, input random no 150, sum of digits 6  
...
total time taken  20.01081009 seconds
```

隨著 worker goroutine 數量增加, 完成工作的總時間會減少, 可以透過 `main()` 函數修改 `noOfJobs` 和 `noOfWorkers` 值測試


## Close channel
- 使用 close(ch) 關閉 channel
- 讀端可以判斷 channel 是否關閉

```go
if num, ok := <-ch; ok == true {
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
			ch <- i
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
∏
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

# Lock

## Deadlock

Deadlock 指兩個或兩個以上 process 在執行過程中, 由於競爭系統資源或者由於彼此通信而造成 blocking 的現象. 若無外力干預, process 無法繼續執行, 此時稱系統產生 deadlock

單 goroutine deadlock: channel 應至少在兩個以上 goroutine 間操作避免 deadlock

channel 寫端 blocking 造成 deadlock
```go
func main() {

	ch := make(chan int)
	ch <- 789

	num := <-ch
	fmt.Println("num: ", num)
}
```

goroutine channel 訪問順序造成 deadlock: 使用 channel 一端讀(寫)要確保另一端有辦法讀(寫)

channel 讀端 blocking 造成 deadlock
```go
func main() {

	ch := make(chan int)
	num := <-ch

	fmt.Println("num: ", num)
	
	go func() {
		ch <- 789
	}()
}
```

goroutine 間多 channel 交叉訪問造成 deadlock

```go
func main() {

	ch1 := make(chan int)
	ch2 := make(chan int)
	
	go func() {
		for  {
			select {
			case num := <-ch1:
				ch2<- num
			}
		}
	}()

	for  {
		select {
		case num := <-ch2:
			ch1<- num
		}
	}
}
```

>❗️盡可能不要將 Mutex, RWMutex 與 channel 混用 - deadlock
```go
var rwMutex sync.RWMutex

func readGo(in <-chan int, i int) {
	for {
		// read mode lock
		rwMutex.RLock()

		num := <-in
		fmt.Printf("%dth read goroutine, read %d\n", i, num)

		// read mode unlock
		rwMutex.RUnlock()
	}
}

func writeGo(out chan<- int, i int) {
	for {
		// generate rand
		num := rand.Intn(1000)

		// write mode lock
		rwMutex.Lock()

		out <- num
		fmt.Printf("%dth write goroutine, write %d\n", i, num)

		time.Sleep(time.Millisecond * 300)

		// write mode lock
		rwMutex.Unlock()
	}
}

func main() {
	// rand seed
	rand.Seed(time.Now().UnixNano())

	dataCh := make(chan int)
	quitCh := make(chan bool)

	for i := 0; i < 5; i++ {
		go readGo(dataCh, i+1)
	}

	for i := 0; i < 5; i++ {
		go writeGo(dataCh, i+1)
	}

	<-quitCh
}
```


## Mutex

兩個 goroutine 共能訪問共享資料, 由於 cpu 隨機調度, 需要對共享資料訪問順序加以限定 (sync)

創建 mutex, 訪問共享資料前加鎖, 訪問結束解鎖. 在 A goroutine 加鎖期間 B goroutine 加鎖會失敗而 blocking, 直至 A goroutine 解鎖後 B 才能從 blocking 恢復執行

## RWMutex

Mutex 本質是當 goroutine 訪問時其他 goroutine 無法訪問, 這樣在資源同步上可以避免 race concondition, 但同時也降低了 concurrency performance, 由 Concurrency 變成了 Serializability

但一個不操作資料的 read 操作不存在 race condiction 問題, 所以在需要注意的是修改資料的同步. 真正的互斥應該是 RW, WW 之間, RR 之間是沒有互斥操作的必要

由此衍生出另一種鎖, 即 RWMutex

RWMutex 可以讓多個 read operation concurrecy, 但對於 write operation 完全互斥. 即當一個 goroutine 進行 write operation 時其他 goroutine 不能 read 也不能 write

**讀時共享, 寫時獨佔. write lock piority higher than read lock
**
go 中 RWMutex 由 sync.RWMutex 定義, 此類型包含兩對 methods:

> 一組是對 write operation 的鎖定及解鎖, 簡稱『寫鎖定』及『寫解鎖』

```go
func (*RWMute)Lock()
func (*RWMute)Unlock()
```

> 另一組表示對 read operation 的鎖定及解鎖, 簡稱『讀鎖定』及『讀解鎖』
```go
func (*RWMute)RLock()
func (*RWMute)RUnlock()
```

> RWMutex Implement shared data

```go
var rwMutex sync.RWMutex

// global variable to shared data
var globalVal int

func readGo(i int) {
	for {
		// read mode lock
		rwMutex.RLock()

		num := globalVal
		fmt.Printf("%dth read goroutine, read %d\n", i, num)

		// read mode unlock
		rwMutex.RUnlock()
	}
}

func writeGo(i int) {
	for {
		// generate rand
		num := rand.Intn(1000)

		// write mode lock
		rwMutex.Lock()

		globalVal = num
		fmt.Printf("%dth write goroutine, write %d\n", i, num)

		time.Sleep(time.Millisecond * 300)

		// write mode lock
		rwMutex.Unlock()
	}
}

func main() {
	// rand seed
	rand.Seed(time.Now().UnixNano())

	quitCh := make(chan bool)

	for i := 0; i < 5; i++ {
		go readGo(i + 1)
	}

	for i := 0; i < 5; i++ {
		go writeGo(i + 1)
	}

	<-quitCh
}
```

## sync.Cond

先來回顧一下之前的生產者消費者模型：

```go
package main

import (
    "fmt"
    "time"
)

func producer(out chan <- int) {
    for i:=0; i<5; i++ {
        fmt.Println("producer, produce: ", i)
        out <- i
    }
    close(out)
}

func consumer(in <- chan int) {
    for num := range in {
        fmt.Println("---consumer, consume: ", num)
    }
}

func main() {
    ch := make(chan int)
    go producer(ch)
    go consumer(ch)
    time.Sleep(5 * time.Second)
}
```

之前都是一個生產者和一個消費者，如果是多個生產者和多個消費者的情況呢？

```go

package main

import (
    "fmt"
    "math/rand"
    "time"
)

func producer(out chan <- int, idx int) {
    for i:=0; i<10; i++ {
        num := rand.Intn(800)
        fmt.Printf("%dth producer, produce: %d\n", idx, num)
        out <- num
    }
}

func consumer(in <- chan int, idx int) {
    for num := range in {
        fmt.Printf("---%dth consumer, consume: %d\n", idx, num)
    }
}

func main() {
    ch := make(chan int)
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < 5; i++ {
        go producer(ch, i + 1)
    }
    for i := 0; i < 5; i++ {
        go consumer(ch, i + 1)
    }
    time.Sleep(5 * time.Second)
}

```

Output
```go
2th goroutine, Write: 115
----3th goroutine, Read: 709
----4th goroutine, Read: 115
1th goroutine, Write: 709
3th goroutine, Write: 204
4th goroutine, Write: 711
```

如果是按照上面的程式碼寫的話, 就又會出現之前的錯誤. 即通過結果我們可以知道, 當寫入 115 時, 由於創建的是無緩衝的 channel , 應該先把這個數讀出來, 然後才可以繼續寫數據, 但是結果顯示, 讀到的是 709, 709 在下面才顯示寫入啊, 怎麼會先讀出來呢? 出現這個情況的問題在於, 當運行到 `num := <- in` 時, 已經把 709 寫進去了, 但是這個時候還沒有來得及打印, 就失去了CPU, 失去CPU之後, 緩衝區中的數據就會被覆蓋掉, 這時被 115 所覆蓋.

上面已經說過了, 解決這種錯誤有兩種方法: 用鎖或者用條件變量.

這次就用條件變量來解決一下.

首先, 強調一下. 條件變量本身不是鎖!! 但是經常與鎖結合使用!!

還有另外一個問題, 如果消費者比生產者多, 倉庫中就會出現沒有數據的情況. 我們需要不斷的通過循環來判斷倉庫隊列中是否有數據, 這樣會造成 cpu 的浪費. 反之, 如果生產者比較多, 倉庫很容易滿, 滿了就不能繼續添加數據, 也需要循環判斷倉庫滿這一事件, 同樣也會造成 cpu 的浪費.

我們希望當倉庫滿時, 生產者停止生產, 等待消費者消費; 同理, 如果倉庫空了, 我們希望消費者停下來等待生產者生產. 為了達到這個目的, 這裡就引入了條件變量. 需要注意如果倉庫隊列用是不存在以上情況的因為被填滿後就阻塞了或者中沒有數據也會阻塞.

條件變量: 條件變量的作用並不保證在同一時刻僅有一個協程線程訪問某個共享的數據資源, 而是在對應的共享數據的狀態發生變化時, 通知阻塞在某個條件上的協程線程. 條件變量不是鎖, 在並發中不能達到同步的目的, 因此條件變量總是與鎖一塊使用.

例如, 我們上面說的, 如果倉庫隊列滿了, 我們可以使用條件變量讓生產者對應的 goroutine 暫停阻塞, 但是當消費者消費了某個產品後, 倉庫就不再滿了, 應該喚醒發送通知給阻塞的生產者 goroutine 繼續生產產品.

Go標準庫中的 sync.Cond 類型代表了條件變量. 條件變量要與鎖互斥鎖或者讀寫鎖一起使用. 成員變量L代表與條件變量搭配使用的鎖.

```go
type Cond struct {
    noCopy noCopy
    L Locker
    notify notifyList
    checker copyChecker
}
```

對應的有3個常用的方法, Wait, Signal, Broadcast

1) func (c *Cond) Wait()

    該函數的作用可歸納為如下三點:

    - 釋放已掌握的互斥鎖相當於cond.L.Unlock
    - blocking & 等待條件變量滿足 (與第一步是 `atomic operation`)
    - 當被喚醒時, Wait() 函數返回時, 解除阻塞並重新獲取互斥鎖. 相當於cond.L.Lock

2) func (c *Cond) Signal()

    單發通知, 給一個正等待阻塞在該條件變量上的goroutine線程發送通知.

3) func (c *Cond) Broadcast()

    廣播通知, 給正在等待阻塞在該條件變量上的所有goroutine線程發送通知

條件變量實現生產者消費者模型

```go

package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

var cond sync.Cond  // 定義全局變量

func producer2(out chan<- int, idx int) {
    for {
        // 先加鎖
        cond.L.Lock()
        // 判斷緩衝區是否滿
        for len(out) == 3 {
            cond.Wait()
        }
        num := rand.Intn(800)
        out <- num
        fmt.Printf("%dth producer, produce: %d\n", idx, num)
        // 訪問公共區結束, 並且打印結束, 解鎖
        cond.L.Unlock()
        // 喚醒阻塞在條件變量上的 消費者
        cond.Signal()
    }
}

func consumer2(in <- chan int, idx int) {
    for {
        // 先加鎖
        cond.L.Lock()
        // 判斷緩衝區是否為 空
        for len(in) == 0 {
            cond.Wait()
        }
        num := <- in
        fmt.Printf("---%dth consumer, consume: %d\n", idx, num)
        // 訪問公共區結束後, 解鎖
        cond.L.Unlock()
        // 喚醒阻塞在條件變量上的生產者
        cond.Signal()
    }
}

func main() {
    // 設置隨機種子數
    rand.Seed(time.Now().UnixNano())

    ch := make(chan int, 3)

    cond.L = new(sync.Mutex)

    for i := 0; i < 5; i++ {
        go producer2(ch, i + 1)
    }
    for i := 0; i < 5; i++ {
        go consumer2(ch, i + 1)
    }
    time.Sleep(time.Second * 1)
}

```

1）定義 ch 作為隊列, 生產者產生數據保存至隊列中, 最多存儲3個數據, 消費者從中取出數據模擬消費

2）條件變量要與鎖一起使用, 這裡定義全局條件變量 cond, 它有一個屬性: L Locker, 是一個互斥鎖.

3）開啟5個消費者 goroutine, 開啟5個生產者 goroutine.

4）producer2 生產者, 在該方法中開啟互斥鎖, 保證數據完整性. 並且判斷隊列是否滿, 如果已滿, 調用 cond.Wait() 讓該 goroutine 阻塞. 當消費者取出數據後執行 cond.Signal(), 會喚醒該 goroutine, 繼續產生數據.

5）consumer2 消費者, 同樣開啟互斥鎖, 保證數據完整性. 判斷隊列是否為空, 如果為空, 調用 cond.Wait() 使得當前 goroutine blocking. 當生產者產生數據並添加到隊列, 執行 cond.Signal() 喚醒該 goroutine.

條件變量使用流程:

1. 創建條件變量: var cond sync.Cond
2. 指定條件變量用的鎖: cond.L = new
3. 給公共區加鎖互斥鎖: cond.L.Lock
4. 判斷是否到達阻塞條件(緩衝區滿/空) --> for循環判斷

    ```go 
    for len(ch) == cap(ch) { cond.Wait() }
    或者 for len(ch) == 0 { cond.Wait() }

    1) unlock 2)blocking 3)lock
    ```

5. 訪問公共區 --> 讀、寫數據、打印
6. 解鎖條件變量用的鎖: cond.L.Unlock
7. 喚醒阻塞在條件變量上的對端: cond.Signal cond.Broadcast