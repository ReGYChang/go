- [Interface](#interface)
  - [Nil Interface](#nil-interface)
  - [Interface Implement Open Closed Principle](#interface-implement-open-closed-principle)
  - [Embedding Interface](#embedding-interface)

# Interface

Interface 是一組抽象方法集合(未實現方法/僅包含方法名參數返回值的方法), 如果 implement interface 中所有 methods, 即該類型物件實現了此 interface

interface 定義語法：

```go
type interfaceName interface {  
    //methods 
}  
```

```go
package main

import "fmt"

type Human struct {
    name string
    age int
    phone string
}

type Student struct {
    Human //anonymous field
    school string
    loan float32
}

type Employee struct {
    Human //anonymous field
    company string
    money float32
}

//Human implement SayHi method
func (h Human) SayHi() {
    fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

//Human implement SayHi method
func (h Human) Sing(lyrics string) {
    fmt.Println("La la la la...", lyrics)
}

//Employee overwrite Human SayHi method
func (e Employee) SayHi() {
    fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
        e.company, e.phone)
    }

// Interface Men 被 Human,Student和Employee implement
// 因為這三個 struct 都 implement 下面兩個 methods
type Men interface {
    SayHi()
    Sing(lyrics string)
}

func main() {
    mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
    paul := Student{Human{"Paul", 26, "111-222-XXX"}, "Harvard", 100}
    sam := Employee{Human{"Sam", 36, "444-222-XXX"}, "Golang Inc.", 1000}
    tom := Employee{Human{"Tom", 37, "222-444-XXX"}, "Things Ltd.", 5000}

    //定義 Men type 的變數 i
    var i Men

    //i 能承載 Student
    i = mike
    fmt.Println("This is Mike, a Student:")
    i.SayHi()
    i.Sing("November rain")

    //i 也能承載 Employee
    i = tom
    fmt.Println("This is tom, an Employee:")
    i.SayHi()
    i.Sing("Born to be wild")

    //宣告 slice Men
    fmt.Println("Let's use a slice of Men and see what happens")
    x := make([]Men, 3)
    //三個不同 type element 都 implement Men interface
    x[0], x[1], x[2] = paul, sam, mike

    for _, value := range x{
        value.SayHi()
    }
}
```

## Nil Interface

Nil interface (interface{}) 不包含任何 method, 所有的類型都實現了 nil interface. 

Nil interface 在需要承載任意類型的值時相當有用, 類似 C 中的 `void*`

```go
// define a nil interface
var a interface{}
var i int = 5
s := "Hello world"
// a 可以承載任意類型值
a = i
a = s
```

那如何知道這個 interface 變數實際保存了哪種類型的物件？目前常用的兩種方法：
- switch
- Comma-ok

Switch check
```go
type Element interface{}
type List [] Element

type Person struct {
    name string
    age int 
}

//print
func (p Person) String() string {
    return "(name: " + p.name + " - age: "+strconv.Itoa(p.age)+ " years)"
}

func main() {
    list := make(List, 3)
    list[0] = 1 //an int 
    list[1] = "Hello" //a string
    list[2] = Person{"Dennis", 70} 

    for index, element := range list{
        switch value := element.(type) {
            case int:
                fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
            case string:
                fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
            case Person:
                fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
            default:
                fmt.Println("list[%d] is of a different type", index)
        }   
    }   
}
```

Comma-ok check

```go
func main() {
    list := make(List, 3)
    list[0] = 1 // an int
    list[1] = "Hello" // a string
    list[2] = Person{"Dennis", 70}

    for index, element := range list {
        if value, ok := element.(int); ok {
            fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
        } else if value, ok := element.(string); ok {
            fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
        } else if value, ok := element.(Person); ok {
            fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
        } else {
            fmt.Printf("list[%d] is of a different type\n", index)
        }
    }
}
```

## Interface Implement Open Closed Principle

OCP 描述軟體實體(類型, 模組, 函數等)應該是可擴展但不能被修改的

```go
package main

import (  
    "fmt"
)

type SalaryCalculator interface {  
    CalculateSalary() int
}

type Permanent struct {  
    empId    int
    basicpay int
    pf       int
}

type Contract struct {  
    empId  int
    basicpay int
}

//salary of permanent employee is sum of basic pay and pf
func (p Permanent) CalculateSalary() int {  
    return p.basicpay + p.pf
}

//salary of contract employee is the basic pay alone
func (c Contract) CalculateSalary() int {  
    return c.basicpay
}

/*
total expense is calculated by iterating though the SalaryCalculator slice and summing  
the salaries of the individual employees  
*/
func totalExpense(s []SalaryCalculator) {  
    expense := 0
    for _, v := range s {
        expense = expense + v.CalculateSalary()
    }
    fmt.Printf("Total Expense Per Month $%d", expense)
}

func main() {  
    pemp1 := Permanent{1, 5000, 20}
    pemp2 := Permanent{2, 6000, 30}
    cemp1 := Contract{3, 3000}
    employees := []SalaryCalculator{pemp1, pemp2, cemp1}
    totalExpense(employees)

}
```

`SalaryCalculator` interface 的設計使得 `totalExpense` 可以擴展新的員工類型而不需要修改任何程式碼.

若公司新增了一種員工類型 `Freelancer` 且有著不同的薪資結構, `Freelancer` 只需傳遞到 `totalExpense` slice parameter, 不需要修改 `totalExpense` method 本身.

只要 `Freelancer` 也實現 `SalaryCalculator` interface 即可實現

## Embedding Interface

Go 中沒有提供繼承機制, 但可以通過 embedding interface 來創建一個新 interface 實現

```go
package main

import (  
    "fmt"
)

type SalaryCalculator interface {  
    DisplaySalary()
}

type LeaveCalculator interface {  
    CalculateLeavesLeft() int
}

type EmployeeOperations interface {  
    SalaryCalculator
    LeaveCalculator
}

type Employee struct {  
    firstName string
    lastName string
    basicPay int
    pf int
    totalLeaves int
    leavesTaken int
}

func (e Employee) DisplaySalary() {  
    fmt.Printf("%s %s has salary $%d", e.firstName, e.lastName, (e.basicPay + e.pf))
}

func (e Employee) CalculateLeavesLeft() int {  
    return e.totalLeaves - e.leavesTaken
}

func main() {  
    e := Employee {
        firstName: "Naveen",
        lastName: "Ramanathan",
        basicPay: 5000,
        pf: 200,
        totalLeaves: 30,
        leavesTaken: 5,
    }
    var empOp EmployeeOperations = e
    empOp.DisplaySalary()
    fmt.Println("\nLeaves left =", empOp.CalculateLeavesLeft())
}
```

創建新 interface `EmployeeOperations` 並嵌套兩個 interface: `SalaryCalculator` & `LeaveCalculator`.

如果一個類型定義了 `SalaryCalculator` 及 `LeaveCalculator` interface 中的 method, 就稱該類型實現了 `EmployeeOperations` interface
