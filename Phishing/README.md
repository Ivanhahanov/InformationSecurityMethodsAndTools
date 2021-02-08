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
Mail server должен сообщить об ошибках, так как у нас нет пользователей и нет подписанного сертификата.
### Конфигурация
#### Создание сертификатов
Нам необходимо создать самоподписанный TLS сертификат. Такие сертификаты используются для шифрования данных, самый распространённый пример это HTTPS.

В нашем случае он нужен для зашифровки писем, чтобы никто не могу их прочитать при перехвате трафика

В контейнере **docker-mailserver** есть утилита для генерации сертификатов.

Запустим контейнер в интерактивном режиме (флаг **-it**).

Примонтируем раздел чтобы сертификат создался у нас в директории MailServer/config/ssl локально.
```
docker run -ti --rm -v "$(pwd)"/config/ssl:/tmp/docker-mailserver/ssl -h mail.domain.com -t tvial/docker-mailserver generate-ssl-certificate
```
Создание сертификата происходит в два этапа
1. создание CA сертификата
1. создание конечного сертификата

Вводимые поля должны быть идентичны кроме одного поля `Common Name (e.g. server FQDN or YOUR name)`. Здесь в первым случае надо указать **domain.com**, во втором **mail.domain.com**

Поля `extra` оставить пустыми

Пример генерации сертификата:

```shell
CA certificate filename (or enter to create)

Making CA certificate ...
====
openssl req  -new -keyout ./demoCA/private/cakey.pem -out ./demoCA/careq.pem 
Generating a RSA private key
...........................+++++
.................................................................................................................................+++++
writing new private key to './demoCA/private/cakey.pem'
Enter PEM pass phrase:
Verifying - Enter PEM pass phrase:
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:RU
State or Province Name (full name) [Some-State]:Moscow
Locality Name (eg, city) []:Moscow
Organization Name (eg, company) [Internet Widgits Pty Ltd]:My-company
Organizational Unit Name (eg, section) []:Red Team      
Common Name (e.g. server FQDN or YOUR name) []:domain.com
Email Address []:admin@domain.com

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
==> 0
====
====
openssl ca  -create_serial -out ./demoCA/cacert.pem -days 1095 -batch -keyfile ./demoCA/private/cakey.pem -selfsign -extensions v3_ca  -infiles ./demoCA/careq.pem
Using configuration from /usr/lib/ssl/openssl.cnf
Enter pass phrase for ./demoCA/private/cakey.pem:
Check that the request matches the signature
Signature ok
Certificate Details:
        Serial Number:
            5d:ac:b5:3f:cc:b6:26:bf:6a:20:dd:c1:ff:08:a0:be:14:1b:e4:3c
        Validity
            Not Before: Feb  8 13:40:27 2021 GMT
            Not After : Feb  8 13:40:27 2024 GMT
        Subject:
            countryName               = RU
            stateOrProvinceName       = Moscow
            organizationName          = My-company
            organizationalUnitName    = Red Team
            commonName                = domain.com
            emailAddress              = admin@domain.com
        X509v3 extensions:
            X509v3 Subject Key Identifier: 
                E9:4F:D1:CC:9D:09:14:A3:9C:23:68:8E:0E:76:9E:35:AE:22:56:61
            X509v3 Authority Key Identifier: 
                keyid:E9:4F:D1:CC:9D:09:14:A3:9C:23:68:8E:0E:76:9E:35:AE:22:56:61

            X509v3 Basic Constraints: critical
                CA:TRUE
Certificate is to be certified until Feb  8 13:40:27 2024 GMT (1095 days)

Write out database with 1 new entries
Data Base Updated
==> 0
====
CA certificate is in ./demoCA/cacert.pem
Ignoring -days; not generating a certificate
Generating a RSA private key
..................................................+++++
..............................................................+++++
writing new private key to '/tmp/docker-mailserver/ssl/mail.domain.com-key.pem'
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:RU
State or Province Name (full name) [Some-State]:Moscow
Locality Name (eg, city) []:Moscow
Organization Name (eg, company) [Internet Widgits Pty Ltd]:My-company
Organizational Unit Name (eg, section) []:Red Team       
Common Name (e.g. server FQDN or YOUR name) []:mail.domain.com
Email Address []:admin@domain.com

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
Using configuration from /usr/lib/ssl/openssl.cnf
Enter pass phrase for ./demoCA/private/cakey.pem:
Check that the request matches the signature
Signature ok
Certificate Details:
        Serial Number:
            5d:ac:b5:3f:cc:b6:26:bf:6a:20:dd:c1:ff:08:a0:be:14:1b:e4:3d
        Validity
            Not Before: Feb  8 13:41:55 2021 GMT
            Not After : Feb  8 13:41:55 2022 GMT
        Subject:
            countryName               = RU
            stateOrProvinceName       = Moscow
            organizationName          = My-company
            organizationalUnitName    = Red Team
            commonName                = mail.domain.com
            emailAddress              = admin@domain.com
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:FALSE
            Netscape Comment: 
                OpenSSL Generated Certificate
            X509v3 Subject Key Identifier: 
                2F:6B:63:BE:40:11:03:4E:EC:E7:27:9E:E7:F8:3B:A8:82:9C:84:D9
            X509v3 Authority Key Identifier: 
                keyid:E9:4F:D1:CC:9D:09:14:A3:9C:23:68:8E:0E:76:9E:35:AE:22:56:61

Certificate is to be certified until Feb  8 13:41:55 2022 GMT (365 days)
Sign the certificate? [y/n]:y


1 out of 1 certificate requests certified, commit? [y/n]y
Write out database with 1 new entries

```

