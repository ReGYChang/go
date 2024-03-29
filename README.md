<h1 align="center"> :penguin: Backend Developer in Go

<p align="center">
  <a href="#Go"><img src="https://img.shields.io/badge/language-Go-blue.svg" alt="Go"></a>
  <a href="https://regy.dev"><img src="https://img.shields.io/badge/Blog-ReGY's Inspiration-critical.svg" alt="Blog"></a>
  <a href="https://github.com/ReGYChang/LeetCode"><img src="https://img.shields.io/badge/algo-leetcode-brightgreen.svg" alt="leetcode"></a>
</p>

# Go
- [package, init, tools, repo, dependency](go/pkg_init_tools_repo_dependency.md)
- [var, const, type, pointer](go/var_const_type_pointer.md)
- [function, defer, closure](go/function_defer_closure.md)
- [array, slice, map](go/array_slice_map.md)
- [packages](go/packages.md)
- [struct, method](go/struct_methods.md)
- [interface, reflection, generic](go/interface_reflection_generic.md)
- [I/O](go/io.md)
- [error, panic, recover, test](go/error_panic_recover_test.md)
- [goroutine, channel, buffer, select, mutex](go/go_channel_buffer_select_mutex.md)
- [interview](go/interview.md)

# Data Structures And Algorithms
  - String
  - Array, Linked List
  - Stack, Queue
  - [Tree](algo/tree.md)
  - Graph
  - Analyzing Algorithms
  - Searching Algorithm
  - Sorting Algorithm
  - Divide and Conquer Algorithm
  - Greedy Mehtodology
  - Recursion
  - Backtracking Algorithm
  - Dynamic Programming

# General Development Skills
  - [Git](general/git.md)
  - [HTTP / HTTPS](general/http_https.md)
  - SQL Fundamentals
  - Data Structures and Algorithms
  - Scrum, Kanban or other project strategies
  - [Basic Authentication, OAuth, JWT, etc](general/authentication.md)
  - SOLIC, YAGNI, KISS
  - [System Design](general/system_design.md)
  - [Design Pattern](general/design_pattern.md)
  - [Domain-driven Design](general/ddd.md)
  - [Microservice](general/microservice.md)
# CLI
  - cobra
  - [urfave/cli](cmd/urfave_cli.md)
# Web Frameworks & Routers
  - Echo
  - Fiber
  - Gin
  - [gorilla/mux](routers/gorilla_mux.md)
  - [Gee](routers/gee.md)
# Network
  - [net/http](network/net_http.md)
  - [socket](network/socket.md)
  - [WebSocket](network/websocket.md)
  - [RPC](network/rpc.md)
# Gateway
  - nginx
  - [traefik](gateway/traefik.md)
  - kong
  - gRPC-Gateway
# ORMs
  - Gorm
# Database
  - Relational
      - PostgreSQL
  - Document
      - [MongoDB](database/mongodb.md)
  - Serach Engine
      - [Elasticsearch](database/elasticsearch.md)
  - Key-value
      - [Redis](database/redis.md)
  - Graph
      - [Neo4j](database/neo4j.md)
      - [NebulaGraph](database/nebula.md)
  - OLAP
    - [Doris](database/doris.md)
# Message Queue
  - [Kafka](mq/kafka.md)
  - [RabbitMQ](mq/rabbitmq.md)
# CDC
  - [Debezium](cdc/debezium.md)
# Data Flow Engine
  - [Flink](data-flow/flink.md)
# Caching
  - GCache
  - Distributed Cache
      - [go-redis](go_redis.md)
  - [GeeCache](caching/gee_cache.md)
# Logging
  - [zerolog](logging/zerolog.md)
  - Zap
# Real-Time Communication
  - Melody
  - Centrifugo
# API Clients
  - GraphQL
  - REST
# Testing
  - [Unit Testing](testing/unit_test.md)
  - [gomock](testing/gomock.md)
  - [testify](testing/testify.md)
  - [bxcodec/faker](library/bxcodec_faker.md)
  - Benchmarking
# Virtualizaion
  - [Docker](virtualization/docker.md)
    - [Neo4j](virtualization/docker/neo4j/docker-compose.yml)
    - [Redis](virtualization/docker/redis/docker-compose.yml)
    - [RabbitMQ](virtualization/docker/rabbitmq/docker-compose.yml)
    - [Kafka(arm64)](virtualization/docker/kafka/arm64/docker-compose.yml)
  - [Kubernetes](virtualization/k8s.md)
# Good to Know Libraries
  - [golang-migrate/migrate](library/migrate.md)
  - Encoding/Decoding
  - Input and output
  - Validator
  - Glow
  - GJson
  - Authboss
  - Go-Underscore
  - MicroServices
      - Message-Broker
          - Kafka
          - RabbitMQ
      - Frameworks
          - rpcx
          - Go-kit
          - Micro
          - go-zero
      - Building event-driven
          - Watermill
      - RPC
          - Protocol Buffers
          - gRPC-Go
          - gRPC-gateway
  - Task Scheduling
      - gron