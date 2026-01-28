# telemetry

gRPC сервис приема телеметрии автомобилей. Принимает данные, обновляет актуальное состояние в Redis и публикует события в Kafka.

## Ответственность

- прием телеметрии по gRPC
- запись текущего состояния в Redis
- публикация телеметрии/нарушений в Kafka

## Запуск локально

1) Подними Redis и Kafka.
2) Создай `telemetry/.env` из `telemetry/.env.example`.
3) Запусти:

```bash
go run ./cmd
```

## Запуск в Docker Compose

Используется общий `deployments/.env` и compose из `deployments/`.

## gRPC

Контракт описан в `proto/telemetry/telemetry.proto`, сгенерированный код лежит в `gen/proto/telemetry`.

## Kafka

Топики по умолчанию:

- `telemetry.raw`
- `telemetry.violations`

Имена настраиваются через `KAFKA_TOPIC_TELEMETRY_RAW` и `KAFKA_TOPIC_VIOLATIONS`.

## Переменные окружения

См. `telemetry/.env.example`. Ключевые:

- `GRPC_PORT`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_DB`
- `KAFKA_BROKERS`, `KAFKA_TOPIC_TELEMETRY_RAW`, `KAFKA_TOPIC_VIOLATIONS`
