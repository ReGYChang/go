

# Distributed Caching

> 商業世界中現金為王; 架構世界中緩存為王

在計算機系統中, cahce 無處不在: 

當我們訪問一個網頁, 網頁和引用的 js/css 等靜態文件會依據不同策略cache 在 browser local 或 CDN server, 當第二次再進行訪問時就會覺得加載速度快了很多

或是在 twitter 的按讚數量, 不可能每個人每次訪問都要從資料庫中查找所有點讚的紀錄再統計, 資料庫操作十分耗時, 無法支撐如此龐大的流量, 所以一般點讚這種資料是 cache 在 Redis server 中實現的

Caching 最簡單的莫過於儲存在記憶體中的 key-value pair, 在 Go 中為 map, 那直接使用 map 來做 cache 會有什麼問題?

- 記憶體不足怎麼辦?

    那就需要刪除一些資料, 重點是該怎麼刪? 應該隨機刪除還是按照時間順序刪除? 不同資料的請求頻率不同, 是否應該優先刪除請求頻率低的資料? 資料的請求頻率也可能隨著時間變化

- 併發寫入 conflict 怎麼辦?

    對於 cache request 一般不可能 Serializabe, map 操作並不是 thread-safe, 針對 concurrency 的場景做修改操作需要加鎖

- Standalone 性能不足怎麼辦?

    Standalone 資源有限, 隨著業務需求及請求增加很容易遇到瓶頸, 大部分情況下會選擇利用多台 server 資源並行處理以提升性能

...

# GeeCache

設計一個 distributed caching system 需要考慮資源控制, 淘汰策略, 併發, distributed nodes communication 等各方面問題

而且針對不同的應用場景還需要在不同特性間權衡, 如是否需要支持 cache update? 還是假設 cache 在淘汰前不允許改變?

