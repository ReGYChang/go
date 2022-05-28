- [Higher-order Function](#higher-order-function)
	- [Decorator Pattern](#decorator-pattern)
	- [Recursion](#recursion)
		- [Fibonacci](#fibonacci)
		- [Summary](#summary)

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

