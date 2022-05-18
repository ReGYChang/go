# 匿名函數

> 擁有函數名的函數隻能在包級語法塊中被聲明，通過函數字面量（function literal），我們可繞過這一限制，在任何表達式中表示一個函數值。函數字面量的語法和函數聲明相似，區别在於func關鍵字後沒有函數名。函數值字面量是一種表達式，它的值被成爲匿名函數（anonymous function）。函數字面量允許我們在使用時函數時，再定義它。通過這種技巧，我們可以改寫之前對strings.Map的調用：

```go
strings.Map(func(r rune) rune { return r + 1 }, "HAL-9000")
```

更爲重要的是，通過這種方式定義的函數可以訪問完整的詞法環境（lexical environment），這意味着在函數中定義的內部函數可以引用該函數的變量

## 匿名函數與閉包

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

## 閉包優點

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

## 閉包 Best Practice

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

## 閉包應用場景

- 工廠模式 - 生成器
- 單例模式
- 限流模式