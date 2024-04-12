grpc contract https://github.com/RautaruukkiPalich/go_auth_grpc_contract

Простое приложение авторизации реализованное на gRPC

Реализован механизм сброса пароля, если его забыли: 
генерируется новый пароль,
отправляется в почтовый сервис через kafka,
новый пароль отправляется на указанный вами email

(smtp реализован в https://github.com/RautaruukkiPalich/go_auth_grpc_smtp)


Можно стартануть в 3 команды:
1) Запустить докер
```sh 
make dockerrun
```
2) Провести миграции
```sh 
make migrate
```
3) Запустить приложение
```sh 
make run
```
