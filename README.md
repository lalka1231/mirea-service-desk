# MIREA Service Desk

Сервис для регистрации и обработки заявок на обслуживание инфраструктуры.

## Стек
- Backend: Go + Gin + JWT
- Frontend: HTML/CSS/JavaScript
- DB: SQLite
- DevOps: Docker, Docker Compose, GitHub Actions

## Локальный запуск без Docker
```bash
cd backend
cp .env.example .env
go mod download
go run ./cmd/server
```

Открой `frontend/index.html` в браузере или раздай папку `frontend` любым статическим сервером.

## Запуск через Docker Compose
```bash
docker compose up --build
```

После запуска:
- backend: `http://localhost:8080`
- frontend: `http://localhost`

## Переменные окружения
Пример находится в `backend/.env.example`:
- `PORT` — порт backend
- `DB_PATH` — путь к SQLite базе
- `JWT_SECRET` — секрет для JWT

## CI/CD
Workflow находится в `.github/workflows/ci-cd.yml` и выполняет:
- проверку форматирования Go-кода;
- сборку backend;
- запуск тестов;
- сборку Docker-образов backend и frontend.
