- [Creational](#creational)
  - [Singleton Design Pattern](#singleton-design-pattern)
  - [Factory Design Pattern](#factory-design-pattern)
  - [Builder Design Pattern](#builder-design-pattern)
  - [Prototype Design Pattern [x]](#prototype-design-pattern-x)
- [Structural](#structural)
  - [Proxy Design Pattern](#proxy-design-pattern)
  - [Bridge Design Pattern](#bridge-design-pattern)
  - [Decorator Design Pattern](#decorator-design-pattern)
  - [Adapter Design Pattern](#adapter-design-pattern)
  - [Facade Design Pattern [x]](#facade-design-pattern-x)
  - [Composite Design Pattern [x]](#composite-design-pattern-x)
  - [Flyweight Design Pattern [x]](#flyweight-design-pattern-x)
- [Behavioral](#behavioral)
  - [Observer Design Pattern](#observer-design-pattern)
  - [Template Method Design Pattern](#template-method-design-pattern)
  - [Strategy Method Design Pattern](#strategy-method-design-pattern)
  - [Chain Of Responsibility Design Pattern](#chain-of-responsibility-design-pattern)
  - [State Design Pattern](#state-design-pattern)
  - [State Design Pattern](#state-design-pattern-1)
  - [Iterator Design Pattern](#iterator-design-pattern)
  - [Visitor Design Pattern [x]](#visitor-design-pattern-x)
  - [Memento Design Pattern [x]](#memento-design-pattern-x)
  - [Command Design Pattern [x]](#command-design-pattern-x)
  - [Interpreter Design Pattern [x]](#interpreter-design-pattern-x)
  - [Mediator Design Pattern [x]](#mediator-design-pattern-x)
- [Synchronization](#synchronization)
- [Concurrency](#concurrency)
- [Messaging](#messaging)
- [Stability](#stability)
  
# Creational

## Singleton Design Pattern

當某個 struct 只允許有一個 instance 存在時就會使用 singleton pattern, 以下是需要創建 singleton instance 的一些場景:

- database instance: 一般對於一個應用只需要一個 database object instance
- logger instance: 對於一個應用也只需要一個 logger object instance

Singleton instance 通常在 struct 初始化的時候創建, 通常若某個 struct 只需要創建一個 instance 時會為其定義一個 `getInstance()` 方法, 創建的 singleton instance 會通過這個方法返回給調用者

## Factory Design Pattern

## Builder Design Pattern

## Prototype Design Pattern [x]

# Structural

## Proxy Design Pattern

## Bridge Design Pattern

## Decorator Design Pattern

## Adapter Design Pattern

## Facade Design Pattern [x]

## Composite Design Pattern [x]

## Flyweight Design Pattern [x]

# Behavioral

## Observer Design Pattern

## Template Method Design Pattern

## Strategy Method Design Pattern

## Chain Of Responsibility Design Pattern

## State Design Pattern

## State Design Pattern

## Iterator Design Pattern

## Visitor Design Pattern [x]

## Memento Design Pattern [x]

## Command Design Pattern [x]

## Interpreter Design Pattern [x]

## Mediator Design Pattern [x]

# Synchronization
# Concurrency
# Messaging
# Stability