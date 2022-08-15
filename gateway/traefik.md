# What is API Gateway?

API gateway 是一個 server, 封裝了系統的內部架構並為每個 client 提供一個訂製的 API

其可能具有其他職責, 如 `authentication`, `authorization`, `monitor`, `load balance`, `cache`, `protocol converter`, `rate limit`, `circuit break`, `static response handle` 等等

API gateway 核心價值為讓所有的 client or consumer 都透過統一的 gateway 接入 microservices, 在 gateway layer 處理所有的非業務邏輯功能, 通常提供 `REST/HTTP` API

# Traefik

`Træfɪk` 是一個為了使微服務部署更加簡易便捷而誕生的現代化 HTTP reverse proxy, load balance 的工具

其支援多種平台如 `Kubernetes` 或 `Swarm`, 也支援多種 service registry 如 `etcd` 或 `consul`, 且能自動化套用動態配置

其具有以下特點:
- Continuously updates its configuration (No restarts!)
- Supports multiple load balancing algorithms
- Provides HTTPS to your microservices by leveraging Let's Encrypt (wildcard certificates support)
- Circuit breakers, retry
- See the magic through its clean web UI
- Websocket, HTTP/2, GRPC ready
- Provides metrics (Rest, Prometheus, Datadog, Statsd, InfluxDB)
- Keeps access logs (JSON, CLF)
- Fast
- Exposes a Rest API
- Packaged as a single binary file (made with ❤️ with go) and available as an official docker image