

# What is Elasticsearch?

> Elasticsearch 是一個非常強大的 real-time distributed serach engine, 基於 Lucene 實現

其主要被應用於 `full-text search`, `structural search` 及 `analytics`

除了 search, 結合 `Kibana`, `Logstash`, `Beats` 等 open source 產品, Elastic Stack(ELK) 還被廣泛運用在大數據實時分析領域, 包括 log analysis, metrics monitor, infomation security 等

ELK 可以完成海量結構化/非結構化資料搜尋, 創建可視化報表, 對監控資料設定報警閾值, 或通過 mechina learning 來自動識別異常狀況

Elasticsearch 是基於 `Restful API`, 使用 `Java` 開發的 search engine, 並作為 `Apache Lisence` release

> 根據 [DB Engine](https://db-engines.com/en/ranking) 排名顯示, Elasticsearch 為最受歡迎的企業級搜尋引擎

![db_engine_elasticsearch](img/db_engine_elasticsearch.png)

- 在當前軟體產業中, 搜尋為一個軟體系統或平台最基本的功能, 使用 Elasticsearch 可以創造出良好的搜尋體驗
- Elasticsearch 具備非常強大的大數據分析能力
- Elasticsearch 方便易用, 且可進行水平擴展
- Development community 活躍, 使用者基數龐大
- 擁有接近 real-time 的搜尋及分析能力

`Lucene` 可以算是當前最先進, 高性能且全功能的 search engine

但是 `Lucene` 僅為一個 library, 需要使用 Java 將 `Lucene` 集成到 application 中, 另外 `Lucene` 的工作原理十分的複雜

`Elasticsearch` 內部使用 `Lucene` 實現索引及搜尋, 透過 `RESTful API` 來封裝 `Lucene` 的複雜性, 使 `full-text search` 變得簡單易用

然而 Elasticsearch 不僅僅為 Lucene, 且也不僅僅為一個 full-text search engine, 其可以被以下三點精準定位:
- distributed real-time documents storage, 每個 field 都可以被索引及搜尋
- distributed real-time analytic search engine
- 能支撐上百個節點擴充, 並支持 PB 級結構化或非結構化資料儲存

# Elasticsearch Basic Concept

- Near Realtime(NRT): 接近 real-time, 資料在 summit index 後馬上就可以搜尋到
- Cluster: 一個 cluster 有一個 unique identifier, default 為 `elasticsearch`, 具有相同 cluster name 的 nodes 才會組成一個 cluster
- Node: 儲存 cluster data, 參與 cluster 索引和搜尋功能, node name default 為啟動時以一個隨機的 UUID 前七個字符, 通過 cluster name 在網絡中發現 member 並組成 cluster, single node 也可以為 cluster
- Index: 一個 index 為一個 document 集合, 每個 index 有 unique name, 一個 cluster 中可以有任意多個 index
- Document: 被索引的一筆資料, 索引的基本資料單元, 以 `JSON` 格式表示
- Shard: 在創建一個 index 時可以指定分成多少個 shard 來儲存, 每個 shard 本身也是一個功能完善且獨立的 `"index"`, 可以被放置在 cluster 的任意 node 上

| RDBMS               | Elasticserach          |
| ------------------- | ---------------------- |
| database            | index                  |
| table               | type(6.0.0 deprecated) |
| row                 | document               |
| column              | field                  |
| schema              | mapping                |
| index               | reverse index          |
| SQL                 | DSL                    |
| SELECT * FROM table | GET http://...         |
| UPDATE table SET    | PUT http://...         |
| DELETE              | DELETE http://...      |

# Elastic Stack

> Beats + Logstash + Elasticsearch + Kibana

![elastic_search](img/elastic_stack.png)

## Beats

`Beats` 是一個輕量型採集器平台, 這些採集器可以從 edge mechine 向 `Logstash` 或 `Elasticsearch` 發送資料, 期由 Go 進行開發, 運行效率較高, 不同的 beats 套件針對不同的 data source

## Logstash

`Logstash` 是動態資料收集管道, 擁有可擴充的 plugin 生態, 支持從不同來源收集資料並轉換, 最後將資料發送到不同的資料庫中, 能與 Elasticsearch 產生強大的協同作用, 在 2013 年被 Elastic 公司收購

其具有以下特性:
- 實時解析與轉化資料
- 可擴展性
- 可用性, 會通過持久話隊列來保證至少將運行中的事件送達一次
- 安全性, 可對資料進行傳輸加密
- 可監控

## Elasticsearch

`Elasticsearch` 可對資料進行搜尋, 分析和儲存, 其是基於 `JSON` 的分散式搜尋和分析引擎, 專門為了實現水平擴展性, 高可用性及管理便攜性而設計

其實現原理主要分為以下幾個步驟:
- 將資料提交到 Elasticsearch 中
- 通過分詞器將對應語句分詞
- 將分詞結果及權重一並存入, 在搜尋資料時根據權重將結果排名並返回

## Kibana

`Kibana` 實現資料可視化, 其作用為將 Elasticsearch 中的資料以圖表的形式呈現, 且具有可擴展的使用者介面, 可以配置並管理 Elasticsearch

Kibana 最早是基於 Logstash 創建的工具, 後被 Elastic 公司於 2013 年收購

# Search and Aggregation

可以通過 RESTful API 的方式對 Elasticsearch 進行操作

## Add data to Elasticserach

新增一筆資料到 Elasticsearch 中:

```shell
curl -X POST "localhost:9200/customer/_doc/1?pretty" -H 'Content-Type: application/json' -d'
{
  "name": "John Doe"
}
'
```

查詢剛才插入的 document:

```shell
curl -X GET "localhost:9200/customer/_doc/1?pretty"
```

output:

```json
{
    "_index" : "customer",
    "_type" : "_doc",
    "_id" : "1",
    "_version" : 1,
    "_seq_no" : 0,
    "_primary_term" : 1,
    "found" : true,
    "_source" : {
        "name" : "John Doe"
    }
}
```

## Add data in bulk

> ES 提供了批量操作, 使用批量處理 document 相對快很多, 節省了網絡往返的時間

```shell
curl -X PUT "localhost:9200/bank/_bulk?pretty" -H 'Content-Type: application/json' -d'
{ "create":{ } }
{ "account_number":1,"balance":39225,"firstname":"Amber","lastname":"Duke","age":32,"gender":"M","address":"880 Holmes Lane","employer":"Pyrami","email":"amberduke@pyrami.com","city":"Brogan","state":"IL" }
{ "create":{ } }
{ "account_number":6,"balance":5686,"firstname":"Hattie","lastname":"Bond","age":36,"gender":"M","address":"671 Bristol Street","employer":"Netagy","email":"hattiebond@netagy.com","city":"Dante","state":"TN" }
{ "create":{ } }
{ "account_number":13,"balance":32838,"firstname":"Nanette","lastname":"Bates","age":28,"gender":"F","address":"789 Madison Street","employer":"Quility","email":"nanettebates@quility.com","city":"Nogal","state":"VA" }
{ "create":{ } }
{ "account_number":18,"balance":4180,"firstname":"Dale","lastname":"Adams","age":33,"gender":"M","address":"467 Hutchinson Court","employer":"Boink","email":"daleadams@boink.com","city":"Orick","state":"MD" }
{ "create":{ } }
{ "account_number":20,"balance":16418,"firstname":"Elinor","lastname":"Ratliff","age":36,"gender":"M","address":"282 Kings Place","employer":"Scentric","email":"elinorratliff@scentric.com","city":"Ribera","state":"WA" }
{ "create":{ } }
{ "account_number":25,"balance":40540,"firstname":"Virginia","lastname":"Ayala","age":39,"gender":"F","address":"171 Putnam Avenue","employer":"Filodyne","email":"virginiaayala@filodyne.com","city":"Nicholson","state":"PA" }
{ "create":{ } }
{ "account_number":32,"balance":48086,"firstname":"Dillard","lastname":"Mcpherson","age":34,"gender":"F","address":"702 Quentin Street","employer":"Quailcom","email":"dillardmcpherson@quailcom.com","city":"Veguita","state":"IN" }
{ "create":{ } }
{ "account_number":37,"balance":18612,"firstname":"Mcgee","lastname":"Mooney","age":39,"gender":"M","address":"826 Fillmore Place","employer":"Reversus","email":"mcgeemooney@reversus.com","city":"Tooleville","state":"OK" }
{ "create":{ } }
{ "account_number":44,"balance":34487,"firstname":"Aurelia","lastname":"Harding","age":37,"gender":"M","address":"502 Baycliff Terrace","employer":"Orbalix","email":"aureliaharding@orbalix.com","city":"Yardville","state":"DE" }
{ "create":{ } }
{ "account_number":49,"balance":29104,"firstname":"Fulton","lastname":"Holt","age":23,"gender":"F","address":"451 Humboldt Street","employer":"Anocha","email":"fultonholt@anocha.com","city":"Sunriver","state":"RI" }
'
```

# Index Modules

> Index Modules are modules created per index and control all aspects related to an index.

## Index Management

在之前新增 document 時, 使用下面的方式會動態創建一個 customer 的 index:

```shell
curl -X POST "localhost:9200/customer/_doc/1?pretty" -H 'Content-Type: application/json' -d'
{
  "name": "John Doe"
}
'
```

這個 index 實際上已經自動創建了一個 mapping:

```json
{
  "mappings": {
    "_doc": {
      "properties": {
        "name": {
          "type": "text",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        }
      }
    }
  }
}
```

如果需要對建立 index 的過程做更多的控制, 如想要確保這個 index 有數量適中的主分片, 且在新增任何資料之前分析器和 mapping 都已經被建立好, 就需要 import 兩點:
- 禁止自動創建 index
- 手動設定並創建 index

可以通過在 `config/elasticsearch.yaml` 的每個節點下添加下面的配置:

```yaml
action.auto_create_index: false
```

## Index Format

在 request body 中添加設置或是型別 mapping, 如下所示:

```json
PUT /my_index
{
    "settings": { ... any settings ... },
    "mappings": {
        "properties": { ... any properties ... }
    }
}
```

- settings: 設置 shards, replications 等配置資訊
- mappings: field mapping, type 等
  - properties: object fields or nested fields

## Create Index

首先創建一個 `user index test-index-users`, 其中包含三個屬性: `name`, `age`, `remarks`, 儲存在一個 shard 及 一個 replication 上

```json
PUT /test-index-users
{
  "settings": {
		"number_of_shards": 1,
		"number_of_replicas": 1
	},
  "mappings": {
    "properties": {
      "name": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "age": {
        "type": "long"
      },
      "remarks": {
        "type": "text"
      }
    }
  }
}
```

output:

```json
{
    "acknowledged" : true,
    "shards_acknowledged" : true,
    "index" : "test-index-users"
}
```

新增資料測試:

```json
POST /test-index-users/_doc
{
    "name" : "pdai test name",
    "age" : 18,
    "remarks" : "hello world"
}
```

查詢資料測試:

```json
GET /test-index-users/_search
{
    "query" : {"match_all" : {}}
}
```

測試不匹配的資料型別(age):

```json
POST /test-index-users/_doc
{
  "name": "test user",
  "age": "error_age",
  "remarks": "hello eeee"
}
```

會報 `mapper parsing exception` 的錯誤

## Update Index

查看剛才創建的 index:

```shell
curl 'localhost:9200/_cat/indices?v' | grep users
```

output:

```
yellow open test-index-users                          LSaIB57XSC6uVtGQHoPYxQ 1 1     1    0   4.4kb   4.4kb
```

這邊需注意剛創建的 index status 為 `yellow`, 因為測試環境為單節點環境, 無法創建 replication, 但在上述 `number_of_replicas` 中設置了 replication 數量為 1

這時可以按需要修改 index 配置, 將 replication 數量修改為 0:

```
PUT /test-index-users/_settings
{
  "settings": {
    "number_of_replicas": 0
  }
}
```

再次查看狀態:

```
green open test-index-users                          LSaIB57XSC6uVtGQHoPYxQ 1 1     1    0   4.4kb   4.4kb
```

## Open/Close Index

一旦 index 被關閉, 則此 index 只能顯示 metadata, 無法進行任何讀寫操作:

```json
POST /test-index-users/_close
```

關閉 index 後再插入資料:

```json
POST /test-index-users/_doc
{
    "name" : "test user2",
    "age" : 18,
    "remarks" : "hello user2"
}
```

會報 `index_closed_exception` 錯誤

再打開 index:

```json
POST /test-index-users/_open
```

output:

```json
{
    "acknowledged" : true,
    "shards_acknowledged" : true
}

此時又可以重新寫入資料:

```json
POST /test-index-users/_doc
{
    "name" : "test user2",
    "age" : 18,
    "remarks" : "hello user2"
}
```

## Delete Index

將前面創建的 `test-index-users` 刪除:

```json
DELETE /test-index-users
```