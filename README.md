# CryptocoinPrice

HTTP API для сбора, хранения и отображения стоимости криптовалют.

## Предварительная настройка (перед запуском)

Склонировать репозиторий и перейти в папку с репозиторием

```shell
git clone https://github.com/ej-you/CryptocoinPrice.git
cd ./CryptocoinPrice
```

В корне проекта необходимо создать файл .env и поместить в него переменные окружения:

```dotenv
# пример настроек БД
POSTGRES_USER="cryptoprice"
POSTGRES_PASSWORD="p@SSw0rd"
POSTGRES_DB="cryptoprice_db"
PGDATA="/var/lib/postgresql/data/pgdata"

POSTGRES_HOST="postgresql"
POSTGRES_PORT="5432"

# создан специально на период проверки задания, чтобы не нужно было возиться,
# но по хорошему ключи не хранят в открытом доступе
COINGECKO_API_KEY="CG-tnq5533ANABFRGWBcws1ewEV"
```

## Запуск

```shell
docker compose -f ./docker-compose.yml up -d
```

> !! Затем необходимо провести миграции, используя команду !!
>
> ```shell
> docker compose -f ./docker-compose.yml exec server sh -c "/app/migrator up"
> ```

По умолчанию сервер запускается на `8000` порту.

Swagger документация — `/api/v1/docs`.

[Ссылка](http://127.0.0.1:8000/api/v1/docs) на swagger документацию (локальный адрес).

Чтобы остановить сервер используйте

```shell
docker compose -f ./docker-compose.yml down
```

## Дополнительно

### Логирование

У приложения можно настроить уровень логирования и формат логов.

Доступные форматы логов:

1. text (по умолчанию)
2. json

Доступные уровни логирования

1. info (по умолчанию)
2. warn
3. error

Настройка происходит с помощью переменных окружения. Пример:

```dotenv
LOG_FORMAT=text
LOG_LEVEL=error
```

### Настройка интервала обновления

Интервал сбора новых цен отслеживаемых монет настраивается
через переменную окружения `PRICE_COLLECT_INTERVAL`.
По умолчанию он равер 5 секундам. Пример настройки для 30 секунд:

```dotenv
PRICE_COLLECT_INTERVAL=30s
```

### Особенности работы программы

Время везде использует часовой пояс `UTC`. Это единый часовой пояс.
Предполагается перевод времени в локальный(ые) часовой(ые) пояс(а) на клиенте.
