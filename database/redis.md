- [Introduction](#introduction)
  - [Why Redis](#why-redis)
  - [Redis Use Cases](#redis-use-cases)
- [Data Structure](#data-structure)
  - [Redis Data Types](#redis-data-types)
    - [String](#string)
    - [Hash](#hash)
    - [List](#list)
    - [Set](#set)
    - [ZSet](#zset)
    - [HyperLogLog](#hyperloglog)
    - [Bitmap](#bitmap)
    - [Geo](#geo)
  - [Stream Types(v5.0)](#stream-typesv50)
  - [Redis Object](#redis-object)
    - [redisObject](#redisobject)
    - [Shared Object](#shared-object)
  - [Data Structure Implementation](#data-structure-implementation)
- [Persistence](#persistence)
  - [RDB](#rdb)
  - [AOF](#aof)
- [Pub/Sub](#pubsub)
  - [Based on Channel](#based-on-channel)
  - [Based on Pattern](#based-on-pattern)
- [Event](#event)
- [Transaction](#transaction)
  - [Standard Transaction Execution](#standard-transaction-execution)
  - [CAS Lock Implementation](#cas-lock-implementation)
- [High Availability](#high-availability)
- [Scalability](#scalability)
- [Application](#application)

# Introduction

> Redis 是一種支持 key-value 等多種資料結構的存儲系統, 可用於 cache, event publish 或 subscribe, 且支持網絡, 提供 string, hash, list, set 等資料類型, 且基於記憶體, 資料可實現持久化

全名為 **Remote Dictionary Server**, 使用 C 撰寫

![redis_object_impl](img/redis_object_impl.png)

## Why Redis

Redis 有以下幾種特點:
- r/w 性能優異: read 的速度可達 110k/s, write 可達 81k/s
- 資料類型豐富: 支持 binary Strings, Lists, Hashes, Sets 及 Ordered Sets 資料類型
- atomic: 所有操作都是 atomic operation, 還支持對幾個操作全併後的 atomic 執行
- 豐富的功能: 支持 publish/subscribe, notification, key expiration 等
- persistence: 支持 RDB, AOF 等持久化方式
- pub/sub
- distributed

官方 bench-mark 根據以下條件獲得測試結果(read 110k/s, write 81k/s)
- 50 parallel clients 併發執行 100k requests
- 更新讀取的值為 256 bytes string
- Linux 2.6, X3320 Xeon 2.5 ghz
- 文本執行使用 loopback interface(127.0.0.1)

## Redis Use Cases

> Cache hotspots

Cache 是 Redis 最常見的應用場景, 主要因為 r/w 性能優異且逐漸取代 memcached, 成為首選 server side cache component

且 Redis 支持 transaction, 能有效保證資料一致性

作為 cache 使用時一般有兩種方式保存資料:
- 取資料前先讀取 redis, 若無資料再讀取 database 並將資料寫進 redis
- insert data 時同時寫入 redis

第一種方式實現簡單, 但須注意:
- 避免 `cache breakdown`
- 資料實時性較差

第二種方式資料實時性強, 但開發時不便於統一處理

> Timelimit 業務應用

redis 支持 key expiration, 可以利用這一功能在限時的優惠活動訊息, 手機驗證碼等場景

> Counter 應用

redis 由於 `incrby` 指令可以實現 atomic 遞增, 可以運用於高併發秒殺活動, 分散式序列號生成, 具體業務還體現在如限制一個手機號碼發多少條簡訊, 一個 api 一分鐘限制多少 requests, 一個 api 一天限制調用幾次等

> Distributed Lock

主要利用 redis `setnx` 指令完成, `setnx(set if not exists)` 就是若不存在則成功設置緩存並返回 1, 否則返回 0

在 server cluster 中可能兩台機器上運行定時任務, 首先可以通過 `setnx` 設置一個 lock, 若成功設置則執行; 否則表明該定時任務已執行

此場景主要用在如秒殺系統等場景

> 延遲操作

比如在訂單產生後先佔用庫存, 10 分鐘後去檢查 user 是否真正購買, 若沒有購買則將此訂單設置無效並釋放庫存

由於 redis 從 2.8.0 後提供 `keyspace notifications`, 允許客戶訂閱 pub/sub channel, 以便以某種方式接收 redis 資料集的 event

上述需求可以在訂單產生時設置一個 key, 同時設置 10 分鐘過期, 並實現一個監聽器監聽 key 時效, key 失效後再做後續操作

此需求也可以利用 message queue 的延遲隊列來實現

> Ranking

Relational Database 在 ranking 方面查詢速度普遍偏慢, 可以借住 redis 的 `SortedSet` 進行熱點資料的排序

比如按讚排行榜, 可以做一個 `SortedSet`, 以 user openid 作為 username, 以 user 按讚數作為 score, 然後針對每個 user 做一個 `hash`, 通過 `zrangebyscore` 就可以按照點讚數獲取排行榜, 再根據 username 獲取 user hash 訊息

> 按讚, 好友等相互關係儲存

利用 redis 集合的一些指令, 如交集, 聯集, 差集等

在微薄應用中, 每個 user 關注的人存在一個集合中, 很容易實現找出兩人共同好友的功能

> Simple Queue

由於 redis 有 `list push` 和 `list pop` 這樣的指令, 所以能很方便的執行 queue operation

# Data Structure

## Redis Data Types

### String

### Hash

### List

### Set

### ZSet

### HyperLogLog

### Bitmap

### Geo

## Stream Types(v5.0)

## Redis Object

### redisObject

### Shared Object

## Data Structure Implementation

# Persistence

## RDB

## AOF

# Pub/Sub

## Based on Channel

## Based on Pattern

# Event

# Transaction

## Standard Transaction Execution

## CAS Lock Implementation

# High Availability
# Scalability
# Application
