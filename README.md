## CDEK Gamification Platform

Demo-проект платформы геймификации для внутреннего портала CDEK.

### Стек

- React + Vite
- BFF на Go + Fiber
- gRPC между BFF и внутренними сервисами
- clean architecture по слоям `presentation -> application -> domain <- infrastructure`
- Prometheus metrics через `promhttp`
- OpenTelemetry hooks
- Dockerfiles для сервисов

### Структура

- `frontend` - React SPA
- `bff` - REST BFF сервис
- `services/user-service` - сервис пользовательских данных
- `services/gamification` - сервис геймификации
- `shared/contracts` - общие gRPC-контракты для демо
- `proto` - proto-описания контрактов

### Что уже реализовано

- экраны `Дашборд`, `Профиль`, `Задания`, `Достижения`, `Рейтинг`, `Магазин наград`, `Покупки`
- BFF-агрегация данных из двух gRPC-сервисов
- demo JWT через `Bearer demo-token`
- рабочие действия для показа на демо:
  - принять задание
  - продвинуть задание
  - завершить задание с начислением XP
  - открыть достижение
  - обменять баллы на награду
- fallback фронта в локальный demo-режим, если BFF временно недоступен

### Запуск фронта локально

```bash
cd frontend
npm ci
npm run dev
```

Фронт по умолчанию ожидает BFF на `http://localhost:8080`.

### Запуск backend локально

На этой машине Go доступен по пути `C:\Program Files\Go\bin\go.exe`.

Запусти сервисы в отдельных терминалах:

```powershell
cd services/user-service
& "C:\Program Files\Go\bin\go.exe" run ./cmd/user-service
```

```powershell
cd services/gamification
& "C:\Program Files\Go\bin\go.exe" run ./cmd/gamification
```

```powershell
cd bff
& "C:\Program Files\Go\bin\go.exe" run ./cmd/bff
```

### Docker Compose

```bash
docker compose up --build
```

После старта:

- frontend: `http://localhost:4173`
- bff: `http://localhost:8080`
- metrics: `http://localhost:9090`

### Демо-сценарий

1. Открыть `Дашборд` и показать общий прогресс пользователя.
2. Перейти в `Задания` и принять или продвинуть одно из активных заданий.
3. Показать, что изменились XP, coins, история активности и профиль.
4. Перейти в `Достижения` и показать открытие эпического бейджа после выполнения нужного числа задач.
5. Перейти в `Рейтинг` и показать позицию пользователя.
6. Перейти в `Магазин наград` и обменять coins на награду.
7. Открыть `Покупки` и показать появившуюся запись.

### Контракты

- `proto/user/v1/user.proto`
- `proto/gamification/v1/gamification.proto`

### Архитектурная идея

- `Presentation`: React SPA и BFF
- `Application`: orchestration/use cases в сервисах
- `Domain`: сущности, статусы, ошибки, level logic
- `Infrastructure`: gRPC transport, memory-репозитории для демо, hooks для observability

Для показа архитектуры CDEK это уже даёт end-to-end сценарий. Следующий шаг после демо - заменить memory-репозитории на реализации `pgx` и `ruedis`, подключить Vault и вынести telemetry в полноценный OTEL pipeline.
