- [Docker Introduction](#docker-introduction)
  - [Why Docker?](#why-docker)
  - [What's docker](#whats-docker)
  - [Docker History](#docker-history)
  - [Talk about Docker](#talk-about-docker)
  - [Docker Features](#docker-features)
- [Docker Installation](#docker-installation)
  - [Runtime Environment](#runtime-environment)
  - [Environment Check](#environment-check)
  - [Installation Steps](#installation-steps)
  - [Theory](#theory)
    - [How Docker Work？](#how-docker-work)
    - [Why Docker Faster Than VM？](#why-docker-faster-than-vm)
- [Docker CMD](#docker-cmd)
  - [Image CMD](#image-cmd)
  - [Container CMD](#container-cmd)
  - [Other CMD](#other-cmd)
  - [Summary](#summary)
- [Deploy Nginx](#deploy-nginx)
- [Deploy Tomcat](#deploy-tomcat)
- [Depoly ES and Kibana](#depoly-es-and-kibana)
- [Docker Volume](#docker-volume)
  - [Why Need Volume](#why-need-volume)
  - [How to Volume?](#how-to-volume)
- [DockerFile](#dockerfile)
  - [Build Step](#build-step)
  - [DockerFile Command](#dockerfile-command)
    - [FROM](#from)
    - [MAINTAINER](#maintainer)
    - [RUN](#run)
    - [CMD](#cmd)
    - [EXPOSE](#expose)
    - [ENV](#env)
    - [ADD](#add)
    - [COPY](#copy)
    - [ENTRYPOINT](#entrypoint)
    - [VOLUME](#volume)
    - [USER](#user)
    - [WORKDIR](#workdir)
    - [ONBUILD](#onbuild)
  - [Testing](#testing)
- [Docker Network](#docker-network)
  - [Docker0](#docker0)
  - [Custom Network](#custom-network)
    - [Network Mode](#network-mode)
- [Docker Compose](#docker-compose)

# Docker Introduction

## Why Docker?

- 環境配置麻煩，每台機器要部屬環境
- 環境更換、版本更新導致服務不可用
- 發布 jar ，專案連同環境打包
- 伺服器配置應用環境，無法跨平台
- 開發打包佈署上限線一套流程
- Java — jar (環境) — 打包專案與環境 ( 鏡像 ) — Docker 鏡像倉庫 — 下載 — 運行

## What's docker

- Docker 思想來源於集裝箱
- JRE — 多個服務( 端口衝突 )
- 隔離：核心思想，打包裝箱，箱子間互相隔離
- Docker 通過隔離機制，將伺服器利用到極致

![img/Untitled.png](img/Untitled.png)

- 鏡像 ( Image )：Docker 鏡像好比一個模板，可通過模板來創建多個容器服務
- 容器 ( Container )： 獨立運行一個或一組應用，通過鏡像來創建。
- 倉庫 ( Repository )：存放鏡像的地方，分為公有倉庫與私有倉庫。 — Docker Hub

## Docker History

- 2010 — dotCloud
- 做 PaaS 服務、Linux 相關容器技術
- 將容器化技術命名為 Docker
- Docker 剛誕生時，沒有引起行業注意
- 開源 — 每個月更新一個版本
- 2014.04.09 Docker 1.0 發布
- 為甚麼 Docker 流行?
    - VM 技術：笨重，linux Cent OS 原生鏡像( 一個主機 ) ，隔離需要多個虛擬機
    - 容器技術：隔離，鏡像( 最核心環境 4M + JDK + MySQL ) 十分輕巧，秒級啟動
    

## Talk about Docker

- 基於 Go 語言開發、開源項目
- 文檔地址：[https://docs.docker.com/](https://docs.docker.com/)
- 倉庫地址：[https://hub.docker.com/](https://hub.docker.com/)

## Docker Features

- 虛擬機缺點
    - 資源佔用多
    - 冗餘步驟多
    - 啟動慢
- 容器化技術
    - 應用直接運行在宿主主機內核中，容器自己本身沒有內核
    - 每個容器間相互隔離
- DevOps
    
    應用更快速交付與佈署
    
    - 傳統：一堆文檔，安裝程式
    - Docker：打包鏡像發布測試
    
    更便捷的升級與擴縮容
    
    更簡單系統運維
    
    更高效計算資源利用
    
    - Docker 是內核級別虛擬化，可以在一個物理機上運行很多容器實體
    

# Docker Installation

## Runtime Environment

- Linux 基礎
- CentOS 7

## Environment Check

```bash
#系統內核是 3.10 以上的
[root@regy ~]# uname -r
3.10.0-1127.el7.x86_64
```

```bash
#系統版本
[root@regy ~]# cat /etc/os-release
NAME="CentOS Linux"
VERSION="7 (Core)"
ID="centos"
ID_LIKE="rhel fedora"
VERSION_ID="7"
PRETTY_NAME="CentOS Linux 7 (Core)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:centos:centos:7"
HOME_URL="[https://www.centos.org/](https://www.centos.org/)"
BUG_REPORT_URL="[https://bugs.centos.org/](https://bugs.centos.org/)"
CENTOS_MANTISBT_PROJECT="CentOS-7"
CENTOS_MANTISBT_PROJECT_VERSION="7"
REDHAT_SUPPORT_PRODUCT="centos"
REDHAT_SUPPORT_PRODUCT_VERSION="7"
```

## Installation Steps

更新安裝包：

```bash
sudo yum update -y
```

卸載舊版本 Docker：

```bash
$ sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
```

需要的安裝包：

```bash
$ sudo yum install -y yum-utils
```

設置鏡像倉庫：

```bash
$ sudo yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
```

安裝 Docker Engine：

```bash
$ sudo yum install docker-ce docker-ce-cli [containerd.io](http://containerd.io/)
```

啟動 Docker：

```bash
$ sudo systemctl start docker
```

查看是否安裝成功：

```bash
docker version
```

hello-world：

```bash
docker run hello-world
```

![img/Untitled%201.png](img/Untitled%201.png)

- Docker 在 Local 尋找鏡像
    - True → 運行
    - False → 去 Docker Hub 下載
        - Docker 是否可以找到
            - True → 下載鏡像到 Local
            - False → 返回錯誤

查看一下 hello-world Image：

```bash
# docker images
REPOSITORY    TAG      IMAGE ID        CREATED         SIZE
hello-world   latest   bf756fb1ae65    6 months ago    13.3kB
```

卸載 Docker：

```bash
#卸載依賴
$ sudo yum remove docker-ce docker-ce-cli [containerd.io](http://containerd.io/)
```

```bash
#刪除資源
#Docker 默認工作路徑
$ sudo rm -rf /var/lib/docker
```

## Theory

### How Docker Work？

![img/Untitled%202.png](img/Untitled%202.png)

- Docker 是一個 Client - Server 架構系統，Docker 的守護進程運行在主機上，通過 Socket 從客戶端訪問
- Docker-Server 接收到 Docker-Client 請求並執行

### Why Docker Faster Than VM？

![img/Untitled%203.png](img/Untitled%203.png)

- 比 VM 更少的抽象層
- Docker 運用宿主機內核，VM 則需要 Guest OS

> 所以新建一個容器時不需要像 VM 一樣重新加載一個 OS kernal，省略複雜流程
> 

# Docker CMD

- 顯示 Docker 版本信息

```bash
docker version
```

- 顯示 Docker 系統信息

```bash
docker info
```

- 幫助指令

```bash
docker -- help
```

[Reference documentation](https://docs.docker.com/reference/)

## Image CMD

- 查看 Local 主機上的鏡像

```bash
[root@regy ~]# docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
hello-world         latest              bf756fb1ae65        6 months ago        13.3kB
```

- 解釋
    - REPOSITORY：鏡像倉庫源
    - TAG：鏡像標籤
    - IMAGE ID：鏡像ID
    - CREATED：鏡像創建時間
    - SIZE：鏡像大小
- 選項
    - -a, --all        列出所有鏡像
    - -q, --quiet   只顯示鏡像ID

---

- 搜索指令

```bash
[root@regy ~]# docker search mysql
NAME                              DESCRIPTION                                                                 STARS               OFFICIAL            AUTOMATED
mysql                             MySQL is a widely used, open-source relation…                               9757                [OK]
mariadb                           MariaDB is a community-developed fork of MyS…                               3564                [OK]
mysql/mysql-server                Optimized MySQL Server Docker images. Create…                               717                                     [OK]
percona                           Percona Server is a fork of the MySQL relati…                               498                 [OK]
centos/mysql-57-centos7           MySQL 5.7 SQL database server                                               77
mysql/mysql-cluster               Experimental MySQL Cluster Docker images. Cr…                               73
centurylink/mysql                 Image containing mysql. Optimized to be link…                               61                                      [OK]
bitnami/mysql                     Bitnami MySQL Docker Image                                                  44                                      [OK]
deitch/mysql-backup               REPLACED! Please use http://hub.docker.com/r…                               41                                      [OK]
tutum/mysql                       Base docker image to run a MySQL database se…                               35
schickling/mysql-backup-s3        Backup MySQL to S3 (supports periodic backup…                               30                                      [OK]
prom/mysqld-exporter                                                                                          29                                      [OK]
databack/mysql-backup             Back up mysql databases to... anywhere!                                     27
linuxserver/mysql                 A Mysql container, brought to you by LinuxSe…                               25
centos/mysql-56-centos7           MySQL 5.6 SQL database server                                               19
circleci/mysql                    MySQL is a widely used, open-source relation…                               19
mysql/mysql-router                MySQL Router provides transparent routing be…                               16
arey/mysql-client                 Run a MySQL client from a docker container                                  14                                      [OK]
fradelg/mysql-cron-backup         MySQL/MariaDB database backup using cron tas…                               8                                       [OK]
openshift/mysql-55-centos7        DEPRECATED: A Centos7 based MySQL v5.5 image…                               6
devilbox/mysql                    Retagged MySQL, MariaDB and PerconaDB offici…                               3
ansibleplaybookbundle/mysql-apb   An APB which deploys RHSCL MySQL                                            2                                       [OK]
jelastic/mysql                    An image of the MySQL database server mainta…                               1
widdpim/mysql-client              Dockerized MySQL Client (5.7) including Curl…                               1                                       [OK]
monasca/mysql-init                A minimal decoupled init container for mysql                                0
```

- 選項
    - - -filter=stars=3000    #stars 大於 3000

---

- 下載鏡像

```bash
docker pull mysql[:tag]

#沒寫 tag 默認是 latest
#指定版本下載 
docker pull mysl:5.7
```

---

- 刪除鏡像

```bash
docker rmi -f image-id
docker rmi -f image-id image-id image-id
docker rmi -f $(docker images -q)
```

## Container CMD

> 有鏡像才可以創建容器
> 
- 新建容器並啟動

```bash
docker run [可選參數] image

#參數說明
--name="Name"  容器名稱  tomcat01 tomcat02  用來區分容器
-d              後台方式運行
-it             使用交互方式運行，進入容器查看內容
-p              指定容器端口  -p 8080:8080
		-p ip:主機端口:容器端口
		-p 主機端口:容器端口 (常用)
		-p 容器端口
		容器端口

#測試、啟動並進入容器
[root@regy ~]# docker run -it centos /bin/bash
Unable to find image 'centos:latest' locally
latest: Pulling from library/centos
6910e5a164f7: Pull complete
Digest: sha256:4062bbdd1bb0801b0aa38e0f83dece70fb7a5e9bce223423a68de2d8b784b43b
Status: Downloaded newer image for centos:latest

[root@11873871bdae /]# ls
bin  etc   lib    lost+found  mnt  proc  run   srv  tmp  var
dev  home  lib64  media       opt  root  sbin  sys  usr

#從容器中退出
[root@11873871bdae /]# exit  
```

- 列出運行容器

```bash
#當前運行容器
[root@regy ~]# docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES

#容器運行歷史紀錄
[root@regy ~]# docker ps -a
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS                            PORTS               NAMES
11873871bdae        centos              "/bin/bash"         3 minutes ago       Exited (127) About a minute ago                       cranky_agnesi
7cc7b326772e        bf756fb1ae65        "/hello"            18 hours ago        Exited (0) 18 hours ago

#顯示最近創建的2個容器
[root@regy ~]# docker ps -a -n=2

#只顯示容器編號
[root@regy ~]# docker ps -a -q
```

- 退出容器

```bash
exit #容器停止退出
Ctrl + P + Q #容器不停止退出
```

- 刪除容器

```bash
docker rm container-id           #不能刪除正在運行容器，強制刪除使用 rm -f
docker rm -f $(docker ps -aq)    #刪除所有容器
docker ps -a -q|xargs docker rm  #刪除所有容器
```

- 啟動和停止容器操作

```bash
docker start container-id     #啟動容器
docker restart container-id   #重啟容器
docker stop container-id      #停止正在運行容器
docker kill container-id      #強制停止正在運行容器
```

## Other CMD

- 後台啟動容器

```bash
#命令 docker run -d image-name
[root@regy ~]# docker run -d centos:7
da639ea77003605b194b46ae27fa068f278ad142a5467f219eddf0aa7e335842

#問題: docker ps 發現 centos 停止了
#docker container 使用後台運行，必須要有個前台進程，docker 發現沒有應用，就會自動停止
#nginx 容器啟動後，發現自己沒有提供服務，就會立刻停止
```

- 查看日誌

```bash
docker logs --help
docker logs -tf container-id  #顯示日誌
--tail number                 #顯示日誌條數

```

- 查看容器中進程信息

```bash
#docker top container-id
```

- 查看容器 meta-data

```bash
#docker inspect container-id
[root@regy ~]# docker inspect da639ea77003
[
    {
        "Id": "da639ea77003605b194b46ae27fa068f278ad142a5467f219eddf0aa7e335842"                 ,
        "Created": "2020-07-23T03:56:58.44951221Z",
        "Path": "/bin/bash",
        "Args": [],
        "State": {
            "Status": "exited",
            "Running": false,
            "Paused": false,
            "Restarting": false,
            "OOMKilled": false,
            "Dead": false,
            "Pid": 0,
            "ExitCode": 0,
            "Error": "",
            "StartedAt": "2020-07-23T03:56:58.896960912Z",
            "FinishedAt": "2020-07-23T03:56:58.903903647Z"
        },
        "Image": "sha256:b5b4d78bc90ccd15806443fb881e35b5ddba924e2f475c1071a38a3                 094c3081d",
        "ResolvConfPath": "/var/lib/docker/containers/da639ea77003605b194b46ae27                 fa068f278ad142a5467f219eddf0aa7e335842/resolv.conf",
        "HostnamePath": "/var/lib/docker/containers/da639ea77003605b194b46ae27fa                 068f278ad142a5467f219eddf0aa7e335842/hostname",
        "HostsPath": "/var/lib/docker/containers/da639ea77003605b194b46ae27fa068                 f278ad142a5467f219eddf0aa7e335842/hosts",
        "LogPath": "/var/lib/docker/containers/da639ea77003605b194b46ae27fa068f2                 78ad142a5467f219eddf0aa7e335842/da639ea77003605b194b46ae27fa068f278ad142a5467f21                 9eddf0aa7e335842-json.log",
        "Name": "/elegant_colden",
        "RestartCount": 0,
        "Driver": "overlay2",
        "Platform": "linux",
        "MountLabel": "",
        "ProcessLabel": "",
        "AppArmorProfile": "",
        "ExecIDs": null,
        "HostConfig": {
            "Binds": null,
            "ContainerIDFile": "",
            "LogConfig": {
                "Type": "json-file",
                "Config": {}
            },
            "NetworkMode": "default",
            "PortBindings": {},
            "RestartPolicy": {
                "Name": "no",
                "MaximumRetryCount": 0
            },
            "AutoRemove": false,
            "VolumeDriver": "",
            "VolumesFrom": null,
            "CapAdd": null,
            "CapDrop": null,
            "Capabilities": null,
            "Dns": [],
            "DnsOptions": [],
            "DnsSearch": [],
            "ExtraHosts": null,
            "GroupAdd": null,
            "IpcMode": "private",
            "Cgroup": "",
            "Links": null,
            "OomScoreAdj": 0,
            "PidMode": "",
            "Privileged": false,
            "PublishAllPorts": false,
            "ReadonlyRootfs": false,
            "SecurityOpt": null,
            "UTSMode": "",
            "UsernsMode": "",
            "ShmSize": 67108864,
            "Runtime": "runc",
            "ConsoleSize": [
                0,
                0
            ],
            "Isolation": "",
            "CpuShares": 0,
            "Memory": 0,
            "NanoCpus": 0,
            "CgroupParent": "",
            "BlkioWeight": 0,
            "BlkioWeightDevice": [],
            "BlkioDeviceReadBps": null,
            "BlkioDeviceWriteBps": null,
            "BlkioDeviceReadIOps": null,
            "BlkioDeviceWriteIOps": null,
            "CpuPeriod": 0,
            "CpuQuota": 0,
            "CpuRealtimePeriod": 0,
            "CpuRealtimeRuntime": 0,
            "CpusetCpus": "",
            "CpusetMems": "",
            "Devices": [],
            "DeviceCgroupRules": null,
            "DeviceRequests": null,
            "KernelMemory": 0,
            "KernelMemoryTCP": 0,
            "MemoryReservation": 0,
            "MemorySwap": 0,
            "MemorySwappiness": null,
            "OomKillDisable": false,
            "PidsLimit": null,
            "Ulimits": null,
            "CpuCount": 0,
            "CpuPercent": 0,
            "IOMaximumIOps": 0,
            "IOMaximumBandwidth": 0,
            "MaskedPaths": [
                "/proc/asound",
                "/proc/acpi",
                "/proc/kcore",
                "/proc/keys",
                "/proc/latency_stats",
                "/proc/timer_list",
                "/proc/timer_stats",
                "/proc/sched_debug",
                "/proc/scsi",
                "/sys/firmware"
            ],
            "ReadonlyPaths": [
                "/proc/bus",
                "/proc/fs",
                "/proc/irq",
                "/proc/sys",
                "/proc/sysrq-trigger"
            ]
        },
        "GraphDriver": {
            "Data": {
                "LowerDir": "/var/lib/docker/overlay2/667e52930dee2a778fa69b3b50                 3c6e366fccfe782ee9362ebf970a4455ca5793-init/diff:/var/lib/docker/overlay2/ecb40d                 d9a6aa4346e323f475b444d8a058ad8d9562714c59fba367524f2076f7/diff",
                "MergedDir": "/var/lib/docker/overlay2/667e52930dee2a778fa69b3b5                 03c6e366fccfe782ee9362ebf970a4455ca5793/merged",
                "UpperDir": "/var/lib/docker/overlay2/667e52930dee2a778fa69b3b50                 3c6e366fccfe782ee9362ebf970a4455ca5793/diff",
                "WorkDir": "/var/lib/docker/overlay2/667e52930dee2a778fa69b3b503                 c6e366fccfe782ee9362ebf970a4455ca5793/work"
            },
            "Name": "overlay2"
        },
        "Mounts": [],
        "Config": {
            "Hostname": "da639ea77003",
            "Domainname": "",
            "User": "",
            "AttachStdin": false,
            "AttachStdout": false,
            "AttachStderr": false,
            "Tty": false,
            "OpenStdin": false,
            "StdinOnce": false,
            "Env": [
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/b                 in"
            ],
            "Cmd": [
                "/bin/bash"
            ],
            "Image": "centos:7",
            "Volumes": null,
            "WorkingDir": "",
            "Entrypoint": null,
            "OnBuild": null,
            "Labels": {
                "org.label-schema.build-date": "20200504",
                "org.label-schema.license": "GPLv2",
                "org.label-schema.name": "CentOS Base Image",
                "org.label-schema.schema-version": "1.0",
                "org.label-schema.vendor": "CentOS",
                "org.opencontainers.image.created": "2020-05-04 00:00:00+01:00",
                "org.opencontainers.image.licenses": "GPL-2.0-only",
                "org.opencontainers.image.title": "CentOS Base Image",
                "org.opencontainers.image.vendor": "CentOS"
            }
        },
        "NetworkSettings": {
            "Bridge": "",
            "SandboxID": "d95301e20052483ae695ebffb3c731d4bb41c3e020db13544271d8                 ff18aa40e3",
            "HairpinMode": false,
            "LinkLocalIPv6Address": "",
            "LinkLocalIPv6PrefixLen": 0,
            "Ports": {},
            "SandboxKey": "/var/run/docker/netns/d95301e20052",
            "SecondaryIPAddresses": null,
            "SecondaryIPv6Addresses": null,
            "EndpointID": "",
            "Gateway": "",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "IPAddress": "",
            "IPPrefixLen": 0,
            "IPv6Gateway": "",
            "MacAddress": "",
            "Networks": {
                "bridge": {
                    "IPAMConfig": null,
                    "Links": null,
                    "Aliases": null,
                    "NetworkID": "0dda68df241beea16fb5a5be63173d7c180208ec61cb06                 03c8c5f6c0b48a7f22",
                    "EndpointID": "",
                    "Gateway": "",
                    "IPAddress": "",
                    "IPPrefixLen": 0,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "MacAddress": "",
                    "DriverOpts": null
                }
            }
        }
    }
]
```

- 進入當前正在運行的容器

```bash
#需要進入後台運行的容器修改配置

#方法一
#進入容器後開啟新的 Terminal
docker exec -it container-id /bin/bash

#方法二
#進入容器正在執行的 Terminal
docker attach contain-id 
```

- 從容器內複製文件到主機上

```bash
#創建 test.java 文件
touch test.java

#docker cp container-id:容器內位置 主機位置
docker cp container-id:/home/test.java /home

#copy 是一個手動過程，可通過實現 -v 數據卷的技術實現自動同步 /home /home
```

## Summary

![img/Untitled%204.png](img/Untitled%204.png)

[Docker Command](https://www.notion.so/f557e81d4ec444f5af6abd330c5de447)

# Deploy Nginx

- docker search nginx
- docker pull nginx
- run container and test

```bash
[root@regy ~]# docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
nginx               latest              8cf1bfb43ff5        30 hours ago        132MB
centos              7                   b5b4d78bc90c        2 months ago        203MB

#-d 後台運行
#--name 容器名稱
#-p 主機端口:容器端口
[root@regy ~]# docker run -d --name nginx-01 -p 5566:80 nginx
2e25e782345bb18fb5a977fce64506b8b12fa23afedc1ad69cb2412a7addb6cf

#主機發送 url 請求到 port 5566 並響應
[root@regy ~]# curl localhost:5566
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>

#進入 Nginx Container
[root@regy ~]# docker exec -it nginx-01 /bin/bash
root@2e25e782345b:/# whereis nginx
nginx: /usr/sbin/nginx /usr/lib/nginx /etc/nginx /usr/share/nginx
root@2e25e782345b:/# cd /etc/nginx
root@2e25e782345b:/etc/nginx# ls
conf.d          koi-utf  mime.types  nginx.conf   uwsgi_params
fastcgi_params  koi-win  modules     scgi_params  win-utf
```

![img/Untitled%205.png](img/Untitled%205.png)

> 每次改動 Nginx 配置文件，都需要進入容器內部。可以在容器外部提供一個映射路徑，達到在容器外部修改文件，容器內部就可以自動修改  -v 數據卷
> 

# Deploy Tomcat

- 下載 tomcat 鏡像並啟動測試

```bash
#-rm 用於測試用，用完即刪
docker run -it --rm tomcat:9.0
```

- 後台啟動運行

```bash
docker run -d -p 1115:8080 --name tomcat-01 tomcat:9
```

- 進入容器

```bash
docker exec -it tomcat-01 /bin/bash

#linux command 少了
#webapps 沒文件
#默認最小可運行環境
```

> 可通過容器外映射路徑，在外部放置項目自動同步到容器內部
> 

# Depoly ES and Kibana

- ES 暴露的端口多
- 十分消耗記憶體
- ES 的數據一般放到安全目錄

- 啟動ES

```bash
#--net somenetwork 網路配置

#啟動ES
docker run -d --name elasticsearch  -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.6.2
```

- 修改內存限制

```bash
#-e 修改ES環境配置
docker run -d --name elasticsearch  -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -e ES_JAVA_OPTS="-Xms64m -Xms512m" elasticsearch:7.6.2
```

> 使用 Kibana 如何連接到 ES？ — Docker 網絡原理
> 

# Docker Volume

## Why Need Volume

如果數據存在容器中，當容器刪除時數據就會丟失 → Requirement：Data Persistence

容器間可以實現數據共享的技術，Docker Container 產生的數據同步到 Local

將 Container Directory 掛載到 Linux

## How to Volume?

> 使用 Command 來掛載 -v
> 

```bash
docker run -it -v linux dir : container dir -p linux port : container port

# 匿名掛載
docker run -d -p --name nginx-ts -v /etc/nginx nginx

# 查看所有 volume 情況
docker volume ls

-v 容器內路徑 #匿名掛載
-v volume name : container path #具名掛載
-v /linux path : container path #指定路徑掛載

# extention
# ro: 路徑只能通過 linux 修改
docker run -d -p --name nginx-ts -v /etc/nginx:ro nginx # read only
docker run -d -p --name nginx-ts -v /etc/nginx:rw nginx # read write

# 掛載 volume container
docker run -d -p --name nginx-ts --volume-from container-name
```

> Dockerfile
> 

```bash

FROM centos:centos7

VOLUME ["volume","volume2"]

CMD echo "------end------"
CMD /bin/bash
```

# DockerFile

> Command argument script
> 

## Build Step

- create dockerfile
- docker build 構建成 image
- docker run image
- docker push (dockerHub)

## DockerFile Command

---

### FROM

格式為 `FROM <image>`或`FROM <image>:<tag>`

第一條指令必須為 `FROM` 指令。並且，如果在同一個Dockerfile中建立多個映像檔時，可以使用多個 `FROM` 指令（每個映像檔一次）

---

### MAINTAINER

格式為 `MAINTAINER <name>`，指定維護者訊息

---

### RUN

格式為 `RUN <command>` 或 `RUN ["executable", "param1", "param2"]`

前者將在 shell 終端中運行命令，即 `/bin/sh -c`；後者則使用 `exec` 執行。指定使用其它終端可以透過第二種方式實作，例如 `RUN ["/bin/bash", "-c", "echo hello"]`

每條 `RUN` 指令將在當前映像檔基底上執行指定命令，並產生新的映像檔。當命令較長時可以使用 `\` 來換行

---

### CMD

支援三種格式

- `CMD ["executable","param1","param2"]` 使用 `exec` 執行，推薦使用；
- `CMD command param1 param2` 在 `/bin/sh` 中執行，使用在給需要互動的指令；
- `CMD ["param1","param2"]` 提供給 `ENTRYPOINT` 的預設參數；

指定啟動容器時執行的命令，每個 Dockerfile 只能有一條 `CMD` 命令。如果指定了多條命令，只有最後一條會被執行

如果使用者啟動容器時候指定了運行的命令，則會覆蓋掉 `CMD` 指定的命令

---

### EXPOSE

格式為 `EXPOSE <port> [<port>...]`

設定 Docker 伺服器容器對外的埠號，供外界使用。在啟動容器時需要透過 -P，Docker 會自動分配一個埠號轉發到指定的埠號

---

### ENV

格式為 `ENV <key> <value>`。 指定一個環境變數，會被後續 `RUN` 指令使用，並在容器運行時保持

例如

```bash
ENV PG_MAJOR 9.3
ENV PG_VERSION 9.3.4
RUN curl -SL http://example.com/postgres-$PG_VERSION.tar.xz | tar -xJC /usr/src/postgress && …
ENV PATH /usr/local/postgres-$PG_MAJOR/bin:$PATH
```

---

### ADD

格式為 `ADD <src> <dest>`

該命令將複製指定的 `<src>` 到容器中的 `<dest>`。 其中 `<src>` 可以是 Dockerfile 所在目錄的相對路徑；也可以是一個 URL；還可以是一個 tar 檔案（其複製後會自動解壓縮）

---

### COPY

格式為 `COPY <src> <dest>`

複製本地端的 `<src>`（為 Dockefile 所在目錄的相對路徑）到容器中的 `<dest>`

當使用本地目錄為根目錄時，推薦使用 `COPY`

---

### ENTRYPOINT

兩種格式：

- `ENTRYPOINT ["executable", "param1", "param2"]`
- `ENTRYPOINT command param1 param2`（shell中執行）

指定容器啟動後執行的命令，並且不會被 `docker run` 提供的參數覆蓋

每個 Dockerfile 中只能有一個 `ENTRYPOINT`，當指定多個時，只有最後一個會生效

---

### VOLUME

格式為 `VOLUME ["/data"]`

建立一個可以從本地端或其他容器掛載的掛載點，一般用來存放資料庫和需要保存的資料等

---

### USER

格式為 `USER daemon`

指定運行容器時的使用者名稱或 UID，後續的 `RUN` 也會使用指定使用者

當服務不需要管理員權限時，可以透過該命令指定運行使用者。並且可以在之前建立所需要的使用者，例如：`RUN groupadd -r postgres && useradd -r -g postgres postgres`。要臨時取得管理員權限可以使用 `gosu`，而不推薦 `sudo`

---

### WORKDIR

格式為 `WORKDIR /path/to/workdir`

為後續的 `RUN`、`CMD`、`ENTRYPOINT` 指令指定工作目錄

可以使用多個 `WORKDIR` 指令，後續命令如果參數是相對路徑，則會基於之前命令指定的路徑。例如

```bash
WORKDIR /a
WORKDIR b
WORKDIR c
RUN pwd

```

則最終路徑為 `/a/b/c`

---

### ONBUILD

格式為 `ONBUILD [INSTRUCTION]`

指定當建立的映像檔作為其它新建立映像檔的基底映像檔時，所執行的操作指令

例如，Dockerfile 使用以下的內容建立了映像檔 `image-A`

```bash
[...]
ONBUILD ADD . /app/src
ONBUILD RUN /usr/local/bin/python-build --dir /app/src
[...]

```

如果基於 image-A 建立新的映像檔時，新的 Dockerfile 中使用 `FROM image-A`指定基底映像檔時，會自動執行 `ONBUILD` 指令內容，等於在後面新增了兩條指令

```docker
FROM image-A
#Automatically run the following
ADD . /app/src
RUN /usr/local/bin/python-build --dir /app/src

```

使用 `ONBUILD` 指令的映像檔，推薦在標籤中註明，例如 `ruby:1.9-onbuild`

## Testing

> 創建 centos image
> 

```bash
# 編寫 dockerfile 文件
[root@regy dockerfile-test]# cat dockerfileTest-centos
FROM centos:centos7
MAINTAINER regy<p714140432@gmail.com>

ENV MYPATH /usr/local
WORKDIR $MYPATH

RUN yum -y install vim
RUN yum -y install net-tools

EXPOSE 80

CMD echo $MYPATH
CMD echo "---end---"
CMD /bin/bash

# 通過 dockerfile 構建 image
[root@regy dockerfile-test]# docker build -f dockerfileTest-centos -t testcentos:0.1 .

#測試運行

```

# Docker Network

## Docker0

- 每啟動一個 docker container, docker 就會為 container 分配一個 ip, 安裝 docker 會有一個網卡 docker0 bridge mode, 使用 veth-pair 技術
- Docker through docker0 route internal network address
- Default, can not connect through domain name

## Custom Network

### Network Mode

- Bridge
- none
- host
- container(少)

```bash
docker network create --help

docker network create --driver bridge --subnet 192.168.0.0/16 --gateway 192.168.0.1 mongonet
```

# Docker Compose

- docker-compose.yml