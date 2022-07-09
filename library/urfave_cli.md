- [urfave/cli](#urfavecli)
- [Arguments](#arguments)
- [Flags](#flags)
- [Variable](#variable)
- [Placeholder Values](#placeholder-values)
- [Alternate Names](#alternate-names)
- [Values from the Environment](#values-from-the-environment)
- [Values from files](#values-from-files)
- [Piority](#piority)

# urfave/cli

`urfave/cli` 是一個 command line framework, 用於構建 cmd process

installation:

```go
$ go get -u github.com/urfave/cli/v2
```


usage:

```go
func main() {
    (&cli.App{}).Run(os.Args)
}
```

只需創建一個 `cli.App` 物件並調用 `Run()` 方法, 傳入 cmd 參數即可產生一個空白的 cli application

```go
package main

import (
  "fmt"
  "log"
  "os"

  "github.com/urfave/cli/v2"
)

func main() {
  app := &cli.App{
    Name:  "hello",
    Usage: "hello world example",
    Action: func(c *cli.Context) error {
      fmt.Println("hello world")
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

這裡設置了 Name/Usage/Action, `Name` 和 `Usage` 都顯示在 help 中, `Action` 是調用該 cmd application 時實際執行的函式, 需要的資訊可以從參數 `cli.Context` 中獲取

compile and run:

```shell
$ go build -o hello
$ ./hello
hello world
```

除此之外 cli 還生成了額外的 help 資訊:

```shell
$ ./hello --help
NAME:
   hello - hello world example

USAGE:
   hello [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

# Arguments

通過 `cli.Context` 相關方法可以獲取到 cmd 的參數資訊:

- NArg(): 返回的參數個數
- Args(): 返回的 `cli.Args` struct

```go
func main() {
  app := &cli.App{
    Name:  "arguments",
    Usage: "arguments example",
    Action: func(c *cli.Context) error {
      for i := 0; i < c.NArg(); i++ {
        fmt.Printf("%d: %s\n", i+1, c.Args().Get(i))
      }
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

output:

```go
$ go run main.go hello world
1: hello
2: world
```

# Flags

`urfave/cli` 設置和獲取 flag 都十分簡單, 在 `cli.App{}` 初始化時, 設置 field `Flags` 即可

`Flags` field 是 `[]cli.Flag` 型別, `cli.Flag` 實際上是 interface 型別

cli 為常用的型別都實現了對應的 Flag, 如 `BoolFlag`, `DurationFlag`, `StringFlag` 等, 其有一些共用 field, 如 `Name`, `Value`, `Usage`

```go
func main() {
  app := &cli.App{
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:  "lang",
        Value: "english",
        Usage: "language for the greeting",
      },
    },
    Action: func(c *cli.Context) error {
      name := "world"
      if c.NArg() > 0 {
        name = c.Args().Get(0)
      }

      if c.String("lang") == "english" {
        fmt.Println("hello", name)
      } else {
        fmt.Println("你好", name)
      }
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

上面是一個打招呼的 cmd application, 可通過選項 `lang` 指定語言, 默認為英文

若有參數則使用第一個參數作為人名, 否則使用 world

注意 flag 參數是透過 `c.Type(name)` 來獲取的, Type 為 flag 型別, name 為 flag name

```shell
$ go build -o flags

# default
$ ./flags
hello world

# setup lang flag
$ ./flags --lang chinese
你好 world

# input argument as people name
$ ./flags --lang chinese dj
你好 dj
```

可以通過 `--help` 來查看選項:

```go
$ ./flags --help
NAME:
   flags - A new cli application

USAGE:
   flags [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --lang value  language for the greeting (default: "english")
   --help, -h    show help (default: false)
```

# Variable

除了通過 `c.Type(name)` 來獲取 flag 參數, 也可以通過將選項存到預先聲明的變數中, 只需要設置 `Destination` field 為 var address 即可:

```go
func main() {
  var language string

  app := &cli.App{
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:        "lang",
        Value:       "english",
        Usage:       "language for the greeting",
        Destination: &language,
      },
    },
    Action: func(c *cli.Context) error {
      name := "world"
      if c.NArg() > 0 {
        name = c.Args().Get(0)
      }

      if language == "english" {
        fmt.Println("hello", name)
      } else {
        fmt.Println("你好", name)
      }
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

效果同上

# Placeholder Values

cli 可以在 `Usage` field 中為 flag 設置佔位符, 通過反引號標示

佔位符有助於生成易於理解的 help 資訊:

```go
func main() {
  app := & cli.App{
    Flags : []cli.Flag {
      &cli.StringFlag{
        Name:"config",
        Usage: "Load configuration from `FILE`",
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

設置佔位符後在 help 中, 佔位值會顯示在對應選項後面, 對短選項也是有效:

```shell
$ go build -o placeholder
$ ./placeholder --help
NAME:
   placeholder - A new cli application

USAGE:
   placeholder [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE  Load configuration from FILE
   --help, -h     show help (default: false)
```

# Alternate Names

Flag 可以設置多個 alias, 只要設置對應的 flag `Aliases` field 即可:

```go
func main() {
  app := &cli.App{
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:    "lang",
        Aliases: []string{"language", "l"},
        Value:   "english",
        Usage:   "language for the greeting",
      },
    },
    Action: func(c *cli.Context) error {
      name := "world"
      if c.NArg() > 0 {
        name = c.Args().Get(0)
      }

      if c.String("lang") == "english" {
        fmt.Println("hello", name)
      } else {
        fmt.Println("你好", name)
      }
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

這裡使用 `--lang chinese`, `--language chiness` 和 `-l chiness` 效果相同, 若通過不同 alias 指定相同 flag 則會報錯:

```shell
$ go build -o aliase
$ ./aliase --lang chinese
你好 world
$ ./aliase --language chinese
你好 world
$ ./aliase -l chinese
你好 world
$ ./aliase -l chinese --lang chinese
Cannot use two forms of the same flag: l lang
```

# Values from the Environment

除了通過執行程式時手動輸入特定參數之外, 也可以通過讀取指定環境變數的方式輸入參數, 只需將環境變數的名子設置到 flag `EnvVars` field 即可

可以同時指定多個環境變數名, cli 會依次查詢, 第一個有值的環境變數會被使用

```go
func main() {
  app := &cli.App{
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:    "lang",
        Value:   "english",
        Usage:   "language for the greeting",
        EnvVars: []string{"APP_LANG", "SYSTEM_LANG"},
      },
    },
    Action: func(c *cli.Context) error {
      if c.String("lang") == "english" {
        fmt.Println("hello")
      } else {
        fmt.Println("你好")
      }
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

compile & run:

```shell
$ go build -o env
$ APP_LANG=chinese ./env
你好
```

# Values from files

cli 也支持從文件中讀取 flag value, 設置 flag `FilePath` field 作為文件路徑:

```go
func main() {
  app := &cli.App{
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:     "lang",
        Value:    "english",
        Usage:    "language for the greeting",
        FilePath: "./lang.txt",
      },
    },
    Action: func(c *cli.Context) error {
      if c.String("lang") == "english" {
        fmt.Println("hello")
      } else {
        fmt.Println("你好")
      }
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

在 `main.go` 同層目錄下創建一個 `lang.txt`, 輸入內容 chinese

compile & run:

```shell
$ go build -o file
$ ./file
你好
```

> cli 還支持從 YAML/JSON/TOML 等配置文件中讀取 flag value

# Piority

上面幾種設置 flag 的方式, 若同時有多種方式生效, 則按照以下優先順序由高到低設置:

- 使用者指定參數
- 環境變數
- 配置文件
- default