###### tags: `Практика`
# Фишинг
[**Фи́шинг**](https://ru.wikipedia.org/wiki/%D0%A4%D0%B8%D1%88%D0%B8%D0%BD%D0%B3) (англ. phishing от fishing «рыбная ловля, выуживание») — вид интернет-мошенничества, целью которого является получение доступа к конфиденциальным данным пользователей — логинам и паролям.

## Используемые инструменты
* [SMTP-сервер](https://ru.wikipedia.org/wiki/SMTP) (Почтовый сервер)
* [Gophish](https://getgophish.com/)
## Инструкция по установке
### Загрузка
```
git clone https://github.com/Ivanhahanov/InformtionsSecurityMethodsAndTools.git
```


### Развертывание
```
docker-compose up -d
```

```
docker-compose logs -f
```

### Конфигурация
#### Создание сертификатов
```
docker run -ti --rm -v "$(pwd)"/config/ssl:/tmp/docker-mailserver/ssl -h mail.domain.com -t tvial/docker-mailserver generate-ssl-certificate
```
#### Создание пользователей
```
# сделать скрипт исполняемым
chmod +x setup.sh

# создание email аккаунтов user/password
./setup.sh -i tvial/docker-mailserver:latest email add admin@domain.com admin123
./setup.sh -i tvial/docker-mailserver:latest email add user@domain.com user123

# список email аккаунтов 
./setup.sh -i tvial/docker-mailserver:latest email list
```
#### Проверка отправки сообщений
Для проверки работы Почтового сервера рекомендуется использовать утилиту [swaks](https://github.com/jetmore/swaks)
```
# Установка утилиты 
sudo apt-get install swaks

# Отправка email для user@domain.com
swaks --from admin@domain.com --to user@domain.com --server <hostIP>:587 -tlso -au admin@domain.com -ap user123 --header "Subject: test from admin" --body "testing 123"

# Отправка email для admin@domain.com
swaks --from user@domain.com --to admin@domain.com --server <hostIP>:587 -tlso -au user@domain.com -ap admin123 --header "Subject: test from user" --body "testing 456"
```
## Инструкция по выполнению

## Инструкция по тестирование
```
как проверить
```
## Отчёт
как сдать?
* список требований
