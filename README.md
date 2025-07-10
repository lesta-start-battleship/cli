# Морской бой CLI

CLI-приложение для игры "Морской бой". Поддерживает регистрацию, авторизацию, гильдии, чат, матчмейкинг и магазин через текстовый интерфейс.

## Основные возможности

- Регистрация и вход (`register`, `login`, `oauth`).
- Управление гильдиями и чат (`guild join`, `guild chat`).
- Поиск соперников (`matchmake`).
- Магазин и лидеры (`shop open`, `scoreboard`).

## Команда

- Евгений: API-клиенты.
- Сергей: Интерфейс CLI.
- Имиль: Чат, матчмейкинг.

## Запуск приложения:

```bash

go run cmd/main.go

```

## Запуск через Docker:

```bash

  docker build -f ./build/cli.dockerfile -t "lesta-battleship-cli:dev" .

  docker run -it "lesta-battleship-cli:dev"
```