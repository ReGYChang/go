

# Kafka

傳統上來說 Kafka 是一個分散式, 基於 `pub/sub` pattern 的 message queue

Message Producer 不會將 message 直接發送給特定的 subscriber, 而是會將 message 分為不同的 topics, subscriber 只接收特定 topic 的 message

Kafka 也被定義為一個 open source 的 `event streaming platform`, 被廣泛使用於高性能的 data pipeline, stream analysis, data integration 及 service application

# Message Queue

目前比較常見的 message queue 主要有 Kafka, ActiveMQ, RabbitMQ, RocketMQ 等

在大數據應用場景下主要採用 Kafka 作為 message queue 的解決方案

傳統的 message queue 主要應用場景包括 `caching`, `decoupling` 及 `asynchronous communication`

- Caching: 有助於控制和優化 data stream 經過系統的速度, 解決 producer 和 consumer 處理速度不一致的問題
- Decoupling: 允許獨立擴展或修改 data source 及 destination 雙邊的處理過程, 只要確保其遵守同樣的 interface 約束
- asyn: 允許 message 放入 queue, 等到需要處理的時候再進行處理, 不需同步處理

