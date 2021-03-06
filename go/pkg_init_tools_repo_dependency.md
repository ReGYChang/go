- [Package Management](#package-management)
  - [GOPATH](#gopath)
  - [Go Vendor](#go-vendor)
  - [Go Modules](#go-modules)

# Package Management

## GOPATH

- bin：存放編譯後生成的二進制可執行文件
- pkg：存放編譯後生成的 .a 文件
- src：存放項目 source code，可以是自己寫的 or go get

> 將你的 package or 別人的 package 全部放在 **$GOPATH/src** 進行管理的方式稱之為 GOPATH 模式。此模式下，使用 go install 生成的可執行文件會放在 **\$GOPATH/bin** 下

問題：

- 無法在項目中使用指定版本的 package
- 其他人運行你開發的程式無法保證他下載的包版本是你期望的版本，當對方使用不同版本的包可能導致程式無法正常運行
- 本地中一個包只能保留一個版本，意味著在本地開發的所有項目都只能用同一個版本的包

## Go Vendor

> 在每個項目下都創建一個 vendor 目錄，每個項目所需的依賴都只會下載到自己的 vendor 目錄下，項目之間的依賴包互不影響。在編譯時，v1.5 Golang 在設置開啟 **GO15VENDOREXPERIMENT=1** 會提升 vendor 目錄依賴包搜索路徑的優先級 (相較於 GOPATH)

搜尋包的優先級順序由高到低：

- 當前包下的 vendor 目錄
- 向上級目錄查找，直到找到 src 下的 vendor 目錄
- 在 GOROOT 目錄下查找
- 在 GOPATH 下查找依賴包

問題：

- 若多個項目用到了同一個包的同一個版本，這個包會存在於機器上的不同目錄下，對硬碟空間是一種浪費，且無法對第三方包進行集中式管理
- 若要開源項目，需要將所有依賴包悉數上傳，別人使用的時候除了你的項目源碼外，還有所有的依賴包全部都要下載，才能保證別人使用時不會因為版本問題導致項目不能如預期運行

## Go Modules

從 v1.11 開始，go env 多了一個環境變數：GO111MODULE，通過他可以開啟或關閉 go mod 模式

他有三個可選項：

- off：禁用模組支持，編譯時會從 GOPATH & vendor 中查找依賴包
- on：啟用模組支持，編譯時會忽略 GOPATH & vendor ，只根據 go.mod 下載依賴
- auto：當項目在 $GOPATH/src 外且項目根目錄有 go.mod 文件時自動開啟模塊支持

在一個空資料夾中初始化一個 Module

```go
$ go mod init example
go: creating new go.mod: module example
```

此時在當前資料夾下生成了 go.mod, 此文件記錄當前 module 的 module name 以及所有的 dependency packages 版本
