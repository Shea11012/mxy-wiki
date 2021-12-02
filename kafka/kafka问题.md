### 从宿主机连接访问容器内的kafka

docker compose file 文件配置

```dockerfile
kafka1:
    image: wurstmeister/kafka:2.13-2.7.0
    container_name: kafka1
    ports:
      - '9093:9093' # 9093 外部访问，9092 内部访问
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_LISTENERS: INTERNAL://0.0.0.0:9092,EXTERNAL://0.0.0.0:9093	# external地址必须设置为 0.0.0.0
      KAFKA_ADVERTISED_LISTENERS: INTERNAL://kafka1:9092,EXTERNAL://localhost:9093	# 关键是要配置第二个地址变为 localhost
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: INTERNAL
```