#### Создание пользователей
Для того чтобы отправлять письма по электронной почте, нам нужны пользователи.

Для создания пользователей у нас есть Bash скрипт **setup.sh**
```
# сделать скрипт исполняемым
chmod +x setup.sh

# создание email аккаунтов user/password
./setup.sh -i tvial/docker-mailserver:latest email add admin@domain.com admin123
./setup.sh -i tvial/docker-mailserver:latest email add user@domain.com user123

# список email аккаунтов 
./setup.sh -i tvial/docker-mailserver:latest email list
```
В результате мы должны увидеть две почты admin@domain.com и user@domain.com
#### Проверка отправки сообщений
Для проверки работы Почтового сервера рекомендуется использовать утилиту [swaks](https://github.com/jetmore/swaks)
```
# Установка утилиты 
sudo apt-get install swaks

# Отправка email для user@domain.com
swaks --from admin@domain.com --to user@domain.com --server 127.0.0.1:587 -tlso -au admin@domain.com -ap admin123 --header "Subject: test from admin" --body "testing 123"

# Отправка email для admin@domain.com
swaks --from user@domain.com --to admin@domain.com --server 127.0.0.1:587 -tlso -au user@domain.com -ap user123 --header "Subject: test from user" --body "testing 456"
```
Описание флагов:
* `--from` - email отправителя
* `--to` - email получателя
* `--server` адрес SMTP сервера
* `--au` почта пользователя для авторизация
* `--ap` пароль пользователя для авторизация
* `--header` заголовок письма
* `--body` тело письма

В результате вы видим полный список действий которые произошли при отправке письма
```shell
=== Trying 127.0.0.1:587...
=== Connected to 127.0.0.1.
<-  220 mail.domain.com ESMTP
 -> EHLO ubuntu
<-  250-mail.domain.com
<-  250-PIPELINING
<-  250-SIZE 10240000
<-  250-ETRN
<-  250-STARTTLS
<-  250-ENHANCEDSTATUSCODES
<-  250-8BITMIME
<-  250-DSN
<-  250 CHUNKING
 -> STARTTLS
<-  220 2.0.0 Ready to start TLS
=== TLS started with cipher TLSv1.3:TLS_AES_256_GCM_SHA384:256
=== TLS no local certificate set
=== TLS peer DN="/C=AU/ST=Some-State/O=Internet Widgits Pty Ltd/CN=mail.domain.com"
 ~> EHLO ubuntu
<~  250-mail.domain.com
<~  250-PIPELINING
<~  250-SIZE 10240000
<~  250-ETRN
<~  250-AUTH PLAIN LOGIN
<~  250-AUTH=PLAIN LOGIN
<~  250-ENHANCEDSTATUSCODES
<~  250-8BITMIME
<~  250-DSN
<~  250 CHUNKING
 ~> AUTH LOGIN
<~  334 VXNlcm5hbWU6
 ~> YWRtaW5AZG9tYWluLmNvbQ==
<~  334 UGFzc3dvcmQ6
 ~> YWRtaW4xMjM=
<~  235 2.7.0 Authentication successful
 ~> MAIL FROM:<admin@domain.com>
<~  250 2.1.0 Ok
 ~> RCPT TO:<user@domain.com>
<~  250 2.1.5 Ok
 ~> DATA
<~  354 End data with <CR><LF>.<CR><LF>
 ~> Date: Mon, 08 Feb 2021 17:00:22 +0300
 ~> To: user@domain.com
 ~> From: admin@domain.com
 ~> Subject: test from admin
 ~> Message-Id: <20210208170022.076657@ubuntu>
 ~> X-Mailer: swaks v20190914.0 jetmore.org/john/code/swaks/
 ~> 
 ~> testing 123
 ~> 
 ~> 
 ~> .
<~  250 2.0.0 Ok: queued as 657C948415C
 ~> QUIT
<~  221 2.0.0 Bye
=== Connection closed with remote host.
```
## Инструкция по выполнению

## Инструкция по тестирование
```
как проверить
```
## Отчёт
как сдать?
* список требований
