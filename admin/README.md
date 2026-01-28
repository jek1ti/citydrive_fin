# admin

gRPC сервис админских запросов CityDrive: выдача текущих координат/состояний и истории телеметрии.

## Ответственность

- текущие данные по автомобилям (из Redis/БД)
- детальная карточка автомобиля
- история телеметрии за период (из PostgreSQL)

## Запуск локально

1) Подними PostgreSQL и Redis.
2) Создай `admin/.env` из `admin/.env.example`.
3) Запусти:

```bash
go run ./cmd
```

## Запуск в Docker Compose

Используется общий `deployments/.env` и compose из `deployments/`. Сервис стартует после применения миграций контейнером `migrate`.

## gRPC

Контракт описан в `proto/admin/admin.proto`, сгенерированный код лежит в `gen/proto/admin`.

## Переменные окружения

См. `admin/.env.example`. Ключевые:

- `GRPC_PORT`
- `DB_URL`
- `REDIS_URL` (или `REDIS_HOST`/`REDIS_PORT` при доработке)
