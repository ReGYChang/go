- [Introduction](#introduction)
- [Installation](#installation)
- [Assert](#assert)
  - [Contains](#contains)
  - [DirExists](#direxists)
  - [ElementsMatch](#elementsmatch)
  - [Empty](#empty)
  - [EqualError](#equalerror)
  - [EqualValues](#equalvalues)
  - [Error](#error)
  - [ErrorAs](#erroras)
  - [ErrorIs](#erroris)
  - [Assertions Object](#assertions-object)
- [Require](#require)

# Introduction

[testify](https://github.com/stretchr/testify) 是非常流行的 Go testing package, 其提供了非常多方便的函式來做斷言及錯誤資訊輸出

`testify` 有三個核心部分組成:

- assert
- mock
- suite

# Installation

使用 `Go Modules` 創建目錄並初始化:

```go
$ mkdir -p testify && cd testify
$ go mod init github.com/darjun/go-daily-lib/testify
```

安裝 `testify`:

```go
$ go get -u github.com/stretchr/testify
```

# Assert

`assert` 提供了便捷的斷言函式, 可以極大簡化測試程式碼的編寫

總結來說, 其將之前需要`判斷 + 輸出`的方式:

```go
if got != expected {
  t.Errorf("Xxx failed expect:%d got:%d", got, expected)
}
```

簡化為一行程式碼:

```go
assert.Equal(t, got, expected, "they should be equal")
```

如此一來程式碼結構更清晰, 可讀性更高

此外, `assert` 中的函式會自動生成較清晰的錯誤描述資訊:

```go
func TestEqual(t *testing.T) {
  var a = 100
  var b = 200
  assert.Equal(t, a, b, "")
}
```

使用 `go test` 運行測試:

```go
$ go test
--- FAIL: TestEqual (0.00s)
    assert_test.go:12:
                Error Trace:
                Error:          Not equal:
                                expected: 100
                                actual  : 200
                Test:           TestEqual
FAIL
exit status 1
FAIL    github.com/darjun/go-daily-lib/testify/assert   0.107s
```

可以看到錯誤訊息更加容易閱讀

## Contains

函式簽名:

```go
func Contains(t TestingT, s, contains interface{}, msgAndArgs ...interface{}) bool
```

`Contains` 斷言 `s` 包含 `contains`, 其中 `s` 可為 string, array/slice, map; 相對 `contains` 可為 sub string, array/slice element, map key 等

## DirExists

函式簽名:

```go
func DirExists(t TestingT, path string, msgAndArgs ...interface{}) bool
```

`DirExists` 斷言路徑 `path` 為一個目錄, 若 `path` 不存在或為一個檔案時則斷言失敗

## ElementsMatch

函式簽名:

```go
func ElementsMatch(t TestingT, listA, listB interface{}, msgAndArgs ...interface{}) bool
```

`ElementMatch` 斷言 `listA` 和 `listB` 包含相同元素, 忽略元素順序

`listA/listB` 必須為 array or slice, 若有重複元素則重複元素出現的次數也必須相等

## Empty

函式簽名:

```go
func Empty(t TestingT, object interface{}, msgAndArgs ...interface{}) bool
```

`Empty` 斷言 `object` 為零值, 根據 `object` 實際型別不同而含意不同:

- pointer: nil
- integer: 0
- float: 0.0
- string: ""
- boolean: false
- slice or channel: length 0

## EqualError

函式簽名:

```go
func EqualError(t TestingT, theError error, errString string, msgAndArgs ...interface{}) bool
```

`EqualError` 斷言 `theError.Error()` 返回值與 `errString` 相等

## EqualValues

函式簽名:

```go
func EqualValues(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool
```

`EqualValues` 斷言 `expected` 與 `actual` 相等, 或可轉換為相同型別且相等

此條件較 `Equal` 更寬, 使用 `reflect.DeapEqual()` 實現:

```go
func ObjectsAreEqual(expected, actual interface{}) bool {
  if expected == nil || actual == nil {
    return expected == actual
  }

  exp, ok := expected.([]byte)
  if !ok {
    return reflect.DeepEqual(expected, actual)
  }

  act, ok := actual.([]byte)
  if !ok {
    return false
  }
  if exp == nil || act == nil {
    return exp == nil && act == nil
  }
  return bytes.Equal(exp, act)
}

func ObjectsAreEqualValues(expected, actual interface{}) bool {
    // 若 `ObjectsAreEqual` return true 則直接 return
  if ObjectsAreEqual(expected, actual) {
    return true
  }

  actualType := reflect.TypeOf(actual)
  if actualType == nil {
    return false
  }
  expectedValue := reflect.ValueOf(expected)
  if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
    // 嘗試型別轉換
    return reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
  }

  return false
}
```

下面例子基於 `int` 定義一個新的型別 `MyInt`, 其值皆為 100, `Equal()` 調用將返回 false, `EaualValues()` 會返回 true:

```go
type MyInt int

func TestEqual(t *testing.T) {
  var a = 100
  var b MyInt = 100
  assert.Equal(t, a, b, "")
  assert.EqualValues(t, a, b, "")
}
```

## Error

函式簽名:

```go
func Error(t TestingT, err error, msgAndArgs ...interface{}) bool
```

`Error` 斷言 `err` 不為 nil

## ErrorAs

函式簽名:

```go
func ErrorAs(t TestingT, err error, target interface{}, msgAndArgs ...interface{}) bool
```

`ErrorAs` 斷言 `err` 表示的 error chain 中至少有一個與 `target` 匹配, 其為 `errors.As` 的封裝

## ErrorIs

函式簽名:

```go
func ErrorIs(t TestingT, err, target error, msgAndArgs ...interface{}) bool
```

`ErrorIs` 斷言 `err` 的 error chain 中有 `target`

## Assertions Object

上面的斷言都是以 `TestingT` 為第一個參數, `testify` 提供了一種更方便的方式: 先以 `*testing.T` 創建一個 `*Assertions` 物件, `Assertions` 實現了前面所有的斷言方法, 只是不需再傳入 `TestingT` 參數:

```go
func TestEqual(t *testing.T) {
  assertions := assert.New(t)
  assertion.Equal(a, b, "")
  // ...
}
```

# Require

`require` 提供了與 `assert` 相同的 interface, 但當遇到錯誤時 `require` 會直接終止測試, 而 `assert` 會返回 false