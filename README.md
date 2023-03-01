# pay_system v1

The first version.

## How to run

1. Environment variables must be set. Create `.env` file:

```bush
$ echo export 'LOG_LEVEL=debug
export MAIL_HOST=smtp.gmail.com
export MAIL_PORT=587
export MAIL_USERNAME=pay.system.bank@gmail.com
export MAIL_APP_PASSWORD=<>
export MYSQL_USER=root
export MYSQL_PASSWORD=<>
export MYSQL_HOST=localhost
export MYSQL_DATABASE=pay_system
export HTTP_PORT=8080
export HMAC_SECRET=pays
export ADMIN_USER_FULLNAME=Serhi Shevchenko
export ADMIN_USER_EMAIL=serhi.shevchenko@gmail.com' >> .env
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

```Methods for user:
(користувач може зареєструватися/залогінитись/відправити лист для скидання паролю/змінити пароль/переглянути свої логи)
- POST {{host}}/api/v1/users/register
- POST {{host}}/api/v1/users/login
- POST {{host}}/api/v1/users/sendemail
- PATCH {{host}}/api/v1/users/resetpassword
- GET {{host}}/api/v1/users/search_logs

(користувач може створити рахунок/переглянути лише свої рахунки/заблокувати чи розблокувати свій рахунок/поповнити свій рахунок)
- GET {{host}}/api/v1/bank_account/search
- POST {{host}}/api/v1/bank_account/create
- PATCH {{host}}/api/v1/bank_account/lock
- PATCH {{host}}/api/v1/bank_account/unlock
- PATCH {{host}}/api/v1/bank_account/top_up

(користувач може створити платіж/переглянути лише свої платежі/надіслати платіж)
- GET {{host}}/api/v1/payment/search
- POST {{host}}/api/v1/payment/create
- PATCH {{host}}/api/v1/payment/sent

Methods for admin
(адміністратор може переглянути усіх корстувачів та рахунки/заблокувати чи розблокувати користувача чи рахунок/переглянути логи користувачів)
- GET {{host}}/api/v1/users/search
- GET {{host}}/api/v1/users/search_logs
- GET {{host}}/api/v1/bank_account/search
- PATCH {{host}}/api/v1/bank_account/lock
- PATCH {{host}}/api/v1/bank_account/unlock
- PATCH {{host}}/api/v1/admin/lock_user
- PATCH {{host}}/api/v1/admin/unlock_user
```