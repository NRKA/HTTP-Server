# Задание

1) Написать CRUD операции для работы с бд
   должны быть реализованы методы
    - GetByID
    - Update
    - Create
    - Delete
2) Написать миграции для создания таблиц используя goose
3) Разработать HTTP сервер. Сервер должен уметь работать с бд и релизовывать CRUD операции.
    - Метод **GetByID**.
        - В Query-параметрах принимать идентификатор ?id= и возвращать данные из бд.
        - если передан идентификатор, который отсутствует в бд, возвращать HTTP-статус 404
        - если Query-параметрах отсутсвует ?id=, возвращать HTTP-статус 400
        - если произошла внутренняя ошибка сервера, возвращать HTTP-статус 500
        - если запрос успешно выполнен, возвращать HTTP-статус 200 и данные в теле ответа

    - Метод **Create**
        - В теле запроса принимать идентификатор и данные: и добавлять в бд
        - если такой идентификатор уже существует, возвращать HTTP-статус 409
        - если произошла внутренняя ошибка сервера, возвращать HTTP-статус 500
        - если запрос успешно выполнен, возвращать HTTP-статус 200

    - Метод **Delete**
        - В Query-параметрах принимать идентификатор ?id= и удалять из бд
        - если такого идентификатора не существует, возвращать HTTP-статус 404
        - если произошла внутренняя ошибка сервера, возвращать HTTP-статус 500
        - если запрос успешно выполнен, возвращать HTTP-статус 200

    - Метод **Update**
        - В теле запроса принимать идентификатор и данные и обновлять бд по ключу
        - если такого идентификатора не существует, возвращать HTTP-статус 404
        - если произошла внутренняя ошибка сервера, возвращать HTTP-статус 500
        - если запрос успешно выполнен, возвращать HTTP-статус 200

В случае ошибки в тело ответа пишем текст ошибки

4) Запустить сервер на порту 9000
5) Покрыть юнит тестами хендлеры. Минимальное покрытие - 40%.
6) Покрыть интеграционными тестами хендлеры дз 3 недели. Минимальные тест кейсы, успешное выполнение и получение ошибки из-за передачи некорректных данных.
7) Подготовить Makefile, в котором будут след команды: запуск тестового окружения при помощи docker-compose, запуск интеграционных тестов, запуск юнит тестов, запуск скрипта миграций, очищение базы от тествых данных
💎 В ридми приложить curl запросы, на каждую ручку. Запросы должны быть валидными и возвращать 200

Ограничения:
Нельзя использовать orm или sql билдеры
Для реализации http сервера можно использовать как net/http так и gin/fasthttp и прочее
Предметную область выбрать самостоятельно. Можно использовать пост/комментарий


Before sending requests you need to follow these steps:
1) Type the following command in your terminal - docker-compose up -d;
2) Type the following command in your terminal make test-migration-up;
3) Make sure that .env variable is in project folder.

<h3 style="text-align:center;">Here is my Curl Requests</h3>

1) **POST**

- **curl -X POST localhost:9000/article -d '{"name":"test","rating":15}' -i**

HTTP/1.1 200 OK

Date: Sun, 08 Oct 2023 15:40:20 GMT

Content-Length: 70

Content-Type: text/plain; charset=utf-8

{"ID":16,"Name":"test","Rating":15,"CreatedAt":"0001-01-01T00:00:00Z"}

3) **GET**

- **curl -X GET localhost:9000/article/1 -i**

HTTP/1.1 200 OK

Date: Sun, 08 Oct 2023 15:41:35 GMT

Content-Length: 85

Content-Type: text/plain; charset=utf-8

{"ID":1,"Name":"testname","Rating":44,"CreatedAt":"2023-10-05T18:12:38.011325+05:00"}



- **curl -X GET localhost:9000/article/10000 -i**

HTTP/1.1 404 Not Found

Date: Sun, 08 Oct 2023 15:42:34 GMT

Content-Length: 59

Content-Type: text/plain; charset=utf-8

Failed to find article:article not found

3) **DELETE**
- **curl -X DELETE localhost:9000/article/14 -i** ----->14 ID has been removed, please replace it with an existing ID

HTTP/1.1 200 OK

Date: Sun, 08 Oct 2023 15:43:22 GMT

Content-Length: 25

Content-Type: text/plain; charset=utf-8


- **curl -X DELETE localhost:9000/article/100000 -i**

HTTP/1.1 404 Not Found

Date: Sun, 08 Oct 2023 15:46:18 GMT

Content-Length: 45

Content-Type: text/plain; charset=utf-8

Failed to find article:article not found

4) **PUT**
- **curl -X PUT localhost:9000/article -d '{"id":1,"name":"testname","rating":44}' -i**

HTTP/1.1 200 OK

Date: Sun, 08 Oct 2023 15:46:54 GMT

Content-Length: 25

Content-Type: text/plain; charset=utf-8

- **curl -X PUT localhost:9000/article -d '{"id":100000,"name":"testname","rating":44}' -i**

HTTP/1.1 404 Not Found

Date: Sun, 08 Oct 2023 11:59:18 GMT

Content-Length: 45

Content-Type: text/plain; charset=utf-8

Failed to find article:article not found

