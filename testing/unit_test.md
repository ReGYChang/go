- [Unit Test](#unit-test)
- [Stress Test](#stress-test)

# Unit Test

> Unit test 重點在於發現程式設計或實現上的錯誤, 讓問題及早暴露以便於問題的定位解決

Go 中自帶有一個輕量級的測試框架 `testing` 和 `go test` 指令來實現 unit test 和效能測試, 可以基於 `testing` framework 寫針對函數的 test case, 也可以基於該框架撰寫對應的壓力測試 test case

另外建議安裝 [gotests](https://github.com/cweill/gotests) plugin 自動產生測試程式碼:

```go
go get -u -v github.com/cweill/gotests/...
```

Unit test 撰寫有以下原則:
- 檔名必須是 `_test.go` 結尾的, 這樣在執行go test的時候才會執行到相應的程式碼
- 必須 import `testing` 這個包
- 所有的測試案例函式必須是 `Test` 開頭
- 測試案例會按照原始碼中寫的順序依次執行
- 測試函式TestXxx()的參數是testing.T, 我們可以使用該型別來記錄錯誤或者是測試狀態
- 測試格式: `func TestXxx (t *testing.T)`, Xxx部分可以為任意的字母數字的組合, 但是首字母不能是小寫字母[a-z], 例如 `Testintdiv` 是錯誤的函式名
- 函式中透過呼叫 `testing.T`的 `Error`, `Errorf`, `FailNow`, `Fatal`, `FatalIf` 方法, 說明測試不通過, 呼叫 Log 方法用來記錄測試的資訊
  
下面是我們的測試案例的程式碼:

gotest.go

```go
package gotest

import (
    "errors"
)

func Division(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("除數不能為 0")
    }

    return a / b, nil
}
```

gotest_test.go

gotest_test.go

```go
package gotest

import (
    "testing"
)

func Test_Division_1(t *testing.T) {
    if i, e := Division(6, 2); i != 3 || e != nil { //try a unit test on function
        t.Error("除法函式測試沒通過") // 如果不是如預期的那麼就報錯
    } else {
        t.Log("第一個測試通過了") //記錄一些你期望記錄的資訊
    }
}

func Test_Division_2(t *testing.T) {
    t.Error("就是不通過")
}
```

在專案目錄下執行 `go test` 就會顯示以下資訊:

```go
--- FAIL: Test_Division_2 (0.00 seconds)
    gotest_test.go:16: 就是不通過
FAIL
exit status 1
FAIL    gotest    0.013s
```

從這個結果顯示測試沒有通過, 因為在第二個測試函式中我們寫死了測試不通過的程式碼 `t.Error`

那麼我們的第一個函式執行的情況怎麼樣呢?

預設情況下執行 `go test` 是不會顯示測試透過的資訊的, 我們需要帶上參數`go test -v`，這樣就會顯示如下資訊：

```go
=== RUN Test_Division_1
--- PASS: Test_Division_1 (0.00 seconds)
    gotest_test.go:11: 第一個測試通過了
=== RUN Test_Division_2
--- FAIL: Test_Division_2 (0.00 seconds)
    gotest_test.go:16: 就是不通過
FAIL
exit status 1
FAIL    gotest    0.012s
```

上面的輸出詳細的展示了這個測試的過程, 我們看到測試函式 1 `Test_Division_1` 測試通過, 而測試函式 2 `Test_Division_2` 測試失敗了, 最後得出結論測試不通過

接下來我們把測試函式 2 修改成如下程式碼：

```go
func Test_Division_2(t *testing.T) {
    if _, e := Division(6, 0); e == nil { //try a unit test on function
        t.Error("Division did not work as expected.") // 如果不是如預期的那麼就報錯
    } else {
        t.Log("one test passed.", e) //記錄一些你期望記錄的資訊
    }
}
```

然後我們執行`go test -v`，就顯示如下資訊，測試通過了:

```go
=== RUN Test_Division_1
--- PASS: Test_Division_1 (0.00 seconds)
    gotest_test.go:11: 第一個測試通過了
=== RUN Test_Division_2
--- PASS: Test_Division_2 (0.00 seconds)
    gotest_test.go:20: one test passed. 除數不能為 0

PASS
ok      gotest    0.013s
```


# Stress Test

壓力測試主要用來檢測函數(方法) 的效能, 和編寫 unit test 方式類似, 需注意以下幾點:
- 壓力測試案例必須遵循以下格式, 其中 XXX 可以是任意字母數字的組合, 但是首字母不能是小寫字母:

```go
func BenchmarkXXX(b *testing.B) { ... }
```

- `go test` default 不會執行壓力測試的函數, 如果要執行壓力測試需要帶上參數 `-test.bench`, 語法為 `-test.bench="test_name_regex"`, 如 `go test -test.bench=".*"` 表示測試全部的壓力測試函數
- 在壓力測試案例中記得在迴圈內使用 `testing.B.N` 以使測試可以正常運行
- 檔名也必須以 `_test.go` 結尾

下面是一個壓力測試檔案 `bench_test.go` :

```go
package gotest

import (
    "testing"
)

func Benchmark_Division(b *testing.B) {
    for i := 0; i < b.N; i++ { //use b.N for looping
        Division(4, 5)
    }
}

func Benchmark_TimeConsumingFunction(b *testing.B) {
    b.StopTimer() //呼叫該函式停止壓力測試的時間計數

    //做一些初始化的工作，例如讀取檔案資料，資料庫連線之類別的,
    //這樣這些時間不影響我們測試函式本身的效能

    b.StartTimer() //重新開始時間
    for i := 0; i < b.N; i++ {
        Division(4, 5)
    }
}
```

執行 `go test bench_test.go -test.bench=".*"` 可以看到以下結果：

```go
Benchmark_Division-4                            500000000          7.76 ns/op         456 B/op          14 allocs/op
Benchmark_TimeConsumingFunction-4            500000000          7.80 ns/op         224 B/op           4 allocs/op
PASS
ok      gotest    9.364s
```

上面結果顯示沒有執行任何 `TestXXX` 的 unit test function, 只執行了壓力測試函數

`Benchmark_Division` 執行了 500000000 次, 每次的執行平均時間是 7.76 ns