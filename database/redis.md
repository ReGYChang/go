- [Introduction](#introduction)
  - [Why Redis](#why-redis)
  - [Redis Use Cases](#redis-use-cases)
- [Data Structure](#data-structure)
  - [Redis Data Types](#redis-data-types)
    - [String](#string)
    - [List](#list)
    - [Set](#set)
    - [Hash](#hash)
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

> 對於 redis 來說, 所有的 key 都是 string

在討論基本資料結構時, 討論的都是儲存 value 的 data types, 主要包括常見 5 種 data type:
- String
- List
- Set
- Zset
- Hash

## Redis Data Types

![redis_data_types](img/redis_data_types.png)

| Type   | Value                                | r/w performance                                                                                                             |
| ------ | ------------------------------------ | --------------------------------------------------------------------------------------------------------------------------- |
| String | String, int, float                   | 對 String 或 String 部分操作; 對 int 或 float 進行遞增或遞減                                                                |
| List   | Linked List, 每個節點都有一個 String | 對 Linked List 頭尾 push 或 pop, 讀取單個或多個元素; 根據值查找或刪除元素                                                   |
| Set    | 包含 String 的無序集合               | String 集合, 確認是否存在新增, 讀取, 刪除; 計算交集, 聯級, 差集等                                                           |
| Hash   | 包含 key-value pair 的無序散列表     | 新增, 讀取, 刪除單個元素                                                                                                    |
| Zset   | 包含 key-value pair 的有序散列表     | String 成員與 float 之間的有序映射; 元素排列順序由 float 大小決定; 新增, 讀取, 刪除單個元素以及根據分值範圍或成員來讀取元素 |

### String

> String 是 redis 中最基本的資料結構, 一個 key 對應一個 value

String 類型是 binary-safe, 即 redis 的 String 可以包含任何資料, 如數字, 字符串, jpg 圖片或者序列化物件

> Command

| Cmd    | Desc                          | Use               |
| ------ | ----------------------------- | ----------------- |
| GET    | 讀取特定 key 中的 value       | GET name          |
| SET    | 設置儲存在特定 key 中的 value | SET name value    |
| DEL    | 刪除儲存在特定 key 中的 value | DEL name          |
| INCR   | 將 key 中 value 加 1          | INCR key          |
| DECR   | 將 key 中 value 減 1          | DECR key          |
| INCRBY | 將 key 中 value 加上整數      | INCRBY key amount |
| DECRBY | 將 key 中 value 加上整數      | DECRBY key amount |

> Execution

```go
127.0.0.1:6379> set hello world
OK
127.0.0.1:6379> get hello
"world"
127.0.0.1:6379> del hello
(integer) 1
127.0.0.1:6379> get hello
(nil)
127.0.0.1:6379> set counter 2
OK
127.0.0.1:6379> get counter
"2"
127.0.0.1:6379> incr counter
(integer) 3
127.0.0.1:6379> get counter
"3"
127.0.0.1:6379> incrby counter 100
(integer) 103
127.0.0.1:6379> get counter
"103"
127.0.0.1:6379> decr counter
(integer) 102
127.0.0.1:6379> get counter
"102"
```

> Use Cases

- Cache: 將常用資料, string, 圖檔或影片等資料放到 redis 中作為 cache 以降低 db 讀寫壓力
- Counter: redis 是 single-thread model, 同時資料可以一步落地到其他 data source
- Session: session 共享

### List

> Redis 中 List 是以 double linked-list 實現

使用 List 資料結構可以輕鬆實現最新消息排隊的功能(TimeLine), 其另一個應用是 message queue, 可以利用 List push 將任務放在 List 中, 然後 worker thread 再用 pop 將任務取出執行

> Command

| Cmd    | Desc                                                                           | Use              |
| ------ | ------------------------------------------------------------------------------ | ---------------- |
| RPUSH  | 將給定 value push 至 List 右側                                                 | RPUSH key value  |
| LPUSH  | 將給定 value push 至 List 左側                                                 | LPUSH key value  |
| RPOP   | 從 List 右側 pop 出一個 value 並返回                                           | RPOP key         |
| LPOP   | 從 List 左側 pop 出一個 value 並返回                                           | LPOP key         |
| LRANGE | 讀取 List 在給定範圍內所有 value                                               | LRANGE key 0 -1  |
| LINDEX | 通過 index 讀取 List 元素; 也可以用負數下標, -1 表示 List 最後一個元素以此類推 | LINDEX key index |

> 使用 List 技巧

- LPUSH + LPOP = Stack
- LPUSH + RPOP = Queue
- LPUSH + ITRIM = Capped Collection
- LPUSH + BRPOP = Message Queue

> Execution

```go
127.0.0.1:6379> lpush mylist 1 2 ll ls mem
(integer) 5
127.0.0.1:6379> lrange mylist 0 -1
1) "mem"
2) "ls"
3) "ll"
4) "2"
5) "1"
127.0.0.1:6379> lindex mylist -1
"1"
127.0.0.1:6379> lindex mylist 10        # index不在 mylist 的区间范围内
(nil)
```

> Use Cases

- TimeLine: 發布新䩞文用 `LPUSH` 加入 timeline, 展示新的 List 訊息
- Message Queue

### Set

> Redis Set 是 String 類型的無序集合, 集合成員是唯一, 集合中不能出現重複資料

Redis Set 透過 hash table 實現, 新增, 刪除, 查詢的時間複雜度都是 O(1)

> Command

| Cmd       | Desc                                  | Use                  |
| --------- | ------------------------------------- | -------------------- |
| SADD      | 向 Set 新增一個或多個 item            | SADD key value       |
| SCARD     | 讀取 Set 成員數                       | SCARD key            |
| SMEMBERS  | 返回 Set 中所有成員                   | SMEMBERS key member  |
| SISMEMBER | 判斷 member 元素是否是 Set key 的成員 | SISMEMBER key member |

其他 set operation 參考: [https://www.runoob.com/redis/redis-sets.html](https://www.runoob.com/redis/redis-sets.html)

> Execution

```go
127.0.0.1:6379> sadd myset hao hao1 xiaohao hao
(integer) 3
127.0.0.1:6379> smembers myset
1) "xiaohao"
2) "hao1"
3) "hao"
127.0.0.1:6379> sismember myset hao
(integer) 1
```

> Use Cases
- Tag: 為 user 或 訊息新增 tag, 可以推薦同一 tag 或類似 tag 給關注的 user
- 按讚, 收藏可以以 set 實現

### Hash

> Redis hash 是一個 String 類型的 field 和 value 映射表, 適合用於儲存 object

> Command

| Cmd     | Desc                          | Use                           |
| ------- | ----------------------------- | ----------------------------- |
| HSET    | 新增 key-value pair           | HSET hash-key sub-key1 value1 |
| HGET    | 讀取指定 hash-key value       | HGET hash-key key1            |
| HGETALL | 讀取 hash 所有 key-value pair | HGETALL hash-key              |
| HDEL    | 若指定 key 存於 hash 中即移除 | HDEL hash-key sub-key1        |

> Execution

```go
127.0.0.1:6379> hset user name1 hao
(integer) 1
127.0.0.1:6379> hset user email1 hao@163.com
(integer) 1
127.0.0.1:6379> hgetall user
1) "name1"
2) "hao"
3) "email1"
4) "hao@163.com"
127.0.0.1:6379> hget user user
(nil)
127.0.0.1:6379> hget user name1
"hao"
127.0.0.1:6379> hset user name2 xiaohao
(integer) 1
127.0.0.1:6379> hset user email2 xiaohao@163.com
(integer) 1
127.0.0.1:6379> hgetall user
1) "name1"
2) "hao"
3) "email1"
4) "hao@163.com"
5) "name2"
6) "xiaohao"
7) "email2"
8) "xiaohao@163.com"
```

> Use Cases
- Cache: 相比 String 更節省空間, 更直觀

### ZSet

> Redis 有序集合, 與集合一樣也是 String 類型元素的集合, 且不允許重複成員; 不同的是每個元素都會關聯一個 double 類型的分數, redis 是透過分數來為集合成員進行排序

ZSet 成員是唯一的, 但分數(score)卻可以重複

ZSet 通過兩種資料結構實現:
- ziplist: 為了提高儲存效率而設計的一種特殊編碼的 double linked-list, 可以儲存字符串或整數, 儲存整數時是採用整數的 binary 而不是字符串形式; 能在 O(1) 的時間複雜度下完成 list 兩端的 pop 和 push, 但因為每次操作都需要重新分配 ziplist 記憶體空間, 空間複雜度較高
- zSkiplist: 其性能可以保證在查詢, 刪除, 新增等操作的時候時間複雜度為 O(log(n)) 

> Command

| Cmd    | Desc                                       | Use                            |
| ------ | ------------------------------------------ | ------------------------------ |
| ZADD   | 將一個帶有指定 score 的成員新增到 ZSet     | ZADD zset-key 178 member1      |
| ZRANGE | 根據元素在 ZSet 所處的位置從中查詢多個元素 | ZRANGE zset-key 0-1 withccores |
| ZREM   | 若指定元素成員存在於 ZSet 即移除此元素     | ZREM zset-key member1          |

更多指令參考: [https://www.runoob.com/redis/redis-sorted-sets.html](https://www.runoob.com/redis/redis-sorted-sets.html)

> Execution

```go
127.0.0.1:6379> zadd myscoreset 100 hao 90 xiaohao
(integer) 2
127.0.0.1:6379> ZRANGE myscoreset 0 -1
1) "xiaohao"
2) "hao"
127.0.0.1:6379> ZSCORE myscoreset hao
"100"
```

> Use Cases
- Ranking: Zset 適合實現各種排行榜場景

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
