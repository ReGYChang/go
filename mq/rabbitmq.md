- [Overview](#overview)

# Overview

`RabbitMQ` 是基於 `AMQP` 構建的 message queue

`AMQP` 是一個提供統一 message service 的 application layer 標準高級訊息隊列協議, 為 application layer 的一個開放協議

`AMQP` 中幾個重要概念如下:

- Server: 負責接收 client connection 以實現 AMQP service
- Connection: application 與 MQ server 的 network connection, TCP connection
- Channel: `Message` 讀寫操作在 channel 中進行, client 可以建立多個 `Channel`, 每個 `Channel` 代表一個 session job
- Message: application 與 server 間傳送的資料, 由 `Properties` 和 `Body` 組成
- Property: 為 `Message` 包裝, 可以對 `Message` 進行修飾, 如 message priority, delay 等
- Body: `Message` 實體
- Virtual Host: 用於邏輯隔離, 一個 `Virtual Host` 中可以有若干個 `Exchange` 和 `Queue`, 同個 `Virtual Host` 中不能有相同名稱的 `Exchange` 或 `Queue`
- Exchange: 接收 `Message`, 並按照 routing rules 將 `Message` routing 到一個或多個 queue; 若 routing 失敗則返回給 producer 或直接丟棄; RabbitMQ 常用 `Exchange` 類型有 `direct`, `topic`, `fanout`, `headers` 四種
- Binding: `Exchange` 和 `Queue` 之間的 virtual connection, `Binding` 中可以包含一個或多個 `RoutingKey`
- RoutingKey: producer 將 `Message` 發送給 `Exchanger` 時會發送一個 `RoutingKey`, 用來指定 routing rules, 如此一來 `Exchange` 就知道將 `Message` 發送到哪個 `Queue`; `RoutingKey` 通常為一個 **.** 分割的字符串, 如 **com.rabbitmq**
- Queue: message queue, 用來保存 `Message` 供 consumer consume