- [Channel](#channel)
- [Map](#map)
- [json](#json)
- [const](#const)
- [string](#string)
- [method](#method)
- [select](#select)
- [slice](#slice)
- [goroutine](#goroutine)

# Channel

1. 下面關於通道(channel)的描述正確的是（單選）?
    - 讀nil通道會觸發panic
    - 寫nil通道會觸發panic 
    - 讀關閉的通道會觸發panic
    - 寫關閉的通道會觸發panic

2. 下面的函數輸出為何(單選)?
    ```golang
    func SelectExam2() {
        c := make(chan int)
        select {
            case <-c:
            fmt.Println("readable")
            case c-< 1:
            fmt.Println("writeable")
        }
    }
    ```
    - 函數輸出readable
    - 函數輸出writeable
    - 函數什麼也不輸出, 正常返回
    - 函數什麼也不輸出, 陷入阻塞(blocking)


# Map

1. 下面的代碼存在什麼問題？
   ```golang
    var FruitColor map[string]string
    
    func AddFruit(name, color string) {
        FruitColor[name] =color
    }
    ```

# json

1. 使用標準json package操作下面的結構時, 何者的描述是正確的(單選)?
    ```golang
    type Fruit struct {
    Name string `json:&quot;name&quot;`
    Color string `json:&quot;color,omitempty&quot;`
    }
    var f = Fruit{Name:&quot;apple&quot;, Color:&quot;red&quot;}
    var s = `{&quot;Name&quot;:&quot;balana&quot;, &quot;Weight&quot;:100}`
    ```

    - 執行json.Marsha(f) 時會忽略Color字段
    - 執行json.Marsha(f)時 不會忽略Color字段
    - 執行json.Unmarsha([]byte(s), &amp;f)時 ,會忽略Color字段
    - 執行json.Unmarsha([]byte(s), &amp;f)時會出錯 ,因為Fruit類型沒有Weight字段

# const

1. 下面代碼中每個constant的值是多少?
    ```golang
    const (
        i=1<<iota
        j
        k
        l=iota
        m=1e6
    )
    ```

# string

1. 針對下面函數中的字串長度的描述正確的是(單選)?
    ```golang 
        func StringExam1() {
            var s string 
            s="台灣"
            fmt.Println(len(s))
        }
    ```
    - 字串長度表示字符個數,長度為2
    - 字串長度表示unicode編碼字節數, 長度大於2
    - 不可以針對中文字串計算長度
    - 不確定, 與運行環境有關

# method

1. 針對下列代碼的描述正確的是(單選)?
    ```golang
    type Kid struct {
        Name string
        Age int
    }

    func (k Kid) SetName(name string) {
        k.Name = name
    }

    func (k *Kid) SetAge(age int ) {
        k.Aget = age
    }
    ```
    - 編譯錯誤, 類型和類型指針不能同時作為方法的接收者
    - SetName()無法修改名字
    - SetAget()無法修改年齡
    - SetName()和SetAge()工作正常

# select

1. 針對下面的函數描述正確的是(單選)?
    ```golang
    func SelectExam5(){
        select{}
    }
    ```
    - 編譯錯誤,select語句非法
    - 運行時錯誤, 觸發panic
    - 函數陷入阻塞(blocking)
    - 函數什麼也不做直接返回

# slice

1. 下面的函數輸出為何?
    ```golang 
    func SliceRise(s []int) {
        s=append(s,0)
        for i:=range s{
        s[i]++
        }
    }

    func SlicePrint() {
        s1:=[]int{1,2}
        s2:=s1
        s2=append(s2,3)
        SliceRise(s1)
        SliceRise(s2)
        fmt.Println(s1,s2)
    }

    func main() {
        SlicePrint()
    } 
    ```

# goroutine

1. 下面的函數輸出為何?
    ```golang
    func PrintSlice() {
        s := []int{1,2,3}
        var wg sync.WaitGroup
        wg.Add(len(s))
        for _, v := range s {
            go func(){
                fmt.Println(v)
            }()
        }
        wg.Wait()
    }
    ```


