# Демонстрационный сервис с Kafka, PostgreSQL, кешем

## Запуск проекта

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/ilam072/wbtech-l0.git
cd wbtech-l0
````

### 2. Переменные окружения

Создайте файл .env в корне проекта и заполните его по примеру ниже:

```bash
HTTP_PORT=8082

PGUSER=postgres
PGPASSWORD=postgres
PGHOST=localhost
PGPORT=5433
PGDATABASE=orders_service
PGSSLMODE=disable

KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=orders
KAFKA_GROUP_ID=order-consumer

CACHE_PRELOAD_LIMIT=1000
```

## Запуск приложения
Введите команду:
```bash
make up
```

## Тестирование

Чтобы запустить все тесты:

```bash
go test ./...
```

---

## Frontend
Перейдите по странице в браузере:
```
http://localhost:5500/
```

---

## Backend
Backend доступен по адресу:
```
http://localhost:8082/
```

Swagger-документация:
```
http://localhost:8082/swagger/
```

---

## Makefile
Полезные команды:
* `make up` — запустить приложние
* `make down` — остановить докер контейнеры
* `make producer` — запустить скрипт для отправки сообщений в Kafka
* `make topics` — посмотреть топики в Kafka
* `make messages` — посмотреть сообщения, отправленные в Kafka
