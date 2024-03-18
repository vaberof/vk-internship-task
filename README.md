# Тестовое задание на стажировку в VK Tech

## Использование

### Запуск окружения с приложением и БД PostgreSQL

    make docker.run

### Запуск юнит-тестов

    make tests.run

### Запуск юнит-тестов с отчетом о покрытии в файл *coverage.html*

    make tests.cover.report

## Описание задачи

Необходимо разработать бэкенд приложения "Фильмотека",
который предоставляет REST API для управления базой данных фильмов

## Описание реализации

### Общее

- Язык реализации - Go
- Для хранения данных используется PostgreSQL (при обращении в БД используется чистый SQL)
- Описана спецификация на API в формате Swagger 2.0 (доступна после запуска приложения по
  URL: http://localhost:8000/swagger)

### Функционал

**Основная часть:**

- Актеры:
    - Добавление информации об актёре
    - Частичное/полное изменение информации об актёре
    - Удаление информации об актёре
    - Получение списка актёров, для каждого актёра выдаётся также список фильмов с его участием
- Фильмы:
    - Добавление информации о фильме
    - Частичное/полное изменение информации о фильме
    - Удаление информации о фильме
    - Получение списка фильмов с возможностью сортировки по названию, по рейтингу, по дате выпуска. По умолчанию
      используется сортировка по рейтингу (по убыванию)
    - Поиск фильма по фрагменту названия, по фрагменту имени актёра
- API закрыт авторизацией:
    - поддерживаются две роли пользователей - обычный пользователь и администратор. Обычный пользователь имеет доступ
      только на получение данных и поиск, администратор - на все действия

**Бонусная часть**:

- Используется подход code-first (генерация спецификации из кода)
- HTTP сервер реализован с использованием стандартной библиотеки
- Логирование (в логах отображаются обрабатываемые HTTP запросы, ошибки)
- Код приложения покрыт юнит-тестами
- Dockerfile для сборки образа приложения
- docker-compose файл для запуска окружения с работающим приложением и СУБД PostgreSQL

### Технологии

Go, PostgreSQL, SQL, Swagger, Docker, Docker Compose