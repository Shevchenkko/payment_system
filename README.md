# pay_system v1

The first version.

## How to run

1. Environment variables must be set. Create `.env` file:

```bush
$ echo export 'LOG_LEVEL=<>
export MAIL_HOST=<>
export MAIL_PORT=<>
export MAIL_USERNAME=<>
export MAIL_APP_PASSWORD=<>
export MYSQL_USER=<>
export MYSQL_PASSWORD=<>
export MYSQL_HOST=<>
export MYSQL_DATABASE=<>
export HTTP_PORT=<>
export HMAC_SECRET=<>
export ADMIN_USER_FULLNAME=<>
export ADMIN_USER_EMAIL=<>' >> .env
```

and run 

```bush
$ source .env
```

2. Run program:

`go run main.go`

### Roles and methods

Roles:
-- admin
-- user

```
Methods for user:
(користувач може зареєструватися/залогінитись/відправити лист для скидання паролю/змінити пароль/переглянути свої логи)
- POST  {{host}}/api/v1/users/register
- POST  {{host}}/api/v1/users/login
- POST  {{host}}/api/v1/users/sendemail
- PATCH {{host}}/api/v1/users/resetpassword
- GET   {{host}}/api/v1/users/search_logs

(користувач може створити рахунок/переглянути лише свої рахунки/заблокувати чи розблокувати свій рахунок/поповнити свій рахунок)
- GET   {{host}}/api/v1/bank_account/search
- POST  {{host}}/api/v1/bank_account/create
- PATCH {{host}}/api/v1/bank_account/lock
- PATCH {{host}}/api/v1/bank_account/unlock
- PATCH {{host}}/api/v1/bank_account/top_up

(користувач може створити платіж/переглянути лише свої платежі/надіслати платіж)
- GET   {{host}}/api/v1/payment/search
- POST  {{host}}/api/v1/payment/create
- PATCH {{host}}/api/v1/payment/sent

Methods for admin
(адміністратор може переглянути усіх корстувачів та рахунки/заблокувати чи розблокувати користувача чи рахунок/переглянути логи користувачів)
- GET   {{host}}/api/v1/users/search
- GET   {{host}}/api/v1/users/search_logs
- GET   {{host}}/api/v1/bank_account/search
- PATCH {{host}}/api/v1/bank_account/lock
- PATCH {{host}}/api/v1/bank_account/unlock
- PATCH {{host}}/api/v1/admin/lock_user
- PATCH {{host}}/api/v1/admin/unlock_user
```