- [What is Elasticsearch?](#what-is-elasticsearch)
- [Elasticsearch Basic Concept](#elasticsearch-basic-concept)
- [Elastic Stack](#elastic-stack)
  - [Beats](#beats)
  - [Logstash](#logstash)
  - [Elasticsearch](#elasticsearch)
  - [Kibana](#kibana)
- [Index Modules](#index-modules)
  - [Index Management](#index-management)
  - [Index Format](#index-format)
  - [Create Index](#create-index)
  - [Search Index](#search-index)
  - [Update Index](#update-index)
  - [Delete Index](#delete-index)
  - [Open/Close Index](#openclose-index)
- [Document Operations](#document-operations)
  - [Create Document](#create-document)
  - [Search Document](#search-document)
  - [Update Document](#update-document)
  - [Delete Document](#delete-document)
  - [Bulk Operations](#bulk-operations)

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

```json
# create index test_index
PUT /test_index?pretty
{
# index settings
  "settings": {
    "index": {
      "number_of_shards": 1, # shard 數量為 1, default 5
      "number_of_replicas": 1 # replication 數量 1, default 1
    }
  },
# index mapping
  "mappings": {
    "_doc": { # 型別, 建議設置為 _doc
      "dynamic": false, # 動態映射配置
# field properties
      "properties": {
        "id": {
          "type": "integer"  # 表示 field id 型別為 integer
        },
        "name": {
          "type": "text",
          "analyzer": "ik_max_word", # 儲存時使用的分詞器
          "search_analyzer": "ik_smart"  # 查詢時使用的分詞器
        },
        "createAt": {
          "type": "date"
        }
      }
    }
  }
}
```

>❗️NOTE: `dynamic` 為動態映射配置, 有三種狀態: true, 動態新增新的 field; false, 忽略新的 field, 不會新增 field mapping, 但會存在於 `_source` 中; strict, 若遇到新 field 會拋出 exception

output:

```json
{
  "acknowledged": true, # 是否在 cluster 中成功創建 index
  "shards_acknowledged": true,
  "index": "test_index"
}
```

## Search Index

```json
# 查看 index
GET /test_index

# 可以同時查看多個 index
GET /test_index,other_index

# 查看所有 index
GET /_cat/indices?v
```

output:

```json
{
  "test_index": {
    "aliases": {},
    "mappings": {
      "_doc": {
        "dynamic": "false",
        "properties": {
          "createAt": {
            "type": "date"
          },
          "id": {
            "type": "integer"
          },
          "name": {
            "type": "text",
            "analyzer": "ik_max_word",
            "search_analyzer": "ik_smart"
          }
        }
      }
    },
    "settings": {
      "index": {
        "creation_date": "1589271136921",
        "number_of_shards": "1",
        "number_of_replicas": "1",
        "uuid": "xueDIxeUQnGBQTms65wA6Q",
        "version": {
          "created": "6050499"
        },
        "provided_name": "test_index"
      }
    }
  }
}
```

## Update Index

> ES 提供了一系列針對 index 修改的語法, 包括 replication 數量, 新增 field, refresh_interval, index parser, aliases 等配置的修改

```json
# 修改 replication
PUT /test_index/_settings
{
    "index" : {
        "number_of_replicas" : 2
    }
}

# 修改 shard 刷新時間, default 為 1s
PUT /test_index/_settings
{
    "index" : {
        "refresh_interval" : "2s"
    }
}

# 新增 field age
PUT /teset_index/_mapping/_doc
{
  "properties": {
    "age": {
      "type": "integer"
    }
  }
}
```

修改完成後再次查看 index config:

```json
GET /test_index
{
  "test_index": {
    "aliases": {},
    "mappings": {
      "_doc": {
        "dynamic": "false",
        "properties": {
          "age": { #
            "type": "integer"
          },
          "createAt": {
            "type": "date"
          },
          "id": {
            "type": "integer"
          },
          "name": {
            "type": "text",
            "analyzer": "ik_max_word",
            "search_analyzer": "ik_smart"
          }
        }
      }
    },
    "settings": {
      "index": {
        "refresh_interval": "2s", #
        "number_of_shards": "1", #
        "provided_name": "test_index",
        "creation_date": "1589271136921",
        "number_of_replicas": "2",
        "uuid": "xueDIxeUQnGBQTms65wA6Q",
        "version": {
          "created": "6050499"
        }
      }
    }
  }
}
```

## Delete Index

```json
# 刪除 index
DELETE /test_index

# 驗證 index 是否存在
HEAD test_index
return: 404 - Not Found
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

# Document Operations

可以通過 RESTful API 的方式對 Elasticsearch document 進行操作

## Create Document

```json
# 新增單筆資料並指定 document id 為 1
PUT /test_index/_doc/1?pretty
{
  "name": "Regy"
}

# 新增單筆資料並自動生成 document id
POST /test_index/_doc?pretty
{
  "name": "Regy2"
}

# 使用 op_type 属性，强制执行某种操作
PUT test_index/_doc/1?op_type=create
{
  "name": "Regy3"
}
```
>❗️NOTE: `op_type=create` 強制執行時若 id 已存在, ES 會報`version_conflict_engine_exception`, `op_type` 主要應用於同步資料場景

此時可以查詢資料:

```json
GET /test_index/_doc/_search
{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": 2,
    "max_score": 1,
    "hits": [
      {
        "_index": "test_index",
        "_type": "_doc",
        "_id": "1",
        "_score": 1,
        "_source": {
          "name": "Regy"
        }
      },
      {
        "_index": "test_index",
        "_type": "_doc",
        "_id": "P7-FCHIBJxE1TMY0WNGN",
        "_score": 1,
        "_source": {
          "name": "Regy2"
        }
      }
    ]
  }
}
```

## Search Document

```json
# 根據 id 查詢單筆資料
GET /test_index/_doc/1

# output--->
{
  "_index": "test_index",
  "_type": "_doc",
  "_id": "1",
  "_version": 5,
  "found": true,
  "_source": {
    "name": "Regy-update",
    "age": 18
  }
}

# 獲取 index 中所有資料
GET /test_index/_doc/_search

# output--->
{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": 3,
    "max_score": 1,
    "hits": [
      {
        "_index": "test_index",
        "_type": "_doc",
        "_id": "P7-FCHIBJxE1TMY0WNGN",
        "_score": 1,
        "_source": {
          "name": "Regy2"
        }
      },
      {
        "_index": "test_index",
        "_type": "_doc",
        "_id": "_update",
        "_score": 1,
        "_source": {
          "name": "Regy3"
        }
      },
      {
        "_index": "test_index",
        "_type": "_doc",
        "_id": "1",
        "_score": 1,
        "_source": {
          "name": "Regy-update",
          "age": 18
        }
      }
    ]
  }
}

# 條件查詢
GET /test_index/_doc/_search
{
  "query": {
    "match": {
      "name": "2"
    }
  }
}

# output--->
{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": 1,
    "max_score": 0.9808292,
    "hits": [
      {
        "_index": "test_index",
        "_type": "_doc",
        "_id": "P7-FCHIBJxE1TMY0WNGN",
        "_score": 0.9808292,
        "_source": {
          "name": "Regy2"
        }
      }
    ]
  }
}
```

## Update Document

```json
# 根據 id 修改單筆資料
PUT /test_index/_doc/1?pretty
{
  "name": "Regy-update-after"
}
# 根據查詢條件 id=10, 修改 name=after name
POST test_index/_update_by_query
{
  "script": {
    "source": "ctx._source.name = params.name",
    "lang": "painless",
    "params":{
      "name":"after name"
    }
  },
  "query": {
    "term": {
      "id": "10"
    }
  }
}
```

>❗️NOTE: 修改語法和新增語法相同, 可以理解為根據 ID, 資料存在則更新; 不存在則新增

## Delete Document

```json
# 根據 id 刪除單筆資料
DELETE /test_index/_doc/1

# delete by query
POST test_index/_delete_by_query
{
  "query": {
    "match": {
     "name": "2"
    }
  }
}
```

## Bulk Operations

```json
POST _bulk
{ "index" : { "_index" : "test_test1", "_type" : "_doc", "_id" : "1" } }
{ "this_is_field1" : "this_is_index_value" }
{ "delete" : { "_index" : "test_test1", "_type" : "_doc", "_id" : "2" } }
{ "create" : { "_index" : "test_test1", "_type" : "_doc", "_id" : "3" } }
{ "this_is_field3" : "this_is_create_value" }
{ "update" : {"_id" : "1", "_type" : "_doc", "_index" : "test_test1"} }
{ "doc" : {"this_is_field2" : "this_is_update_value"} }

# 查詢所有資料
GET /test_test1/_doc/_search

# output--->
{
  "took": 33,
  "timed_out": false,
  "_shards": {
    "total": 5,
    "successful": 5,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": 2,
    "max_score": 1,
    "hits": [
      {
        "_index": "test_test1",
        "_type": "_doc",
        "_id": "1",
        "_score": 1,
        "_source": {
          "this_is_field1": "this_is_index_value",
          "this_is_field2": "this_is_update_value"
        }
      },
      {
        "_index": "test_test1",
        "_type": "_doc",
        "_id": "3",
        "_score": 1,
        "_source": {
          "this_is_field3": "this_is_create_value"
        }
      }
    ]
  }
}
```

>💡 POST _bulk 做了哪些操作?
- 若 index `test_test1` 不存在則創建, 同時若 id=1 document 存在則更新
- 刪除 id=2 document
- 新增 id=3 document; 若 document 存在則報 exception
- 更新 id=1 document

> 實際環境中 bulk operation 使用較多, 其可大幅縮減 IO 以提升效率