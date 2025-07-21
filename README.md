# Simple Marketplace

**Тестовое задание на стажировку в VK**  
REST API маркетплейса: регистрация, авторизация, размещение и просмотр объявлений.

---

## Оглавление
- [Описание](#описание)
- [Технологии](#технологии)
- [Запуск (Docker)](#запуск-docker)
- [Структура API](#структура-api)
- [Документация API (Swagger)](#документация-api-swagger)
- [Примитивный фронтенд](#примитивный-фронтенд)
- [Тесты](#тесты)
- [Структура проекта](#структура-проекта)

---

## Описание
Реализован REST API для условного маркетплейса с возможностью:
- Регистрации и авторизации пользователей (JWT).
- Размещения объявлений (только для авторизованных).
- Просмотра ленты объявлений с фильтрацией, сортировкой, пагинацией.
- Минимальный фронтенд для тестирования API.

---

## Технологии
- **Go 1.24.5**
- **PostgreSQL**
- **Swagger (OpenAPI 3.0)**
- **Docker, docker-compose**
- **Чистая архитектура**

---

## Запуск (Docker)

1. **Склонируйте репозиторий:**
   ```sh
   git clone <repo_url>
   cd simple-marketplace
   ```

2. **Запустите сервисы:**
   ```sh
   docker compose up -d --build db db_test app
   ```
   
3. **Применить миграции к тестовой БД**
   ```sh
   docker compose run --rm migrate_test_db
   ```
   По умолчанию сервис будет доступен на [http://localhost:8080](http://localhost:8080)


4. **(Опционально) Примените миграции вручную:**
   ```sh
   docker-compose run --rm app migrate
   ```

---

## Структура API

- **POST /auth/register** — регистрация пользователя
- **POST /auth/login** — авторизация, получение JWT
- **GET /listings** — лента объявлений (фильтры: цена, сортировка, пагинация)
- **POST /listings** — размещение объявления (требует JWT)
- **GET /health** — проверка статуса сервера

**Подробная спецификация — в [docs/openapi.yaml](docs/openapi.yaml)**

---

## Документация API (Swagger)

- Файл спецификации: [`docs/openapi.yaml`](docs/openapi.yaml)
- Можно открыть в [Swagger Editor](https://editor.swagger.io/) или [Redoc](https://redocly.github.io/redoc/)

**Пример запуска Swagger UI через Docker:**
```sh
docker run -p 8081:8080 -v ${PWD}/docs/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui \
  -e SWAGGER_JSON=/openapi.yaml
```
Документация будет доступна на [http://localhost:8081](http://localhost:8081)

---

## Примитивный фронтенд

- Находится в папке `frontend/`
- Позволяет:
  - Зарегистрироваться, войти, разместить объявление
  - Смотреть ленту, фильтровать, сортировать, пагинировать
  - Массово создавать тестовые объявления
- Для тестирования API без Postman

---

## Тесты

- Юнит-тесты: покрывают бизнес-логику, валидацию, мапперы
- Интеграционные тесты: работа с репозиториями (Postgres)
- Запуск тестов:
  ```sh
  go test ./...
  ```

---

## Структура проекта

```
simple-marketplace/
  docs/               # Swagger-файл для документации
  frontend/           # Примитивный фронтенд (HTML/CSS/JS)
  internal/           # Основная бизнес-логика, delivery, domain, infra
  migrations/         # SQL-миграции для БД
  test/               # Юнит- и интеграционные тесты
  Dockerfile
  docker-compose.yml
  README.md
```