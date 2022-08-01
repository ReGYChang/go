- [Error](#error)
  - [Error Type](#error-type)
  - [Get More Infomation from Error](#get-more-infomation-from-error)
    - [Type Assertion for Struct - Struct Field](#type-assertion-for-struct---struct-field)
    - [Type Assertion for Struct - Call Methods](#type-assertion-for-struct---call-methods)
    - [Compare Error](#compare-error)
  - [Custom Error](#custom-error)
    - [Use New Function](#use-new-function)
    - [Use Errorf Function](#use-errorf-function)
    - [Use Struct Type and Field](#use-struct-type-and-field)
    - [Use Struct Methods](#use-struct-methods)
  - [Error Handling](#error-handling)
- [Panic](#panic)
  - [Panic Use Cases](#panic-use-cases)
  - [Defer and Panic](#defer-and-panic)
- [Recover](#recover)
  - [Pannic, Recover and Goroutine](#pannic-recover-and-goroutine)
  - [Runtime Panic](#runtime-panic)
  - [Get Stack Trace Info After Recover](#get-stack-trace-info-after-recover)

# Error

`error` 表示程式中出現了異常情況：比如說試圖打開一個文件, 而文件系統裡卻並沒有這個文件, 這種異常情況會使用 `error` 類型表示

如同其他 build-in type(int, float64), 錯誤值可以儲存在變數中或作為函數返回值等

試圖打開一個不存在的文件:

```go
package main

import (  
    "fmt"
    "os"
)

func main() {  
    f, err := os.Open("/test.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(f.Name(), "opened successfully")
}
```

調用 `os.Open()` 試圖打開路徑為 `/test.txt` 的文件

```go
func Open(name string) (file *File, err error)
```

若成功打開文件, `Open` 函數會返回一個 file handler 和一個零值(`nil`)的 error; 若打開文件時發生錯誤則會返回一個不為 `nil` 的 error

一般函數或方法返回 error 會作為最後一個值返回, 處理 error 時會將返回的 error 與 nil 比較, nil 值表示沒有 error 發生

output:

```go
open /test.txt: No such file or directory
```

## Error Type

`error` 是一個 interface type:

```go
type error interface {  
    Error() string
}
```

所有實現 `Error()` 方法的類型都可以當作是一個 error 類型, `Error()` 方法給出了錯誤的描述

`fmt.Println()` 在打印錯誤時會在內部調用 `Error() string` 方法來得到錯誤描述

## Get More Infomation from Error

有幾種方法可以透過 error 來獲取更多資訊

舉例前例查找錯誤的文件路徑, 直接解析錯誤的字符串:

```go
open /test.txt: No such file or directory
```

解析這條錯誤訊息雖然獲取了發生錯誤的文件路徑, 但這種方式很不優雅, 隨著語言版本的更新, 這條錯誤的描述隨時都有可能變化導致程式異常

Go library 提供了各種提取錯誤相關資訊的方法

### Type Assertion for Struct - Struct Field

仔細觀察 [Open()](https://pkg.go.dev/os#OpenFile) 文件, 其返回的錯誤類型是 `*PathError`

`PathError` 是 struct type, 在 library implementation 如下:

```go
type PathError struct {  
    Op   string
    Path string
    Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }
```

`*PathError` 通過宣告 `Error() string` 方法而實現 `error` interface

`Error() string` 將文件操作, 路徑及實際錯誤拼接並返回字符串, 於是得到了錯誤訊息:

```go
open /test.txt: No such file or directory
```

struct `PathError` 的 `Path` field 就有導致錯誤的文件路徑, 若修改錯誤處理邏輯:

```go
package main

import (  
    "fmt"
    "os"
)

func main() {  
    f, err := os.Open("/test.txt")
    if err, ok := err.(*os.PathError); ok {
        fmt.Println("File at path", err.Path, "failed to open")
        return
    }
    fmt.Println(f.Name(), "opened successfully")
}
```

這裡使用了 `Type Assertion` 來獲取 `error` interface underlying value, 接著使用 `err.Path` 來打印該路徑

output:

```go
File at path /test.txt failed to open
```

### Type Assertion for Struct - Call Methods

第二種獲得更多錯誤訊息的方法就是對 underlying type 進行斷言, 並通過調用該 struct type methods

standard library 中 `DNSError` struct type 定義如下:

```go
type DNSError struct {  
    ...
}

func (e *DNSError) Error() string {  
    ...
}
func (e *DNSError) Timeout() bool {  
    ... 
}
func (e *DNSError) Temporary() bool {  
    ... 
}
```

`DNSError` struct 還有 `Timeout() bool` 及 `Temporary() bool` 兩個方法, 它們返回一個 bool value, 指出該錯誤是由超時引起, 且是臨時性錯誤

透過斷言 `*DNSError` 類型並調用這些方法來確定該錯誤是臨時性錯誤還是由超時導致的:

```go
package main

import (  
    "fmt"
    "net"
)

func main() {  
    addr, err := net.LookupHost("golangbot123.com")
    if err, ok := err.(*net.DNSError); ok {
        if err.Timeout() {
            fmt.Println("operation timed out")
        } else if err.Temporary() {
            fmt.Println("temporary error")
        } else {
            fmt.Println("generic error: ", err)
        }
        return
    }
    fmt.Println(addr)
}
```

上述程式中試圖獲取 `golangbot123.com` (無效 domain name) 的 ip

通過 `*net.DNSError` 的類型斷言獲取到錯誤的 underlying value, 並分別檢查該錯誤是由超時引起還是一個臨時性的錯誤

本例中錯誤既不是臨時性錯誤, 也不是由超時引起的, 因此程式輸出為:

```go
generic error:  lookup golangbot123.com: no such host
```

### Compare Error

第三種獲取更多錯誤訊息的方式是與 `error` 類型的變數直接比較

`filepath` 中的 `Glob` 用於返回滿足 glob 模式的所有文件名, 如果模式寫的不對函數會返回一個錯誤 `ErrBadPattern`

`filepath` package 對 `ErrBadPattern` 定義如下:

```go
var ErrBadPattern = errors.New("syntax error in pattern")
```

`errors.New()` 用於創建一個新的錯誤

當模式不正確時, `Glob` 函數會返回 `ErrBadPattern`

```go
package main

import (  
    "fmt"
    "path/filepath"
)

func main() {  
    files, error := filepath.Glob("[")
    if error != nil && error == filepath.ErrBadPattern {
        fmt.Println(error)
        return
    }
    fmt.Println("matched files", files)
}
```

上述程式查詢了模式為 `[` 的文件, 為錯誤模式, 並檢查了該錯誤是否為 `nil`

將 `error` 直接與 filepath.ErrBadPattern 比較; 若條件滿足該錯誤就是由模式錯誤所導致

```go
syntax error in pattern
```

## Custom Error

### Use New Function

創建自定義錯誤最簡單的方法就是使用 `errors` package 中的 `New` 函數

`New` 函數實現:

```go
// Package errors implements functions to manipulate errors.
package errors

// New returns an error that formats as the given text.
func New(text string) error {
    return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

`errorString` 是一個 struct type, 成員只有一個 string field `s`

使用 `errorString` pointer receiver 實現 `error` interface 的 `Error() string` method

創建一個計算圓半徑的程式, 若半徑為負則返回錯誤

```go
package main

import (  
    "errors"
    "fmt"
    "math"
)

func circleArea(radius float64) (float64, error) {  
    if radius < 0 {
        return 0, errors.New("Area calculation failed, radius is less than zero")
    }
    return math.Pi * radius * radius, nil
}

func main() {  
    radius := -20.0
    area, err := circleArea(radius)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("Area of circle %0.2f", area)
}
```

上述程式中檢查半徑是否小於 0, 若小於 0 則返回 0 和對應錯誤訊息; 若大於 0 則會計算出面積並返回值為 `nil` 的 error

output:

```go
Area calculation failed, radius is less than zero
```

### Use Errorf Function

`fmt` package 中的 `Errorf` 函數會根據程式說明符, 規定錯誤的格式並返回一個符合該錯誤的 string

```go
package main

import (  
    "fmt"
    "math"
)

func circleArea(radius float64) (float64, error) {  
    if radius < 0 {
        return 0, fmt.Errorf("Area calculation failed, radius %0.2f is less than zero", radius)
    }
    return math.Pi * radius * radius, nil
}

func main() {  
    radius := -20.0
    area, err := circleArea(radius)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("Area of circle %0.2f", area)
}
```

output:

```go
Area calculation failed, radius -20.00 is less than zero
```

### Use Struct Type and Field

錯誤可以實現 `error` interface 的 struct 表示, 這種方式可以更靈活地處理錯誤

可以使用 struct field 來訪問引發錯誤的半徑

先創建一個表示錯誤的 struct type, 錯誤類型命名規範名稱以 `Error` 結尾, 於是命名為 `areaError`

```go
type areaError struct {  
    err    string
    radius float64
}
```

上述結構體類型有一個 `radius` field, 其儲存了與錯誤有關的半徑, `err` field 存儲了實際的錯誤訊息

再來是實現 `error` interface

```go
func (e *areaError) Error() string {  
    return fmt.Sprintf("radius %0.2f: %s", e.radius, e.err)
}
```

使用 pointer receiver `*areaError`, 實現了 `error` interface 及 `Error() string` method, 該方法用於打印半徑及相關錯誤描述

```go
package main

import (  
    "fmt"
    "math"
)

type areaError struct {  
    err    string
    radius float64
}

func (e *areaError) Error() string {  
    return fmt.Sprintf("radius %0.2f: %s", e.radius, e.err)
}

func circleArea(radius float64) (float64, error) {  
    if radius < 0 {
        return 0, &areaError{"radius is negative", radius}
    }
    return math.Pi * radius * radius, nil
}

func main() {  
    radius := -20.0
    area, err := circleArea(radius)
    if err != nil {
        if err, ok := err.(*areaError); ok {
            fmt.Printf("Radius %0.2f is less than zero", err.radius)
            return
        }
        fmt.Println(err)
        return
    }
    fmt.Printf("Area of rectangle1 %0.2f", area)
}
```

`circleArea` 用於計算圓的面積, 檢查半徑若小於零則通過錯誤半徑和對應錯誤訊息創建一個 `areaError` 類型的值, 然後返回 `areaError` 值的地址

在 `main` 函數檢查錯誤是否為 `nil` 並斷言 `*areaError` 類型, 若錯誤是 `*areaError` 類型就可以用 `err.radius` 來獲取錯誤半徑, 並打印出自定義錯誤訊息

>這種作法提供了更多的錯誤訊息(即導致錯誤的半徑), 通過自定義錯誤的 struct field 來定義

output:

```go
Radius -20.00 is less than zero
```

### Use Struct Methods

首先創建一個表示錯誤的 struct

```go
type areaError struct {  
    err    string //error description
    length float64 //length which caused the error
    width  float64 //width which caused the error
}
```

上述結構體類型除了有一個錯誤描述字段, 還有可能引法錯誤的寬和高

接著實現 `error` interface, 並給錯誤類型添加兩個方法, 使其提供更多的錯誤訊息

```go
func (e *areaError) Error() string {  
    return e.err
}

func (e *areaError) lengthNegative() bool {  
    return e.length < 0
}

func (e *areaError) widthNegative() bool {  
    return e.width < 0
}
```

上述程式碼中從 `Error() string` 方法中返回了關於錯誤的描述, 當 `length` 小於零時 `lengthNegative() bool` 方法返回 `true`; 而當 `width` 小於零時 `widthNegative() bool` 方法返回 `true`

這兩個方法都提供了關於錯誤的更多訊息, 其提示了計算面積失敗的原因(長度為負數或寬度為負數)

```go
func rectArea(length, width float64) (float64, error) {  
    err := ""
    if length < 0 {
        err += "length is less than zero"
    }
    if width < 0 {
        if err == "" {
            err = "width is less than zero"
        } else {
            err += ", width is less than zero"
        }
    }
    if err != "" {
        return 0, &areaError{err, length, width}
    }
    return length * width, nil
}
```

`rectArea` 函數分別檢查長或寬是否小於零, 若小於零則返回一個錯誤訊息; 否則 `rectArea` 會返回矩形的面積和一個 `nil` 的錯誤

```go
func main() {  
    length, width := -5.0, -9.0
    area, err := rectArea(length, width)
    if err != nil {
        if err, ok := err.(*areaError); ok {
            if err.lengthNegative() {
                fmt.Printf("error: length %0.2f is less than zero\n", err.length)

            }
            if err.widthNegative() {
                fmt.Printf("error: width %0.2f is less than zero\n", err.width)

            }
            return
        }
        fmt.Println(err)
        return
    }
    fmt.Println("area of rect", area)
}
```

`main` 函數中檢查錯誤是否為 `nil`, 若錯誤值不為 `nil` 接著斷言 `*areaError` 類型, 並使用 `lengthNegative()` 和 `widthNegative()` 方法來檢查錯誤原因是長度小於零還是寬度小於零

完整程式碼:

```go
package main

import "fmt"

type areaError struct {  
    err    string  //error description
    length float64 //length which caused the error
    width  float64 //width which caused the error
}

func (e *areaError) Error() string {  
    return e.err
}

func (e *areaError) lengthNegative() bool {  
    return e.length < 0
}

func (e *areaError) widthNegative() bool {  
    return e.width < 0
}

func rectArea(length, width float64) (float64, error) {  
    err := ""
    if length < 0 {
        err += "length is less than zero"
    }
    if width < 0 {
        if err == "" {
            err = "width is less than zero"
        } else {
            err += ", width is less than zero"
        }
    }
    if err != "" {
        return 0, &areaError{err, length, width}
    }
    return length * width, nil
}

func main() {  
    length, width := -5.0, -9.0
    area, err := rectArea(length, width)
    if err != nil {
        if err, ok := err.(*areaError); ok {
            if err.lengthNegative() {
                fmt.Printf("error: length %0.2f is less than zero\n", err.length)

            }
            if err.widthNegative() {
                fmt.Printf("error: width %0.2f is less than zero\n", err.width)

            }
            return
        }
        fmt.Println(err)
        return
    }
    fmt.Println("area of rect", area)
}
```

output:

```go
error: length -5.00 is less than zero  
error: width -9.00 is less than zero
```

## Error Handling

上述介紹了 Go 中的 error 機制, 再來就是介紹錯誤處理的流程

為了保持簡單, 最好在一個地方對一個錯誤採取一次處理

先看下面一個**同時處理並返回錯誤的**例子:

```go
func someFunc() (Result, error) {
 result, err := repository.Find(id)
 if err != nil {
   log.Errof(err)
   return Result{}, err
 }
  return result, nil
}
```

這段程式碼先記錄錯誤, 並將其返回給函式調用方, 如此一來便處理了兩次錯誤, 當其他團隊成員使用此函式, 便會再記錄一次錯誤, 如此一來在 system log 中會重複記載同個錯誤

想像應用程式有三層, `repository`, `interactor` 跟 `web server`:

```go
// The repository uses an external dependency orm
func getFromRepository(id int) (Result, error) {
  result := Result{ID: id}
  err := orm.entity(&result)
  if err != nil {
    return Result{}, err
  }
  return result, nil 
}
```

根據之前提到的原則, 正確方式是將錯誤返回到頂部來處理, 在 `web server` 獲取到全部的 feedback

不幸的是 Go 內置錯誤並沒有提供 stack trace, 此外錯誤是在 external dependency 上產生, 我們需要了解在我們程式碼中的哪一段導致了這個錯誤

可以使用 `github.com/pkg/errors` package, 通過新增 stack trace 及 `repository` 的錯誤訊息來重構前面的函式:

```go
import "github.com/pkg/errors"

// The repository uses an external depedency orm
func getFromRepository(id int) (Result, error) {
  result := Result{ID: id}
  err := orm.entity(&result)
  if err != nil {
    return Result{}, errors.Wrapf(err, "error getting the result with id %d", id);
  }
  return result, nil 
}
// after the error wraping the result will be
// err.Error() -> error getting the result with id 10: whatever it comes from the orm
```

`errors.Wrapf` 函式的作用為, 在不影響原始錯誤的前提下構建 stack trace 並封裝來自 orm 的錯誤

再來看看其他層如 `interactor` 是如何處理錯誤:

```go
func getInteractor(idString string) (Result, error) {
  id, err := strconv.Atoi(idString)
  if err != nil {
    return Result{}, errors.Wrapf(err, "interactor converting id to int")
  }
  return repository.getFromRepository(id) 
}
```

再來是頂層的 `web server`:

```go
r := mux.NewRouter()
r.HandleFunc("/result/{id}", ResultHandler)
func ResultHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  result, err := interactor.getInteractor(vars["id"])
  if err != nil { 
    handleError(w, err) 
  }
  fmt.Fprintf(w, result)
}
func handleError(w http.ResponseWriter, err error) { 
   w.WriteHeader(http.StatusIntervalServerError)
   log.Errorf(err)
   fmt.Fprintf(w, err.Error())
}
```

如此一來僅在頂層處理錯誤, 但會發現總是收到 500 的 http status code, 另外總是紀錄如 `result not found` 的錯誤只會為 log 增加垃圾資訊

這裡有三個解決方案:
- 提供良好的 error stack trace
- Log the error(e.g. web infrastructure layer)
- 必要時提供 contextual error info(e.g. Email 輸入格式不對)

首先創建一個錯誤型別:

```go
package errors

const(
  NoType = ErrorType(iota)
  BadRequest
  NotFound 
  // add any type you want
)

type ErrorType uint

type customError struct {
  errorType ErrorType 
  originalError error 
  contextInfo map[string]string 
}

// Error returns the mssage of a customError
func (error customError) Error() string {
   return error.originalError.Error()
}

// New creates a new customError
func (type ErrorType) New(msg string) error {
   return customError{errorType: type, originalError: errors.New(msg)}
}

// Newf creates a new customError with formatted message
func (type ErrorType) Newf(msg string, args ...interface{}) error {    
   err := fmt.Errof(msg, args...)

   return customError{errorType: type, originalError: err}
}

// Wrap creates a new wrapped error
func (type ErrorType) Wrap(err error, msg string) error {
   return type.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (type ErrorType) Wrapf(err error, msg string, args ...interface{}) error { 
   newErr := errors.Wrapf(err, msg, args..)   

   return customError{errorType: errorType, originalError: newErr}
}
```

公開 `ErrorType` 及 錯誤型別, 如此一來可以創建新錯誤並封裝現有錯誤

另外有兩件事情需要注意:
- 若沒有 export `customError` 將如何檢查錯誤型別?
- 對於外部依賴中預先存在的錯誤應如何增加或獲取錯誤內容?

可以採用 `github.com/pkg/errors` 並封裝這些 library 方法:

```go
// New creates a no type error
func New(msg string) error {
   return customError{errorType: NoType, originalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
   return customError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap wrans an error with a string
func Wrap(err error, msg string) error {
   return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
   return errors.Cause(err)
}

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
   wrappedError := errors.Wrapf(err, msg, args...)
   if customErr, ok := err.(customError); ok {
      return customError{
         errorType: customErr.errorType,
         originalError: wrappedError,
         contextInfo: customErr.contextInfo,
      }
   }

   return customError{errorType: NoType, originalError: wrappedError}
}
```

接著構建自定義方法來處理一般錯誤的 context 及型別:

```go
// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
   context := errorContext{Field: field, Message: message}
   if customErr, ok := err.(customError); ok {
      return customError{errorType: customErr.errorType, originalError: customErr.originalError, contextInfo: context}
   }

   return customError{errorType: NoType, originalError: err, contextInfo: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
   emptyContext := errorContext{}
   if customErr, ok := err.(customError); ok || customErr.contextInfo != emptyContext  {

      return map[string]string{"field": customErr.context.Field, "message": customErr.context.Message}
   }

   return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
   if customErr, ok := err.(customError); ok {
      return customErr.errorType
   }

   return NoType
}
```

構建完成後, 下面就可以應用這個新的 error library:

```go
import "github.com/our_user/our_project/errors"
// The repository uses an external depedency orm

func getFromRepository(id int) (Result, error) {
  result := Result{ID: id}
  err := orm.entity(&result)
  if err != nil {    
    msg := fmt.Sprintf("error getting the  result with id %d", id)
    switch err {
    case orm.NoResult:
        err = errors.Wrapf(err, msg);
    default: 
        err = errors.NotFound(err, msg);  
    }
    return Result{}, err
  }
  return result, nil 
}
// after the error wraping the result will be
// err.Error() -> error getting the result with id 10: whatever it comes from the orm
```

再來是 `interactor`:

```go
func getInteractor(idString string) (Result, error) {
  id, err := strconv.Atoi(idString)
  if err != nil {
    err = errors.BadRequest.Wrapf(err, "interactor converting id to int")
    err = errors.AddContext(err, "id", "wrong id format, should be an integer")

    return Result{}, err
  }
  return repository.getFromRepository(id) 
}
```

最後是 `web server`:

```go
r := mux.NewRouter()
r.HandleFunc("/result/{id}", ResultHandler)
func ResultHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  result, err := interactor.getInteractor(vars["id"])
  if err != nil { 
    handleError(w, err) 
  }
  fmt.Fprintf(w, result)
}
func handleError(w http.ResponseWriter, err error) { 
   var status int
   errorType := errors.GetType(err)
   switch errorType {
     case BadRequest: 
      status = http.StatusBadRequest
     case NotFound: 
      status = http.StatusNotFound
     default: 
      status = http.StatusInternalServerError
   }
   w.WriteHeader(status) 

   if errorType == errors.NoType {
     log.Errorf(err)
   }
   fmt.Fprintf(w,"error %s", err.Error()) 

   errorContext := errors.GetContext(err) 
   if errorContext != nil {
     fmt.Printf(w, "context %v", errorContext)
   }
}
```

> 透過 export 型別和一些值來優雅地處理錯誤, 這個設計方案的特色是可以顯式表明錯誤類型

# Panic

Go 中一般是使用錯誤來處理異常情況, 但在有些情況當程式發生異常無法繼續運行時, 會透過 `panic` 來終止程式

當函數發生 panic 時會終止運行, 並在執行完所有的 `defer` 函數後返回到函數調用方, 持續過程至當前 goroutine 所有函數都返回退出, 最後程式會打印 panic 資訊及 stack trace, 程式終止

## Panic Use Cases

> 需要注意的是, 應盡可能地使用 `error`, 而不是 `panic` 和 `recover`, 只用當程式無法繼續運行時才使用 `panic` 和 `recover`

panic 有兩個合理的 use cases:
- 發生了一個不能恢復的錯誤, 此時程式無法繼續運行(web server 無法綁定要求的 port)
- 程式上的錯誤(用 `nil` 參數調用一個只能接收合法 pointer 的方法)

Build-in panic 簽名如下:

```go
func panic(interface{})
```

當程式終止時會打印傳入 `panic` 的參數

```go
package main

import (  
    "fmt"
)

func fullName(firstName *string, lastName *string) {  
    if firstName == nil {
        panic("runtime error: first name cannot be nil")
    }
    if lastName == nil {
        panic("runtime error: last name cannot be nil")
    }
    fmt.Printf("%s %s\n", *firstName, *lastName)
    fmt.Println("returned normally from fullName")
}

func main() {  
    firstName := "Elon"
    fullName(&firstName, nil)
    fmt.Println("returned normally from main")
}
```

`fullName` 函數會檢查 `firstName` 和 `lastName` 指針是否為 `nil`; 若為 `nil` `fullName` 函數會調用含有不同錯誤訊息的 `panic`

當程式終止時, 會打印該錯誤訊息:

```go
panic: runtime error: last name cannot be nil

goroutine 1 [running]:  
main.fullName(0x1040c128, 0x0)  
    /tmp/sandbox135038844/main.go:12 +0x120
main.main()  
    /tmp/sandbox135038844/main.go:20 +0x80
```

由於程式在 `fullName` 函數第 12 行發生 panic, 因此首先打印出以下輸出:

```go
main.fullName(0x1040c128, 0x0)  
    /tmp/sandbox135038844/main.go:12 +0x120
```

接著會打印出 stack frame 的下一項, 即 `fullName(&firstName, nil)`, 因此接下來會打印出:

```go
main.main()  
    /tmp/sandbox135038844/main.go:20 +0x80
```

至此已經到達導致 panic 最上層的函數, 因此結束打印並終止程式

## Defer and Panic

若存在 defer 函數, 則會先調用它再返回函數調用方

```go
package main

import (  
    "fmt"
)

func fullName(firstName *string, lastName *string) {  
    defer fmt.Println("deferred call in fullName")
    if firstName == nil {
        panic("runtime error: first name cannot be nil")
    }
    if lastName == nil {
        panic("runtime error: last name cannot be nil")
    }
    fmt.Printf("%s %s\n", *firstName, *lastName)
    fmt.Println("returned normally from fullName")
}

func main() {  
    defer fmt.Println("deferred call in main")
    firstName := "Elon"
    fullName(&firstName, nil)
    fmt.Println("returned normally from main")
}
```

上述程式添加了 `defer` 函數的調用

```go
This program prints,

deferred call in fullName  
deferred call in main  
panic: runtime error: last name cannot be nil

goroutine 1 [running]:  
main.fullName(0x1042bf90, 0x0)  
    /tmp/sandbox060731990/main.go:13 +0x280
main.main()  
    /tmp/sandbox060731990/main.go:22 +0xc0
```

程式發生 panic 時首先執行了 defer 函數, 接著返回函數調用方, 調用方的 defer 函數繼續運行直到最上層調用函數並終止程式

# Recover

`recover` 是一個 build-in 函數, 用於重新獲得 panic goroutine 的控制

`recover` 函數簽名如下所示:

```go
func recover() interface{}
```

只有在 defer 函數內部調用 `recover` 才有用

在 defer 函數內部調用 `recover` 可以取到 `panic` 的錯誤訊息, 並停止 panic 續發事件(Panicking Sequence) 使程式運行恢復正常, 若在 defer 函數外部調用 `recover` 則不能停止 panic 續發事件

```go
package main

import (  
    "fmt"
)

func recoverName() {  
    if r := recover(); r!= nil {
        fmt.Println("recovered from ", r)
    }
}

func fullName(firstName *string, lastName *string) {  
    defer recoverName()
    if firstName == nil {
        panic("runtime error: first name cannot be nil")
    }
    if lastName == nil {
        panic("runtime error: last name cannot be nil")
    }
    fmt.Printf("%s %s\n", *firstName, *lastName)
    fmt.Println("returned normally from fullName")
}

func main() {  
    defer fmt.Println("deferred call in main")
    firstName := "Elon"
    fullName(&firstName, nil)
    fmt.Println("returned normally from main")
}
```

`recoverName()` 函數調用了 `recover()`, 返回了調用 `panic` 的 parameter

這裡打印出 `recover` 的返回值, 並在 `fullName` 函數內 defer 調用了 `recoverNames()`

當 `fullName` 發生 panic 時會調用 defer 函數 `recoverName`, 它會反過來的調用 `recover()` 來重新獲得 panic goroutine 的控制

調用 `recover` 並返回 panic 的 parameter, 因此會打印:

```go
recovered from  runtime error: last name cannot be nil
```

執行完 `recover()` 之後 panic 會停止, 程式控制返回到調用方 (這裡的 `main` 函數)

程式在發生 panic 之後會繼續正常運行, 接續打印 `returned normally from main` 及 `deferred call in main`

## Pannic, Recover and Goroutine

只有在相同的 `goroutine` 中調用 `recover` 才有用, `recover` 不能恢復一個不同 goroutine 的 panic

```go
package main

import (  
    "fmt"
    "time"
)

func recovery() {  
    if r := recover(); r != nil {
        fmt.Println("recovered:", r)
    }
}

func a() {  
    defer recovery()
    fmt.Println("Inside A")
    go b()
    time.Sleep(1 * time.Second)
}

func b() {  
    fmt.Println("Inside B")
    panic("oh! B panicked")
}

func main() {  
    a()
    fmt.Println("normally returned from main")
}
```

函數 `b()` 發生 pannic, 函數 `a()` 調用了一個 defer 函數 `recovery()` 用於恢復 panic

函數 `b()` 作為一個獨立的 goroutine 來調用

結果 panic 並不會恢復, 因為調用 `recovery` 的 goroutine 和 `b()` 中發生 panic 的 goroutine 並不相同, 無法恢復panic

output:

```go
Inside A  
Inside B  
panic: oh! B panicked

goroutine 5 [running]:  
main.b()  
    /tmp/sandbox388039916/main.go:23 +0x80
created by main.a  
    /tmp/sandbox388039916/main.go:17 +0xc0
```

若函數 `b()` 在同一個 goroutine 調用 panic 就可以恢復: 將 `go b()` 改為 `b()`

output:

```go
Inside A  
Inside B  
recovered: oh! B panicked  
normally returned from main
```

## Runtime Panic

Runtime error (array out of index) 也會導致 panic, 等價調用 build-in function `panic`, 其參數由 interface type `runtime.Error` 給出

`runtime.Error` interface 定義如下:

```go
type Error interface {  
    error
    // RuntimeError is a no-op function but
    // serves to distinguish types that are run time
    // errors from ordinary errors: a type is a
    // run time error if it has a RuntimeError method.
    RuntimeError()
}
```

而 `runtime.Error` interface 滿足 build-in interface type `error`

```go
package main

import (  
    "fmt"
)

func a() {  
    n := []int{5, 7, 4}
    fmt.Println(n[3])
    fmt.Println("normally returned from a")
}
func main() {  
    a()
    fmt.Println("normally returned from main")
}
```

上述程式試圖訪問 `n[3]`, 這是一個對 slice 的錯誤引用, 程式錯誤並拋出 panic:

```go
panic: runtime error: index out of range

goroutine 1 [running]:  
main.a()  
    /tmp/sandbox780439659/main.go:9 +0x40
main.main()  
    /tmp/sandbox780439659/main.go:13 +0x20
```

試圖恢復這個 panic:

```go
package main

import (  
    "fmt"
)

func r() {  
    if r := recover(); r != nil {
        fmt.Println("Recovered", r)
    }
}

func a() {  
    defer r()
    n := []int{5, 7, 4}
    fmt.Println(n[3])
    fmt.Println("normally returned from a")
}

func main() {  
    a()
    fmt.Println("normally returned from main")
}
```

output:

```go
Recovered runtime error: index out of range  
normally returned from main
```

## Get Stack Trace Info After Recover

當恢復 panic 時, 就釋放了它的 stack trace

使用 `Debug` package 中的 `PrintStack` 函數可以打印出 stack trace

```go
package main

import (  
    "fmt"
    "runtime/debug"
)

func r() {  
    if r := recover(); r != nil {
        fmt.Println("Recovered", r)
        debug.PrintStack()
    }
}

func a() {  
    defer r()
    n := []int{5, 7, 4}
    fmt.Println(n[3])
    fmt.Println("normally returned from a")
}

func main() {  
    a()
    fmt.Println("normally returned from main")
}
```

上述程式使用 `debug.PrintStack()` 打印 stack trace:

```go
Recovered runtime error: index out of range  
goroutine 1 [running]:  
runtime/debug.Stack(0x1042beb8, 0x2, 0x2, 0x1c)  
    /usr/local/go/src/runtime/debug/stack.go:24 +0xc0
runtime/debug.PrintStack()  
    /usr/local/go/src/runtime/debug/stack.go:16 +0x20
main.r()  
    /tmp/sandbox949178097/main.go:11 +0xe0
panic(0xf0a80, 0x17cd50)  
    /usr/local/go/src/runtime/panic.go:491 +0x2c0
main.a()  
    /tmp/sandbox949178097/main.go:18 +0x80
main.main()  
    /tmp/sandbox949178097/main.go:23 +0x20
normally returned from main
```

從輸出可以看出恢復了 panic, 打印出 `Recovered runtime error: index out of range` 及 stack trace