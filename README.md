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

