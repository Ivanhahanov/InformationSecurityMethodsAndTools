###### tags: `Практика`
# Web Application Firewall
## Часть 1
Даная лабораторная работа разделена на две части, в первой части нам надо будет разобраться с уязвимостями таящимися в Веб приложениях, а во второй части воспользоваться программным решением для того, чтобы оперативно защититься от атак использующих эти уязвимости.
## Инструкция по установке
Для практики будем использовать **Damn Vulnerable Web Application (DVWA)**. Это веб-приложение PHP / MySQL, которое содержит много уязвимостей. Его основная цель - помочь специалистам по безопасности проверить свои навыки и инструменты, помочь веб-разработчикам лучше понять процессы защиты веб-приложений и помочь студентам изучать безопасность веб-приложений.
### Загрузка
```
git clone https://github.com/Ivanhahanov/InformtionsSecurityMethodsAndTools.git
```
### Развертывание
```
сd WAF
docker-compose up -d
```
## Инструкция по выполнению
Теперь по адресу http://localhost у нас развёрнуто веб-приложение.


Данные для входа **admin:password**

По адресу http://localhost/setup.php можно создать (при первом запуске) или перезагрузить базу данных (если что-то сломали)

Слева мы видим список доступных уязвимостей. Для примера возьмём **Command Injection**.

На сайте предоставлен список полезных ссылок, перейдя по которым можно получить более подробную информацию, а также имеется возможность посмотреть исходный код конкретного приложения

В левом нижнем углу можно увидеть уровень защищённости (сложности). 

DVWA предлагает нам 4 уровня сложности:
1. Low
2. Medium
3. Hard
4. Impossible

Изменить уровень сложности можно во вкладке **DVWA Security**.

## Что делать?
Задача состоит в том, чтобы попрактиковаться в эксплуатации веб уязвимостей. Для этого мы с вами сыграем роль **Пентестера**.

Пентестер или Специалист по анализу защищённости это человек в чьи обязанности входит тестирование защищённости веб приложений и инфраструктуры. Проще говоря он ищет уязвимости и использует их, но на законных основаниях, имея на руках исходный коды и права доступа.

Результатом работы такого специалиста является подробный отчёт о найденных уязвимостях.

Результатом вашей работы также будет отчёт в котором вы опишите какие уязвимости были найдены, и как эти уязвимости можно устранить.

Количество описанных уязвимостей зависит от вас, но минимальный порог следующий:
1. Command Injection
2. File Inclusion
3. SQL Injection
4. XSS reflected

Данный уязвимости необходимо проэксплуатировать минимум на трёх уровнях сложности.


## Отчёт
Структура отчёта следующая:

* Информация о типе уязвимости и его краткая характеристика;
* Информация о причине уязвимости;
* Указание сервисов, которые подвержены уязвимости;
* Proof of concept, который позволит подтвердить наличие уязвимости;
* Рекомендации по устранению уязвимости.

Отчёт предоставить в электронном виде. Помните что отчёт вы пишите для заказчика, который не разбирается в безопасности, но ему должно быть понятно какие последствия могут быть если он не устранит эти уязвимости, ну и описание по устранению этих уязвимостей должны быть понятны для разработчиков. 

# Web Application Firewall
## Часть 2
[**Файрвол веб-приложений**](https://ru.wikipedia.org/wiki/%D0%A4%D0%B0%D0%B9%D1%80%D0%B2%D0%BE%D0%BB_%D0%B2%D0%B5%D0%B1-%D0%BF%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B9) (англ. Web application firewall, WAF) — совокупность мониторов и фильтров, предназначенных для обнаружения и блокирования сетевых атак на веб-приложение. WAF относятся к прикладному уровню модели OSI.

Веб-приложение может быть защищено силами разработчиков самого приложения без использования WAF. Это требует дополнительных расходов при разработке. Например, содержание отдела информационной безопасности. WAF вобрали в себя возможность защиты от всех известных информационных атак, что позволяет делегировать ему функцию защиты. Это позволяет разработчикам сосредоточиться на реализации бизнес-логики приложения, не задумываясь о безопасности.

В реальном мире далеко не всегда есть время на устранения уязвимостей средствами разработчиков. Большие команды разработки не могут отреагировать на исправления так быстро как это может требоваться. На исправления критических уязвимостей могут уйти дни, а то и недели. За это время бизнес может потерять всё, и вендоры в области безопасности предлагают элегантное решение данной проблемы.

Решением выступает программное обеспечение способное просматривать входящий трафик и блокировать его в соответствии с набором правил. Важно заметить что данное решение помогает выиграть время, но не решает проблему, уязвимость всё равно надо устранять, но команде разработки проще делать это в спокойной обстановке.
### Загрузка
```
git clone https://github.com/Ivanhahanov/InformtionsSecurityMethodsAndTools.git
```
### Развертывание
```
сd WAF
docker-compose up -d
```

### Инструкция по выполнению
Исходный набор у нас такой же, как и в предыдущей практике, за исключением того что акцент мы будем делать на Web Application Firewall.
В директории `custom/` находится файл с примером различных правил `rules.conf`.

В этом файле необходимо написать правила, которые смогу закрыть уязвимости нашего DVWA.

Давайте представим, что разработчик не могу оперативно устранить уязвимости которые были найдены в процессе пентеста (тестирования на проникновения) и теперь мы выступаем в роли обороняющейся стороны и нам надо защитить наше приложение до того как хакеры до него доберутся, и дать возможность разработчика спокойно выполнять свою работу.

В качестве WAF мы будем использовать ModSecurity, он у нас уже развёрнут, но нам нужны правила для обнаружения атак.

Собственно написанием правил вам и предстоит заняться :)

[Документация](https://github.com/SpiderLabs/ModSecurity/wiki/Reference-Manual-(v2.x))

Для самопроверки можете воспользоваться скриптом exploitation.py
```
python3 exploitation.py localhost
[Command Injection]: Vulnerable 75.0%
[File Inclusion]: Vulnerable 100.0%
[SQL Injection]: Vulnerable 100.0%
[XSS Reflected]: Vulnerable 100.0%
```

В поле Command Injection указано 75% так как у нас уже есть одно правило:
```
SecRule ARGS "@contains &&" "phase:2,log,deny,msg:'Command injection'id:700009"
```
Обратите внимание, что `id` должен быть уникальным для каждого правила.

Для того, чтобы посмотреть логи WAF, выполните следующую команду:
```
docker-compose exec dvwa_with_modsecurity tail -f var/log/apache2/modsec_audit.log
```
Здесь мы можем видеть как отработало наше правило.

Чтобы отменить правило можно написать в конце файла строку `SecRuleRemoveById 700009` и тогда наш WAF не будет использовать это правило при анализе трафика.

# Отчёт 
Работу можно считать выполненной если скрипт exploitation.py покажет везде 0%, но преподаватель оставляет за собой право не принять лабораторную работу если правила помешают нормальной работе приложения