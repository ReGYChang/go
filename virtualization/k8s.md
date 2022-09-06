- [Introduction](#introduction)
  - [Docker Container](#docker-container)

# Introduction

## Docker Container

> Docker provides a way to run applications securely isolated in a container, packaged with all its dependencies and libraries.Build once, Run anywhwere.

`Docker` 提供了一種將應用程式安全且隔離運行的一種方式, 能夠將應用程式 dependency 和 package files 打包於一個容器中, 後續即可在任何地方運行, 達到 **bulild one, run anywhere** 的目的

![docker_architecture](img/Untitled.png)

Docker 組成元件:

- Docker Daemon: 容器管理組件, 負責負載容器, 鏡像, 存儲, 網絡等管理
- Docker Client: 容器客戶端, 負責與 `Docker Daemon` 交互並完成容器生命週期管理
- Docker Registry: 容器鏡像倉庫, 負責儲存, 分發及打包
- Docker Object: 容器物件, 主要包含 `cointainer` 和 `image`

容器為應用程式開發環境帶來極高的便利性, 從根本上解決了容器的環境依賴, 打包等問題

然而容器帶來的便利性同時也夾帶著新的挑戰:

- 容器如何調度, 分發
- 多台機器如何協同工作
- Docker 主機故障時應用如何復原
- 如何保障應用 HA autoscaling