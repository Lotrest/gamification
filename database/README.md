# PostgreSQL

В проект добавлена реальная база `PostgreSQL` для локального запуска через `docker compose`.

## Что лежит в папке

- `init/001_create_schemas.sql` — схемы `user_service` и `gamification`
- `init/002_create_tables.sql` — таблицы для user-service и gamification-service
- `init/003_seed_demo_data.sql` — стартовые данные нового пользователя с нулевым прогрессом

## Как поднять базу

```powershell
docker compose up -d postgres
```

После запуска база будет доступна так:

- host: `localhost`
- port: `5432`
- database: `cdek_platform`
- user: `cdek`
- password: `cdek`

## Что уже создано

Схема `user_service`:

- `users`

Схема `gamification`:

- `user_state`
- `weekly_activity`
- `article_cards`
- `recent_activity`
- `tasks`
- `achievements`
- `leaderboard`
- `rewards`
- `purchases`
- `notifications`

## Важно

Сейчас сервисы все еще используют `memory`-репозитории. PostgreSQL уже подготовлен как реальная база проекта, но следующий шаг — переключить `user-service` и `gamification-service` на `pgx`-репозитории.
