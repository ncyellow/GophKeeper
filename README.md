# GophKeeper


## Работа через https
#### Для работы клиента в режиме https нам необходимо передать флаги или переменные окружения

client -addr "https://localhost"  -crypto-crt "*.crt"  -crypto-key "*.key" -crypto-ca "*.key"

где: \
-addr путь до сервера \
-crypto-ca ключ которым подписаны сертификаты\
-crypto-crt сертификат \
-crypto-key клиентский ключ

#### Для работы Сервера в режиме https нам необходимо передать флаги или переменные окружения

server -addr ":443"  -crypto-crt "*.crt"  -crypto-key "*.key" -dns "user=postgres password=12345 host=localhost port=5433 dbname=gophkeep"

где: \
-addr путь до сервера \
-crypto-crt сертификат \
-crypto-key серверный ключ \
-dns строка подключения к БД



## Работа через grpc
#### Для работы клиента в режиме grpc нам необходимо передать флаги или переменные окружения

client -grpc-addr ":3200"  

где: \
-grpc-addr путь до сервера

#### Для работы Сервера в режиме grpc нам необходимо передать флаги или переменные окружения

server -grpc-addr ":3200"  -dns "user=postgres password=12345 host=localhost port=5433 dbname=gophkeep"

где: \
-grpc-addr путь до сервера \
-dns строка подключения к БД
