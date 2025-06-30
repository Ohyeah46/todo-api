## Содержание

- [Требования](#требования)
- [Установка](#установка)
- [Настройка базы данных PostgreSQL](#настройка-базы-данных-postgresql)
- [Запуск проекта](#запуск-проекта)
- [API эндпоинты](#api-эндпоинты)
- [Особенности](#особенности)
- [Пример запросов](#пример-запросов)
- [Полезные советы](#полезные-советы) 

## Требования

- Go 1.20 или новее
- PostgreSQL 13 или новее
- Любой HTTP клиент (Postman, curl и т.п.)
## Установка

1. Клонировать репозиторий:
```bash
git clone https://github.com/Ohyeah46/todo-api.git
cd todo-api

go mod tidy
```

2. Настройка базы данных PostgreSQL
```bash
CREATE DATABASE todo;
CREATE USER postgres WITH PASSWORD '1111';
GRANT ALL PRIVILEGES ON DATABASE todo TO postgres;
```
3. Настройки подключения к БД в файле database/database.go
```bash
dsn := "host=localhost user=postgres password=1111 dbname=todo port=5432 sslmode=disable"
```
4. Запуск проекта
```bash 
go build -o todo-api
./todo-api

go run main.go
```

5. API эндпоинты
```bash
Публичные (без авторизации)
POST /register — регистрация пользователя (JSON: username, password)

POST /login — вход, возвращает JWT токен (JSON: username, password)

GET /debug/slice — тест слайсов (демо)

GET /debug/map — тест мап (демо)

GET /async-example — демонстрация каналов и контекста

Защищённые (JWT авторизация обязательна)
POST /tasks — создать новую задачу

GET /tasks — получить список задач текущего пользователя

GET /tasks/:id — получить задачу по ID

PUT /tasks/:id — обновить задачу

DELETE /tasks/:id — удалить задачу

POST /enqueue — поставить задачу в очередь на асинхронную обработку
```
6. Особенности реализации

- JWT аутентификация:
- Пользователь получает JWT при входе, который надо передавать в заголовке Authorization: Bearer <token>.

- Асинхронная обработка:
- Есть очередь задач и фоновые горутины-воркеры, которые обрабатывают задачи асинхронно (логируют выполнение с задержкой).

- GORM используется как ORM для работы с PostgreSQL.

- Middleware Gin обеспечивает проверку JWT для защищённых маршрутов.