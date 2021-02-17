###### tags: `Практика`
# Фишинг
[**Фи́шинг**](https://ru.wikipedia.org/wiki/%D0%A4%D0%B8%D1%88%D0%B8%D0%BD%D0%B3) (англ. phishing от fishing «рыбная ловля, выуживание») — вид интернет-мошенничества, целью которого является получение доступа к конфиденциальным данным пользователей — логинам и паролям.

:::danger
Данная лабораторная предоставлена **исключительно в обучающих целях**! Не повторяйте эти действия без письменного согласия заказчика! Составитель лабораторной работы, а также преподаватели не несут никакой ответственности за деяния студентов!
:::

В этой лабораторной работе мы будем использовать инструмент [Gophish](https://getgophish.com/) для проведения фишинговой рассылки. Для этого мы поднимем свой почтовый сервер и будем экспериментировать на нём. Напоминаю, что проведения подобной атаки без письменного разрешения является преступлением, по-этому мы сами создадим себе пользователей и наши фишиговые письма останутся в локальной сети. 

Чем опасен фишинг? Люди по-прежнему остаются самым уязвимым местом информационной системы, в связи с этим один из самых простых способов атаки и дальнейшего распространения является атака на сотрудников, далеко не все сотрудники организаций способны противостоять своим эмоциям и не податься на провокацию злоумышленника. 

Для повышения уровня информационной грамотности следует проводить тренировочные рассылки вредоносных писем, чтобы сотрудники научились распознавать подобные атаки и не повелись на удочку реального злоумышленника.

Подобные проверки стоит проводить раз в две недели или один месяц, чтобы держать коллектив в тонусе

Цена данной услуги варьируется от количества хостов и проработанности аттаки, но удовольствие не из дешёвых :)

### Используемые инструменты
* [SMTP-сервер](https://ru.wikipedia.org/wiki/SMTP) (Почтовый сервер)
* [Gophish](https://getgophish.com/)
## Инструкция по установке
### Загрузка
```
git clone https://github.com/Ivanhahanov/InformationSecurityMethodsAndTools.git
```
После успешного клонирования репозитория на компьютере создаться директория InformationSecurityMethodsAndTools
```shell
cd InformationSecurityMethodsAndTools
```
И перейти в директорию `Phishing`
```shell
cd Phishing
```

### Развертывание
На данном этапе необязательно выполнять эти команды, так как предстоит создать TLS сертификаты и пользователей почтового сервера, но вы можете выполнить их, чтобы проверить что у вас работает `Docker` и `docker-compose`
```
docker-compose up -d
```
Флаг `-d` запустит в режиме демона и вы не будете видеть логи.
Сейчас с ресурса [DockerHub](https://hub.docker.com/) "спулятся" образы почтового сервера и приложения для фишинга.

[DockerHub](https://hub.docker.com/) - это как GitHub только для docker образов, вы можете туда загружать свои образы, а можете загружать чужие. Все образы имеют рейтинг, вы можете узнать сколько раз образ загружали и сколько ему поставили звёздочек. Это полезная информация, так вы можете понять на сколько тот или иной образ популярен и решить стоит ли загружать именно его или поискать другую реализацию.

Чтобы посмотреть что происходит внутри:
```
docker-compose logs -f
```
Mail server должен сообщить об ошибках, так как у нас нет пользователей и нет подписанного сертификата.
Так что можно выполнить эту команду:
```shell
docker-compose down
```
И Mail сервер плавно завершит работу
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
После добавления сертификата надо пересобрать наш mail server, чтобы можно наш сертификат добавился
```shell
docker-compose up --build
```
Флаг `--build` соберёт контейнеры заново, если появились какие-то изменения в файловой структуре или конфигурациях 
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

### Настройка почтового клиета
в качетсве почтового клиента возьмём предустановленный [ThunderBird](https://www.thunderbird.net/ru/)
![](https://i.imgur.com/a4cBEzb.png)
![](https://i.imgur.com/yzweWsU.png)

## Инструкция по выполнению
0. Открыть в браузере веб приложение Gophish по адресу https://localhost:3333
1. Авторизироваться в системе. Пароль находится в логах, а логин - **admin**
![](https://i.imgur.com/0C2udbl.png)
2. Создать новый пароль
3. Теперь приступим к настройке ...
4. Настроить профиль отправителя(**Sending Profile**)
![](https://i.imgur.com/NGfWB1V.png)

:::warning
IP хоста взять из докера с помощью команды **`docker exec mail ip a`**
:::
![](https://i.imgur.com/CA2vlUO.png)
:::info
проверить работоспособность можно с помощью кнопки **Send Test Email** 
:::
![](https://i.imgur.com/7sAcxKR.png)
5. Настройка фальшивой страницы(**Landing Page**)

Нажать на кнопку "Import Site" и ввести туда адрес на ваше усмотрение, я выбрал HackTheBox
![](https://i.imgur.com/tJTVMCn.png)

6. Настройка шаблона письма(**Email Template**)
![](https://i.imgur.com/SMEQ5g1.png)

7. Настройка получателей(User & Groups)
![](https://i.imgur.com/FZkOUst.png)
8. Настройка рассылки компании(**Compaigns**)
![](https://i.imgur.com/I0EsBHM.png)
9. Результат вы можете увидеть у себя во вкладке **Dashboard**. При переходе по ссылке в письме, это будет фиксироваться в круговой диаграмме и на графике, если за основу взять шаблон с логином и паролем, то система автоматически быдет фиксировать ввод пароля
![](https://i.imgur.com/kfC6cIM.png)


## Отчёт
Варианты сдачи, на усмотрение преподавателя
1. Сформировать отчёт в формате .pdf приложить туда скрины поддельных страниц и писем которые пришли созданным вами пользователями. В отчёте предоставить информацию о проведенном тестировании, сколько человек перешло по ссылкам, какой процент сотрудников  вашей вымышленной компании подверглись атаке и описать советы по предотвращении подобных атак. Проявите фантазию, чтобы преподавателю тоже было интересно читать отчёты.
2. Предоставить всё вышеперечисленное преподавателю и ответить на каверзные вопросы :-)

