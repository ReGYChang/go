- [Standard Library](#standard-library)
- [regexp package](#regexp-package)

# Standard Library

像 `fmt`、`os` 等這樣具有常用功能的內建包在 Go 語言中有 150 個以上, 它們被稱為標準函式庫, 大部分(一些底層的除外)內置於 Go 本身, 完整列表可以在 [Go Walker](https://gowalker.org/search?q=gorepos) 檢視

- `unsafe`: 包含了一些打破 Go 語言“型別安全”的命令, 一般的程式中不會被使用, 可用在 C/C++ 程式的呼叫中
- `syscall`-`os`-`os/exec`:  
	- `os`: 提供給我們一個平臺無關性的作業系統功能介面, 採用類別 Unix 設計, 隱藏了不同作業系統間的差異, 讓不同的檔案系統和作業系統物件表現一致
	- `os/exec`: 提供我們執行外部作業系統命令和程式的方式,   
	- `syscall`: 底層的外部套件, 提供了作業系統底層呼叫的基本介面

透過一個 Go 程式讓Linux重啟來體現它的能力:

```go
package main
import (
	"syscall"
)

const LINUX_REBOOT_MAGIC1 uintptr = 0xfee1dead
const LINUX_REBOOT_MAGIC2 uintptr = 672274793
const LINUX_REBOOT_CMD_RESTART uintptr = 0x1234567

func main() {
	syscall.Syscall(syscall.SYS_REBOOT, 
		LINUX_REBOOT_MAGIC1, 
		LINUX_REBOOT_MAGIC2, 
		LINUX_REBOOT_CMD_RESTART)
}
```

- `archive/tar` 和 `/zip-compress`：壓縮（解壓縮）檔案功能, 
- `fmt`-`io`-`bufio`-`path/filepath`-`flag`:  
	- `fmt`: 提供了格式化輸入輸出功能
	- `io`: 提供了基本輸入輸出功能, 大多數是圍繞系統功能的封裝
	- `bufio`: 緩衝輸入輸出功能的封裝
	- `path/filepath`: 用來操作在當前系統中的目標檔名路徑 
	- `flag`: 對命令列引數的操作
- `strings`-`strconv`-`unicode`-`regexp`-`bytes`:  
	- `strings`: 提供對字串的操作
	- `strconv`: 提供將字串轉換為基礎型別的功能
	- `unicode`: 為 unicode 型的字串提供特殊的功能
	- `regexp`: 正則表示式功能
	- `bytes`: 提供對字元型分片的操作
	- `index/suffixarray`: 子字串快速查詢
- `math`-`math/cmath`-`math/big`-`math/rand`-`sort`:  
	- `math`: 基本的數學函式
	- `math/cmath`: 對複數的操作
	- `math/rand`: 偽隨機數產生
	- `sort`: 為陣列排序和自訂集合
	- `math/big`: 大數的實現和計算
- `container`-`/list-ring-heap`: 實現對集合的操作
	- `list`: 雙鏈表
	- `ring`: 環形連結串列

下面程式碼示範瞭如何遍歷一個連結串列(當 l 是 `*List`)：

```go
for e := l.Front(); e != nil; e = e.Next() {
	//do something with e.Value
}
```

- `time`-`log`:  
	- `time`: 日期和時間的基本操作,   
	- `log`: 記錄程式執行時產生的日誌, 我們將在後面的章節使用它, 
- `encoding/json`-`encoding/xml`-`text/template`:
	- `encoding/json`: 讀取並解碼和寫入並編碼 JSON 資料,   
	- `encoding/xml`: 簡單的 XML1.0 解析器
	- `text/template`:產生像 HTML 一樣的資料與文字混合的資料驅動範本  
- `net`-`net/http`-`html`:
	- `net`: 網路資料的基本操作
	- `http`: 提供了一個可擴充套件的 HTTP 伺服器和基礎客戶端, 解析 HTTP 請求和回覆
	- `html`: HTML5 解析器
- `runtime`: Go 程式執行時的互動操作, 例如垃圾回收和協程建立
- `reflect`: 實現透過程式執行時反射, 讓程式操作任意型別的變數

`exp` 套件中有許多將被編譯為新套件的實驗性的套件, 在下次穩定版本發佈的時候, 它們將成為獨立的套件, 如果前一個版本已經存在了, 它們將被作為過時的套件被回收, 然而 Go1.0 發佈的時候並沒有包含過時或者實驗性的套件

**練習 1**

使用 `container/list` 包實現一個雙向連結串列, 將 `101`、`102` 和 `103` 放入其中並打印出來, 

**練習 2**

透過使用 `unsafe` 套件中的方法來測試你電腦上一個整型變數佔用多少個位元組

# regexp package

正則表示式語法和使用的詳細資訊請參考 [維基百科](http://en.wikipedia.org/wiki/Regular_expression)

在下面的程式裡, 我們將在字串中對正則表示式模式 (pattern) 進行匹配

如果是簡單模式, 使用 `Match()` 方法便可：

```go
ok, _ := regexp.Match(pat, []byte(searchIn))
```

變數 `ok` 將返回 `true` 或者 `false`, 我們也可以使用 `MatchString()`:

```go
ok, _ := regexp.MatchString(pat, searchIn)
```

更多方法中, 必須先將正則模式透過 `Compile()` 方法返回一個 `Regexp` 物件,然後我們將掌握一些匹配, 查詢, 替換相關的功能:

```go
package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	//目標字串
	searchIn := "John: 2578.34 William: 4567.23 Steve: 5632.18"
	pat := "[0-9]+.[0-9]+" //正則

	f := func(s string) string {
		v, _ := strconv.ParseFloat(s, 32)
		return strconv.FormatFloat(v * 2, 'f', 2, 32)
	}

	if ok, _ := regexp.Match(pat, []byte(searchIn)); ok {
		fmt.Println("Match Found!")
	}

	re, _ := regexp.Compile(pat)
	//將匹配到的部分替換為"##.#"
	str := re.ReplaceAllString(searchIn, "##.#")
	fmt.Println(str)
	//引數為函式時
	str2 := re.ReplaceAllStringFunc(searchIn, f)
	fmt.Println(str2)
}
```

輸出結果：

	Match Found!
	John: ##.# William: ##.# Steve: ##.#
	John: 5156.68 William: 9134.46 Steve: 11264.36

`Compile()` 函式也可能返回一個錯誤, 在使用時忽略對錯誤的判斷是因為確信自己正則表示式是有效的, 當用戶輸入或從資料中獲取正則表示式的時候, 我們有必要去檢驗它的正確性

另外我們也可以使用 `MustCompile()` 方法, 它可以像 `Compile()` 方法一樣檢驗正則的有效性, 但是當正則不合法時程式將 `panic()`