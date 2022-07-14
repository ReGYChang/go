# Set up Debeizum Connector for Oracle

## Run docker compose

docker-compose-oracle.yaml
```yaml
version: '2'
services:
  zookeeper:
    image: quay.io/debezium/zookeeper:${DEBEZIUM_VERSION}
    ports:
      - 2181:2181
      - 2888:2888
      - 3888:3888
  kafka:
    image: quay.io/debezium/kafka:${DEBEZIUM_VERSION}
    ports:
      - 9092:9092
    links:
      - zookeeper
    environment:
      - ZOOKEEPER_CONNECT=zookeeper:2181
  schema-registry:
    image: confluentinc/cp-schema-registry:7.0.1
    ports:
      - 8181:8181
      - 8081:8081
    environment:
      - SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS=kafka:9092
      - SCHEMA_REGISTRY_HOST_NAME=schema-registry
      - SCHEMA_REGISTRY_LISTENERS=http://schema-registry:8081
    links:
      - zookeeper
  connect:
    image: quay.io/debezium/connect:${DEBEZIUM_VERSION}
    ports:
      - 8083:8083
    links:
      - kafka
      - schema-registry
    environment:
      - BOOTSTRAP_SERVERS=kafka:9092
      - GROUP_ID=1
      - CONFIG_STORAGE_TOPIC=my_connect_configs
      - OFFSET_STORAGE_TOPIC=my_connect_offsets
      - STATUS_STORAGE_TOPIC=my_connect_statuses
      - CONNECT_KEY_CONVERTER_SCHEMA_REGISTRY_URL=http://schema-registry:8081
      - CONNECT_VALUE_CONVERTER_SCHEMA_REGISTRY_URL=http://schema-registry:8081
  kafkaui:
    image: provectuslabs/kafka-ui:latest
    ports:
      - 8811:8080
    links:
      - kafka
    environment:
      - KAFKA_CLUSTERS_0_NAME=test
      - KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka:9092
```

```shell
export DEBEZIUM_VERSION=1.9
docker-compose -f docker-compose-oracle.yaml up
```

## Obtaining the Oracle JDBC driver file

```shell
wget https://repo1.maven.org/maven2/com/oracle/database/jdbc/ojdbc8/19.3.0.0/ojdbc8-19.3.0.0.jar
mv ojdbc8-19.3.0.0.jar ojdbc8.jar
docker cp ojdbc8.jar connect:/kafka/libs
docker restart connect
```

## Debezium Oracle connector configuration

```shell
vim oracle-debezium-connector.json
```

```json
{
  "name": "inventory-connector",
  "config": {
    "connector.class" : "io.debezium.connector.oracle.OracleConnector",
    "tasks.max" : "1",
    "database.server.name" : "oracleserver",
    "database.hostname" : "10.90.1.207",
    "database.port" : "1521",
    "database.user" : "logminer",
    "database.password" : "logminer",
    "database.dbname" : "EMESHY",
    "database.connection.adapter" : "logminer",
    "database.history.kafka.bootstrap.servers" : "kafka:9092",
    "database.history.kafka.topic": "schema-changes.test",
    "database.tablename.case.insensitive": "false",
    "log.mining.strategy": "online_catalog",
    "table.include.list": "EMESP.TP_SN_LOG"
  }
}
```

```shell
➜  curl -i -X POST -H "Accept:application/json" -H  "Content-Type:application/json" http://localhost:8083/connectors/ -d @oracle-debezium-connector.json
HTTP/1.1 201 Created
Date: Thu, 14 Jul 2022 05:09:35 GMT
Location: http://localhost:8083/connectors/inventory-connector
Content-Type: application/json
Content-Length: 632
Server: Jetty(9.4.44.v20210927)

{"name":"inventory-connector","config":{"connector.class":"io.debezium.connector.oracle.OracleConnector","tasks.max":"1","database.server.name":"oracleserver","database.hostname":"10.90.1.207","database.port":"1521","database.user":"logminer","database.password":"logminer","database.dbname":"EMESHY","database.connection.adapter":"logminer","database.history.kafka.bootstrap.servers":"kafka:9092","database.history.kafka.topic":"schema-changes.test","database.tablename.case.insensitive":"false","log.mining.strategy":"online_catalog","table.include.list":"EMESP.TP_SN_LOG","name":"inventory-connector"},"tasks":[],"type":"source"}% 
```

```shell
➜  curl -i  http://localhost:8083/connectors/inventory-connector
HTTP/1.1 200 OK
Date: Thu, 14 Jul 2022 05:36:56 GMT
Content-Type: application/json
Content-Length: 676
Server: Jetty(9.4.44.v20210927)

{"name":"inventory-connector","config":{"connector.class":"io.debezium.connector.oracle.OracleConnector","database.user":"logminer","database.dbname":"EMESHY","tasks.max":"1","database.connection.adapter":"logminer","database.history.kafka.bootstrap.servers":"kafka:9092","database.history.kafka.topic":"schema-changes.test","database.server.name":"oracleserver","database.tablename.case.insensitive":"false","log.mining.strategy":"online_catalog","database.port":"1521","database.hostname":"10.90.1.207","database.password":"logminer","name":"inventory-connector","table.include.list":"EMESP.TP_SN_LOG"},"tasks":[{"connector":"inventory-connector","task":0}],"type":"source"}% 
```

## Open kafkaui dashboard

```shell
http://localhost:8811
```

## Consume Kafka Topic Messages

```shell
docker exec -it connect bash

bin/kafka-console-consumer.sh  --bootstrap-server kafka:9092 --topic oracleserver --from-beginning
```