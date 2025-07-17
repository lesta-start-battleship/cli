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

## Запуск через justfile:

```bash
  just setup-and-build          # Сборка приложения
  just run                      # Запуск приложения
```

Перед использованием убедитесь, что у вас установлен [менеджер just](https://github.com/casey/just)