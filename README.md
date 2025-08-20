# Checklist - Микросервисное ToDo приложение

Микросервисное приложение для управления задачами (ToDo List), построенное на Go с использованием современных технологий и архитектурных паттернов.

## 📋 Описание проекта

Checklist - это распределенное приложение для управления списком задач, состоящее из нескольких микросервисов:

- **API Service** - REST API для взаимодействия с пользователем
- **DB Service** - сервис для работы с базами данных
- **Kafka Service** - сервис для обработки событий и логирования 

### Основные компоненты:

- **API Service**: Принимает HTTP запросы и проксирует их в DB Service
- **DB Service**: Управляет данными в PostgreSQL, опционально кэширует в Redis
- **Kafka Service**: Обрабатывает события и логирует активность пользователей
- **PostgreSQL**: Основная база данных для хранения задач
- **Redis**: Кэш для часто используемых задач (опционально)
- **Kafka**: Message broker для асинхронной обработки событий (опционально)

## 📡 API Endpoints

### POST /create
Создание новой задачи
```json
{
  "title": "Заголовок задачи",
  "description": "Описание задачи"
}
```

### GET /list
Получение списка всех задач
```json
{
  "tasks": [
    {
      "id": "uuid",
      "title": "Заголовок",
      "description": "Описание",
      "completed": false,
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

### DELETE /delete
Удаление задачи
```json
{
  "id": "uuid"
}
```

### PUT /done
Отметка задачи как выполненной
```json
{
  "id": "uuid"
}
```

## 🛠️ Технический стек

- **Язык**: Go 1.21+
- **База данных**: PostgreSQL 15
- **Кэш**: Redis 7
- **Message Broker**: Apache Kafka
- **Контейнеризация**: Docker & Docker Compose
- **Протоколы**: HTTP/REST, gRPC 
- **Тестирование**: Go testing, Testify
- **Линтинг**: golangci-lint

## 🚀 Быстрый старт

### Предварительные требования

- Docker & Docker Compose
- Make
- Go 1.21+ (для разработки)

### Установка и запуск

1. **Клонирование репозитория**
```bash
git clone https://github.com/yourusername/checklist.git
cd checklist
```

2. **Запуск всех сервисов**
```bash
make up
```

3. **Проверка работоспособности**
```bash
curl -X GET http://localhost:8080/list
```

### Основные команды

```bash
# Запуск всех сервисов
make up

# Остановка всех сервисов
make down

# Запуск тестов
make test

# Запуск unit тестов
make test-unit

# Запуск интеграционных тестов
make test-integration

# Запуск линтера
make lint

# Просмотр логов
make logs

# Сборка всех сервисов
make build

# Очистка
make clean

# Перезапуск с пересборкой
make restart
```

## 🧪 Тестирование

Проект включает в себя два типа тестов:

### Unit тесты
```bash
make test-unit
```

### Интеграционные тесты
```bash
make test-integration
```

### Покрытие тестами
```bash
make coverage
```

## 📊 Мониторинг и логирование

### Логи
Все сервисы логируют свою активность в файлы:
- `logs/api-service.log` - логи API сервиса
- `logs/db-service.log` - логи DB сервиса  
- `logs/kafka-service.log` - логи Kafka сервиса

### Просмотр логов в реальном времени
```bash
# Все логи
make logs

# Логи конкретного сервиса
docker-compose logs -f api-service
```

## 🔄 CI/CD

Проект настроен для работы с GitHub Actions:

- Автоматический запуск тестов при push/PR
- Линтинг кода
- Сборка Docker образов
- Проверка безопасности

