

# Kafka

傳統上來說 Kafka 是一個分散式, 基於 `pub/sub` pattern 的 message queue

Message Producer 不會將 message 直接發送給特定的 subscriber, 而是會將 message 分為不同的 topics, subscriber 只接收特定 topic 的 message

Kafka 也被定義為一個 open source 的 `event streaming platform`, 被廣泛使用於高性能的 data pipeline, stream analysis, data integration 及 service application

# Basic Concept

> Message

對於一個 message queue system 來說最基本的自然為 message, 在 `Kafka` 中 message 即 `Message`, 為 `Kafka` 中資料傳輸的基本單位, 多個 messages 會被分批次寫入 kafka, 同一批次的 messages 即為一組 message

> Producer & Consumer

`Producer` 負責產生 `Message`, 而 `Consumer` 負責獲取 `Message`

> Topic & Partition

Kafka 中所有的 messages 並不是存在一條 queue 中, 其通過 `Topic` 進行分類, 而一個 `Topic` 又可以區分為多個 `Partitions`

在同一個 `Partition` 內 `Message` 順序是能夠被保證的, 但在多個 partitions 之間 messages 順序則無法保證

`Partition` 主要作用為 `loading balance`, 當 `Producer` 將 messages 發送到一個 topic 上, Kafka 會將 messages 均衡的發佈到多個 partitions 上

作為 comsumer subscribe topic 時, 需要配置 subscribe 哪些 partitions, 一個 consumer 可以 subscribe 多個 partitions

> Offset

`Offset` 為一個遞增的整數值, 由 kafka 自動遞增後寫入每個 partition, 同個 partition 中一個 offset 對應一條 message, 也可以用來區分多個 message 之間的順序

> Broker & Cluster

一個獨立的 kafka server 稱為一個 `Broker`, 一個或多個 broker 可以組成一個 `Broker cluster`

kafka 雖然為分散式的 message queue system, 但在 cluster 中依然只有一台主要的 broker, 稱為 `controller`

每個 cluster 會自動選出一個 `cluster controller`, 主要負責:
- 管理 cluster
- 將 partition 分配給 broker 並監控 broker

在 cluster 中一個 partition 會從屬於一個 broker, 這個 broker 也可稱為此 partition 的 `leader`

同時一個 partition 也可以分配給多個 broker 進行 `Replication`, 若其中一個 broker 失效, 剩餘的 broker 可以接管 leader 位置

# Message Queue

目前比較常見的 message queue 主要有 Kafka, ActiveMQ, RabbitMQ, RocketMQ 等

在大數據應用場景下主要採用 Kafka 作為 message queue 的解決方案

傳統的 message queue 主要應用場景包括 `caching`, `decoupling` 及 `asynchronous communication`

- Caching: 有助於控制和優化 data stream 經過系統的速度, 解決 producer 和 consumer 處理速度不一致的問題
- Decoupling: 允許獨立擴展或修改 data source 及 destination 雙邊的處理過程, 只要確保其遵守同樣的 interface 約束
- asyn: 允許 message 放入 queue, 等到需要處理的時候再進行處理, 不需同步處理

Pub/sub pattern 有幾個特點:
- 可以有多個 topic
- consumer consume 後不刪除資料
- 每個 consumer 相互獨立, 都可以 consume 到資料

# Kafka Architecture

![kafka_architeture](img/kafka_architecture.png)

- 為了方便擴充並提高 throughput, 一個 topic 可以拆分為多個 partition
- 配合 partition 設計提出 consumer group design, group 內每個 consumer 可並行消費
- 為提升可用性, 可為每個 partition 增加 replication, 類似 NameNode HA(leader, follower)

# Kafka CMD

## kafka-topics.sh

| Arguments            | Type   | Description                            |
| -------------------- | ------ | -------------------------------------- |
| --boostrap-server    | String | 連接 Kafka Broker hostname & port      |
| --topic              | String | 目標 topic name                        |
| --create             | --     | create topic                           |
| --delete             | --     | delete topic                           |
| --alter              | --     | update topic                           |
| --list               | --     | list all topics                        |
| --describe           | --     | inspect topic detail info              |
| --partitions         | --     | setup number of partition              |
| --replication-factor | --     | setup number of replication of partion |
| --config             | --     | update system default config           |


## kafka-console-producer.sh

| Arguments           | Type    | Description                                  | Example                       |
| ------------------- | ------- | -------------------------------------------- | ----------------------------- |
| --bootstrap-server  | String  | `Require: `連接 Kafka Broker hostname & port | host1, host2, host3           |
| --topic             | String  | `Require: `receive topic name                | --                            |
| --broker-list       | String  | `Deprecated: `Source broker server           | host1, host2, host3           |
| --batch-size        | Integer | Messages number in a single batch            | `Default: `200                |
| --compression-codec | String  | Compression encode/decode                    | none, gzip, snappy, lz4, zstd |
| --max-block-ms      | Long    | --                                           | `Default: `60000              |

## kafka-console-consumer.sh

| Arguments          | Type    | Description                                  | Example             |
| ------------------ | ------- | -------------------------------------------- | ------------------- |
| --bootstrap-server | String  | `Require: `連接 Kafka Broker hostname & port | host1, host2, host3 |
| --topic            | String  | Consumed topic                               | --                  |
| --whitelist        | String  | regex                                        | --                  |
| --partition        | Integer | --                                           | --                  |
| --offset           | String  | --                                           | --                  |
| --from-beginning   | --      | --                                           | --                  |
| --max-messages     | Integer | --                                           | --                  |
| --isolation-level  | String  | --                                           | --                  |
| --group            | String  | 指定 consumer group ID                       | --                  |
| --timeout-ms       | Integer | --                                           | --                  |