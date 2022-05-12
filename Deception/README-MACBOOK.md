# Выполняем лабораторную работу Deception на Макбуке

Это заметка к лабораторной работе для тех, кто выполняет её на Макбуке, не используя виртуалки или серверы. На момент написания инструкции удобного варианта запустить виртуальную машину на Макбуке на процессоре M1 нет.

Файл активности хакера скомпилирован под линукс, поэтому не запускается на любых нелинукс-машинках, это касается и макбуков. Поэтому файл надо запустить на линуксе, для этого удобно воспользоваться контейнером – обернуть файл в линукс. Воспользуемся базовым образом Убунту, положим туда файл активности хакера и запустим.

Для взаимодействия файла активности с honeypot'ом нужно изменить настройки сети – присоединить контейнер с файлом активности хакера к сети хоста (докера), оттуда будут видны айпишник и порты honeypot'а.

Используйте эту команду вместо запуска файла `./hackersActivity`. Чтоб создать и запустить контейнер с файлом активности хакера на Маке с Intel, нужно запустить эту команду:

```
docker run -it --rm --network host -v "$(pwd)"/hackersActivity:/app/hackersActivity ubuntu /app/hackersActivity
```

На Маке с M1:

```
docker run -it --rm --platform linux/amd64 --network host -v "$(pwd)"/hackersActivity:/app/hackersActivity ubuntu /app/hackersActivity
```

## Альтернативное решение через отредактированный docker-compose файл

Также вместо этого можно воспользоваться отредактированным файлом docker-compose. В нём будут заданы настройки для honeypot-контейнера и контейнера активности хакера.

Код файла для Макбуков c Intel, используйте его вместо исходного docker-compose:

```
version: "3.7"
services:

  honeypot:
    image: honeytrap/honeytrap:latest
    ports:
      - 22:8022
      - 23:8023
      - 80:8080
      - 21:8021
      - 6379:6379
      - 27016:27016
    volumes:
      - type: bind
        source: ./Honeypot/config/config.toml
        target: /config/config.toml
  hackers_activity:
    image: ubuntu
    command: /app/hackersActivity
    volumes:
      - type: bind
        source: ./hackersActivity
        target: /app/hackersActivity
    network_mode: host
```

Для Макбуков с M1:

```
version: "3.7"
services:

  honeypot:
    image: honeytrap/honeytrap:latest
    ports:
      - 22:8022
      - 23:8023
      - 80:8080
      - 21:8021
      - 6379:6379
      - 27016:27016
    volumes:
      - type: bind
        source: ./Honeypot/config/config.toml
        target: /config/config.toml
  hackers_activity:
    image: ubuntu
    platform: linux/amd64
    command: /app/hackersActivity
    volumes:
      - type: bind
        source: ./hackersActivity
        target: /app/hackersActivity
    network_mode: host
```

Технически это то же самое решение, что и запуск команды, но реализация более удобная. Для запуска honeypot-контейнера используйте команду:

```
docker-compose up honeypot
```

Для запуска контейнера с файлом активности хакера:

```
docker-compose up hackers_activity
```

Спасибо [Кириллу Гашкову](https://github.com/kirillgashkov) за помощь с решением.
