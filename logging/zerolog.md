- [What's zerolog?](#whats-zerolog)
- [zerolog Usage](#zerolog-usage)
  - [Installation:](#installation)
  - [Contextual Logger](#contextual-logger)
  - [Multi-level Logger](#multi-level-logger)

# What's zerolog?

`zerolog` package 提供了一個專門用於 JSON 輸出的 logger, 獨特的鏈式 API 允許通過避免記憶體分配及反射來寫入 JSON(or CBOR) log

此方法由 uber 的 `zap` 開創, `zerolog` 通過更簡單的 interface 及更優秀的性能將此概念提升到更高的層次

# zerolog Usage

## Installation:

```shell
go get -u github.com/rs/zerolog/log
```

## Contextual Logger

```go
func TestContextualLogger(t *testing.T) {
 log := zerolog.New(os.Stdout)
 log.Info().Str("content", "Hello world").Int("count", 3).Msg("TestContextualLogger")


  // add context (file name / column / string)
 log = log.With().Caller().Str("foo", "bar").Logger()
 log.Info().Msg("Hello wrold")
}
```

output:

```go
{"level":"info","content":"Hello world","count":3,"message":"TestContextualLogger"}
{"level":"info","caller":"log_example_test.go:29","message":"Hello wrold"}
```

與 `zap` 相同, 都定義了強型別 field, 不同點在於 `zerolog` 採用鏈式調用

## Multi-level Logger

`zerolog` 提供了從 **Trace** 到 **Panic** 七個 level:

```go
 zerolog.SetGlobalLevel(zerolog.WarnLevel)
 log.Trace().Msg("Trace")
 log.Debug().Msg("Debug")
 log.Info().Msg("Info")
 log.Warn().Msg("Warn")
 log.Error().Msg("Error")
 log.Log().Msg("no level")
 ```

 output:

 ```go
 {"level":"warn","message":"Warn"}
 {"level":"error","message":"Error"}
 {"message":"no level"}
 ```

 >❗️NOTE
 
 1. `zerolog` 不會刪除重複 field

    ```go
    logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
    logger.Info().
        Timestamp().
        Msg("dup")
    ```

    output:

    ```go
    {"level":"info","time":1494567715,"time":1494567715,"message":"dup"}
    ```

2. 鍊式調用必須調用 `Msg` 或 `Msgf`, `Send` 才能輸出 log, `Send` 相當於調用 `Msg("")`
3. 一但調用 `Msg`, `Event` 將會被處理(放回 pool 或丟棄), 不允許二次調用

