# api-gateway

HTTP шлюз CityDrive. Принимает запросы клиентов и вызывает gRPC сервисы `auth`, `telemetry`, `admin`. Делает базовую валидацию JWT и прокидывает trace id.

## Ответственность

- HTTP API для внешних клиентов
- gRPC клиенты к внутренним сервисам
- JWT проверка для пользовательских и машинных токенов

## Запуск локально

1) Создай `api-gateway/.env` из шаблона `api-gateway/.env.example`.
2) Убедись, что gRPC сервисы подняты локально на портах `50051/50052/50053`.
3) Запусти:

```bash
go run ./cmd
```

## Запуск в Docker Compose

Используется общий `deployments/.env` и compose из `deployments/`.

## HTTP эндпоинты

- `GET /health`
- `POST /v1/user/login`
- `POST /v1/user/register`
- `PUT /api/v1/car-info`
- `GET /api/v1/cars/now`
- `GET /api/v1/cars/:id`
- `GET /api/v1/cars/history`
- `GET /api/v1/cars/:id/history`

## Переменные окружения

См. `api-gateway/.env.example`. Ключевые:

- `HTTP_PORT`
- `AUTH_GRPC_ADDR`, `TELEMETRY_GRPC_ADDR`, `ADMIN_GRPC_ADDR`
- `JWT_SECRET_KEY`, `JWT_CAR_SECRET_KEY`
