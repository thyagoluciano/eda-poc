## Setup para usuários de windows

1) Instalar o compilador gcc 64 bits, que é uma dependência da lib  gopkg.in/confluentinc/confluent-kafka-go.v1/kafka utilizado na classe infra\kafka\producer\confluent.go. Sugestão: seguir o manual do site https://www.msys2.org/
2) Assim que concluído o passo anterior, adicione o caminho da sua instalação conforme exemplo abaixo ao Path das variáveis de ambiente
    ```
    C:\msys64\mingw64\bin
    ```
3) Execute o comando abaixo para fazer o download de todas as dependências:
    ```
    go mod tidy
    ```
## Executando o projeto
1) Todos os containers necessários para a execução do projeto estão no diretório /docker do projeto. Acesse o diretório /docker e suba os containers com o comando:
    ```
    docker-compose up
     ```
2) Assim que todos os containers estiverem UP, executar o comando abaixo para conectar ao ksqldb e na sequência, executar os scripts da sessão "Scripts Kafka":
    ```
    docker exec -it ksqldb-cli ksql http://ksqldb-server:8088
    ```
3) Executando os scripts de criação das streams, executar o projeto:
    ```
    go run main.go
    ```

### Scripts Kafka

```
SET 'auto.offset.reset' = 'earliest';
```

```
CREATE STREAM CUSTOMER_ORDER_STREAM (
  eventId VARCHAR,
  eventType VARCHAR,
  metadata STRUCT<
    id VARCHAR,
    domain VARCHAR,
    externalId VARCHAR,
    customerId VARCHAR,
    context STRUCT<
      spanId VARCHAR,
      traceId VARCHAR,
      organization VARCHAR,
      application VARCHAR,
      channel VARCHAR,
      custom VARCHAR>,
    timestamp VARCHAR>,
  payload VARCHAR
) WITH ( 
  KAFKA_TOPIC = 'CUSTOMER_ORDER', 
  VALUE_FORMAT='JSON', 
  partitions=5
);
```

```
CREATE STREAM IF NOT EXISTS PURCHASE_ORDER_CHECKED_OUT WITH (partitions=10) AS SELECT eventId, eventType, metadata, payload FROM CUSTOMER_ORDER_STREAM WHERE eventType = 'PURCHASE_ORDER_CHECKED_OUT' EMIT CHANGES;
CREATE STREAM IF NOT EXISTS CUSTOMER_ORDER_CREATED WITH (partitions=10) AS SELECT eventId, eventType, metadata, payload FROM CUSTOMER_ORDER_STREAM WHERE eventType = 'CUSTOMER_ORDER_CREATED' EMIT CHANGES;
CREATE STREAM IF NOT EXISTS CUSTOMER_ORDER_PAYMENT_COMPLETED WITH (partitions=10) AS SELECT eventId, eventType, metadata, payload FROM CUSTOMER_ORDER_STREAM WHERE eventType = 'CUSTOMER_ORDER_PAYMENT_COMPLETED' EMIT CHANGES;
CREATE STREAM IF NOT EXISTS CUSTOMER_ORDER_INVENTORY_COMPLETED WITH (partitions=10) AS SELECT eventId, eventType, metadata, payload FROM CUSTOMER_ORDER_STREAM WHERE eventType = 'CUSTOMER_ORDER_INVENTORY_COMPLETED' EMIT CHANGES;
CREATE STREAM IF NOT EXISTS CUSTOMER_ORDER_SHIPPING_COMPLETED WITH (partitions=10) AS SELECT eventId, eventType, metadata, payload FROM CUSTOMER_ORDER_STREAM WHERE eventType = 'CUSTOMER_ORDER_SHIPPING_COMPLETED' EMIT CHANGES;
```

```
CREATE TABLE IF NOT EXISTS CUSTOMER_ORDER_AGGREGATION WITH (
    PARTITIONS=10,
    REPLICAS=1,
    VALUE_FORMAT='JSON'
) AS SELECT
  metadata->id AS CUSTOMER_ORDER_ID,
  COLLECT_LIST('{"eventId": "' + eventId + '", "eventType": "' + eventType + '", "payload": ' + payload + '}') AS DATA_AGG
FROM CUSTOMER_ORDER_STREAM
GROUP BY metadata->id
EMIT CHANGES;

CREATE TABLE CUSTOMER_ORDER_AGGREGATION_5 WITH (
  PARTITIONS=10,
  REPLICAS=1,
  VALUE_FORMAT='JSON'
) AS SELECT
  metadata->id AS CUSTOMER_ORDER_ID,
  COLLECT_LIST(eventType) AS TYPE_LIST,
  COLLECT_LIST(payload) AS JSON_DATA_LIST
FROM CUSTOMER_ORDER_STREAM
GROUP BY metadata->id
EMIT CHANGES;
```