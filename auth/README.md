# auth

gRPC сервис аутентификации и регистрации пользователей CityDrive. Хранит пользователей в PostgreSQL и выдает JWT.

## Ответственность

- регистрация пользователя
- логин пользователя
- выпуск JWT
- работа с PostgreSQL

## Запуск локально

1) Подними PostgreSQL.
2) Создай `auth/.env` из `auth/.env.example`.
3) Запусти:

```bash
go run ./cmd
```

## Запуск в Docker Compose

Используется общий `deployments/.env` и compose из `deployments/`. Сервис стартует после применения миграций контейнером `migrate`.

## gRPC

Контракт описан в `proto/auth/auth.proto`, сгенерированный код лежит в `gen/proto/auth`.

## Переменные окружения

См. `auth/.env.example`. Ключевые:

- `GRPC_PORT`
- `DB_URL` / `DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER` / `DB_PASSWORD`
- `JWT_SECRET_KEY`, `JWT_ALG`, `JWT_EXPIRATION`