[groupcache](https://github.com/golang/groupcache) 是用 Go 實現的 memcached, 目的是在某些特定場合替代 memcached, 其作者也是 memcached 的作者

GeeCache 旨在模仿 [groupcache](https://github.com/golang/groupcache) 實現並淬煉其中精華, 支持的特性如下:
- standalone cache 和基於 HTTP 的 distributed cache
- LRU cache strategy
- 使用 Go lock 防止 cache breakdown
- 使用 Consistent hashing 選擇節點以實現 load balance
- 使用 protobuf 優化 nodes communication
- ...

# LRU Cache Eviction Strategy

GeeCache cache 全部都儲存在 memory 中, 存儲空間十分有限, 因此不可能無限制地新增資料, 當記憶體中資料量達到一定 threshold 就需要從 cache 淘汰資料, 那應該怎麼訂淘汰策略?

下面簡單介紹最常用的三種 cache eviction strategies: FIFO, LFU 和 LRU

## FIFO(First In First Out)

FIFO 即淘汰 cache 中最老的資料, 其認為越早被新增的資料不再被使用的可能性可能比最新增加的資料高, 其實現也非常簡單: 創建一個 queue 將新增資料 push, 當記憶體不夠時再從 queue 中 pop data

但是大多數場景下部分資料雖然最早被新增進 cache 但也是最常被請求, 卻因為 FIFO 的關係被淘汰, 頻繁的新增到 cache 中又被淘汰而導致 hit ratio 降低

## LFU(Least Frequently Used)

將 Cache 中請求頻率最低的資料淘汰, 其認為資料在過去被請求越多次則在將來被請求的頻率也會越高, 實現需要維護一個按照請求次數排序的 list, 每次請求則請求數加一, 且 **list re-sorting**, 淘汰時直接選擇請求數最少的資料即可

LFU 演算法的 **hit ratio 較高**, 但缺點為需要維護每個資料的請求次數, 對記憶體的需求很高; 另外如果資料請求模式發生變化, LFU 需要較長時間適應, 受歷史資料的影響比較大, 如某個資料歷史上請求數很高, 但某個時間點之後幾乎不再被請求, 但因為之前請求數很高而遲遲無法被淘汰

## LRU(Least Recently Used)

淘汰最近最少被請求的資料, 相對於僅考慮時間因素的 FIFO 和僅考慮請求頻率的 LFU, LRU 相對是較為平衡的一種 eviction algo

LRU 認為如果資料最近被請求過, 那將來請求的機率也會更高, 其實現只需要維護一個 list, 如果某筆資料被請求則移到尾端, 只需要淘汰 HEAD 的資料即可

![LRU](img/LRU.png)

上圖很好地展示了 LRU 最核心的兩個資料結構:

- 綠色的是 `map`, 儲存 key-value 的映射關係, 如此一來根據 key 查找對應的 value 時間複雜度為 O(1), 插入一筆資料的時間複雜度也是 O(1)
- 紅色的為 `double linked list`, 將所有的值放在 `double linked list` 中, 如此一來當請求某個值時將其移到尾端的複雜度是 O(1), 在尾端插入一筆資料及刪除一筆資料的時間複雜度均為 O(1)

# LRU Implementation

接著來實現 LRU 演算法, 首先需要創建一個包含 map 及 double linked list 的 struct `Cache`, 方便實現後續的 CRUD operations

lru.go

```go
package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}
```

- 使用 Go 標準庫的 `list.List`
- `map[string]*list.Element` 儲存 list 中節點的指針
- `maxBytes` 為允許使用的最大記憶體空間, `nBytes` 為當前使用記憶體空間, `OnEvicted` 指某筆資料被淘汰時的 callbace function
- `entry` 是 linked list 節點的資料結構, 在 linked list 中儲存每個 value 對應的 key 好處在於淘汰 HEAD element 時需要使用 key 從 map 中刪除對應的映射值
- 考量通用性允許 `entry.value` 為實現 `Value` interface 的任意型別, 此 interface 只有一個 `Len()` 方法用於返回 value 所佔用的記憶體大小

實現 `Cache` 構造函數 `New()`:

```go
// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
```

## Search Entry

查找 linked list 節點可以分兩步, 第一步是從 map 中找到對應的 linked list element, 第二步是將該 element 移到 linked list 尾端

```go
// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}
```

- 若 key 對應的 element 存在則將對應 element 移動到 list 尾端並返回找到的 value
- `c.ll.MoveToFront(ele)` 即將 linked list element `ele` 移動到 list 尾端

## Delete Entry

這裡的刪除實際指的是 cache eviction, 即移除最近最少訪問的 element

```go
// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
```

- `c.ll.Back()` 取得 HEAD element, `c.ll.Remove()` 將 element 從 linked list 中刪除
- `delete(c.cache, kv.key)` 從 map `c.cache` 刪除 element 映射關係
- 更新當前所用的記憶體空間 `c.nbytes`
- 若 callbace func `OnEvicted` 不為 nil 實則調用

## Insert & Update Entry

```go
// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}
```

- 若 key 存在則更新對應 element value 並將其移到 linked list 尾端
- 不存在則插入資料, 首先 linked list 尾端新增 element `&entry{key, value}`, 並在 map 中新增 key 和 element 映射關係
- 更新 `c.nbytes`, 如果超過預設 threshold `c.maxBytes` 則移除最少請求的 element

最後為了方便測試, 實現 `Value` interface 用來獲取新增了多少筆資料

```go
// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
```

## Unit Test

測試使用 `Get()` 方法新增資料:

lru_test.go

```go
type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}
```

測試當記憶體超過了 default threshhold 是否會觸發 element eviction:

```go
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}
```

測試 callback function 是否能被調用:

```go
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
```

# Standalone Concurrent Cache

- 使用 `sync.Mutex` 並實現 LRU cache 的 concurrent control
- 實現 GeeCache 核心資料結構 `Group`, 當 cache 不存在時調用 callback function 獲取原始資料

## sync.Mutex

當多個 goroutines 同時讀寫同一個變數時, 在 high concurrency 的情況下有可能會發生衝突, 而確保同一時間只有一個 goroutine 可以訪問該變數以避免衝突, 稱為 mutex, mutex lock 可以解決此問題

> `sync.Mutex` 是一個 mutex lock, 可以由不同的 goroutine 加鎖及解鎖

Go 提供了 mutex lock `sync.Mutex`, 當一個 goroutine 獲得 lock 所有權後, 其他請求 lock 的 goroutine 就會 blocking 在 `Lock()` 方法的調用上, 直到 `Unlock()` 鎖被釋放

假設有十個併發的 goroutines 打印同一個數字 100, 為了避免重複打印, 實現了 `printOnce(num int)` 函數, 並使用集合 set 來記錄打印過的數字, 若數字已打印過則不再打印

```go
var set = make(map[int]bool, 0)

func printOnce(num int) {
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
}

func main() {
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
}
```

這段程式結果會是如何?

有時候會打印 2 次, 有時候會打印 4 次, 有時候甚至還會觸發 panic, 因為對同一個資料結構 set 的訪問衝突, 再來使用 `mutex` 的 `Lock()` 和 `Unlock()` 方法將造成衝突的部分包起來:

```go
var m sync.Mutex
var set = make(map[int]bool, 0)

func printOnce(num int) {
	m.Lock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
	m.Unlock()
}

func main() {
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
}
```

這樣一來相同的數字只會被打印一次, 當一個 goroutine 調用了 `Lock()` 方法時, 其他的 goroutine 則會被 blocking, 直到 `Unlock()` 調用後將 lock 釋放, 如此一來就能避免衝突以實現互斥

`Unlock()` 還有另一種寫法:

```go
func printOnce(num int) {
	m.Lock()
	defer m.Unlock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
}
```

## Support Concurrent R/W

接下來使用 `sync.Mutex` 封裝 LRU 方法, 使其支援 concurrent R/W, 在此之前先抽象了一個 read only 的資料結構 `ByteView` 用來表示 cache 值, 是 GeeCache 核心資料結構之一

geecache/byteview.go

```go
package geecache

// A ByteView holds an immutable view of bytes.
type ByteView struct {
	b []byte
}

// Len returns the view's length
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
```

- `ByteView` 只有一個資料成員 `b []byte`, `b` 將會儲存真實的 cache 值, 使用 `byte` 型別是為了能夠支持任意的資料結構儲存, 如 string, 圖片等
- 實現 `Len() int` 方法, 在 `lru.Cache` 中要求被 cache 的物件必須實現 `Value` interface, 即 `Len()` 方法, 用於返回其所佔的記憶體大小
- `b` 是 read only, 使用 `ByteSlice()` 方法返回一個 copy, 防止 cache 值被外部程式修改

接著就可以幫 `lru.Cache` 新增併發特性了:

geecache/cache.go

```go
package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
```

- `cache.go` 實現非常簡單, 實體化 `lru`, 封裝 `get` 和 `add` 方法並增加 mutex lock `mu`
- `add` 方法中判斷了 `c.lru` 是否為 nil, 若為 nil 再創建實體, 此為 `Lazy Initialization`, 將物件創建延遲到第一次使用該物件的時候, 用於提高性能並減少記憶體需求

## Group

`Group` 是 GeeCache 最核心的資料結構, 負責與使用者互動並且控制 cache value 與 get cache value 的流程

```go
                            true
receive key --> check if be cached -----> return cache value (1)
                |  false                         true
                |-----> if get value from remote server -----> interact with remote server --> return (2)
                            |  false
                            |-----> call callback function to get the value and insert into cache --> return cache value (3)
```

將在 geecache.go 中實現主結構體 `Group`, 目前專案結構雛型已經完成了:

```go
geecache/
    |--lru/
        |--lru.go  // lru cache eviction strategy
    |--byteview.go // cache value abstraction and package
    |--cache.go    // concurrent control
    |--geecache.go // interact with user, control cache storage and get
```

接下來先實現 (1) 和 (3), (2) 的部分後續再實現

## Callback Getter

若 cache 不存在, 應該從 datasource(file, database) 獲取資料並新增到 cache 中, 而 GeeCache 是否應該支持多種 datasource 的配置呢?

結論是不應該, 原因如下:
- Datasource 種類眾多, 無法一一實現
- 擴展性不佳

因此設計一個 callback function, 當 cache 不存在時則調用 callback, 以得到原始資料

geecache/geecache.go

```go
// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
```

- 定義 interface `Getter` 和 callback `Get(key string) ([]byte, error)`
- 定義函數型別 `GetterFunc`, 並實現 `Getter` interface `Get` 方法
- 函數型別實現某一個 interface 則稱為接口型函數, 方便調用者在調用時能夠傳入函數作為參數, 也能夠傳入實現此 interface 的 struct 作為參數

寫一個 test case 保證 callback function 能夠正常運作:

```go
func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
```

- 這個 test case 中借助 `GetterFunc` 型別轉換, 將一個匿名 callback 轉換成 `Getter` interface `f`
- 調用此 interface 方法 `f.Get(key string)`, 實際上即調匿名 callback function

>💡TIP: 定義一個函數型別 F, 且實現 interface A 的方法, 並在此方法中調用自己, 這是 Go 中將其他函數(參數返回值定義與 F 一致)轉化為 interface A 的常用技巧

接下來是最核心資料結構 `Group` 的定義:

geecache/geecache.go

```go
// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}
```

- 一個 `Group` 可以認為是一個 cache namespace, 每個 `Group` 擁有一個唯一的名稱 `name`, 如可以創建兩個 Group, cache 學生成績為 scores, cache 學生資訊則為 info
- 第二個屬性為 `getter Getter`, 即 cache 未命中時獲取原始資料的 callback
- 第三個屬性為 `mainCache cache`, 即一開始實現的 concurrent cache
- 構建函數 `NewGroup` 用來實體化 `Group`, 且將 group 儲存在全局變數 `groups` 中
- `GetGroup` 用來查找特定名稱的 `Group`, 這裡只使用 `RLock()`, 因為不涉及任何衝突變數的寫操作

再來是 GeeCache 最核心的方法 `Get`:

```go
// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
```

- `Get` 方法實現了上述流程中的 (1) 和 (3)
- (1): cache 不存在時則調用 `load` 方法, `load` 調用 `getLocally`(分散式場景會調用 `getFromPeer` 從其他節點獲取), `getLocally` 調用使用者 callback `g.getter.Get()` 來取得原始資料, 且將原始資料新增到 `mainCache` 中(通過 `populateCaceh` 方法)

至此, standalone concurrent caching 即完成

## Testing

首先用一個 map 模擬耗時的 db:

```go
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}
```

創建 `group` instance, 並測試 `Get` 方法:

```go
func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
```

這個 test case 主要測試兩種情況:
- 在 cache 為空的情況下能夠通過 callback 獲取到 source data
- 在 cache 存在情況下是否直接從 cache 中取得資料, 使用 `loadCounts` 統計某個 key 調用 callback function 的次數, 若次數大於 1 則表示調用了多次 callback function, 沒有 cache

測試結果如下:

```go
$ go test -run TestGet
2020/02/11 22:07:31 [SlowDB] search key Sam
2020/02/11 22:07:31 [GeeCache] hit
2020/02/11 22:07:31 [SlowDB] search key Tom
2020/02/11 22:07:31 [GeeCache] hit
2020/02/11 22:07:31 [SlowDB] search key Jack
2020/02/11 22:07:31 [GeeCache] hit
2020/02/11 22:07:31 [SlowDB] search key unknown
PASS
```

可以觀察到當 cache 為空時調用了 callback function, 第二次訪問時則直接從 cache 中讀取