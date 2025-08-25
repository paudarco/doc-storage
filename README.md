# Doc storage

Кеширующий веб-сервер для документов.

## Возможности

*   Регистрация/Аутентификация пользователей.
*   Загрузка, список, получение, удаление документов (файлы/JSON).
*   Управление доступом (публичные, приватные, по логинам).
*   Кэширование в Redis для высокой нагрузки на чтение.

## Технологии

*   **Go 1.21+, Gin, pgx/v5, redis/go-redis/v9, Docker**

## Быстрый старт

1.  Клонируйте репозиторий.
2.  Настройте `.env` (пример в .env.example).
3.  `docker-compose up --build`

Сервис доступен на `http://localhost:8080`.

## API

*   `POST /api/register` (Требует `ADMIN_TOKEN`)
*   `POST /api/auth`
*   `POST /api/docs`
*   `GET/HEAD /api/docs[?login=&key=&value=&limit=]`
*   `GET/HEAD /api/docs/:id`
*   `DELETE /api/docs/:id`
*   `DELETE /api/auth/:token`