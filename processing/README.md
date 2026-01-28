# processing

Сервис обработки телеметрии из Kafka. Читает сообщения, обновляет актуальное состояние в Redis и сохраняет историю в PostgreSQL.

## Ответственность

- чтение телеметрии из Kafka consumer group
- запись текущего состояния в Redis
- сохранение истории телеметрии в PostgreSQL
- HTTP health endpoints

## Запуск локально

1) Подними PostgreSQL, Redis и Kafka.
2) Создай `processing/.env` из `processing/.env.example`.
3) Запусти:

```bash
go run ./cmd
```

## Health endpoints

- `GET /health/liveness`
- `GET /health/readiness`

## Kafka

Топик телеметрии настраивается через `KAFKA_TOPIC_TELEMETRY_RAW` (есть fallback на `KAFKA_TOPIC_TELEMETRY`).

## Переменные окружения

См. `processing/.env.example`. Ключевые:

- `HTTP_PORT`
- `DB_URL` и параметры БД
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_DB`
- `KAFKA_BROKERS`, `KAFKA_CONSUMER_GROUP_ID`, `KAFKA_TOPIC_TELEMETRY_RAW`
