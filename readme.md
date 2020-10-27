# user_balance_service

*Сервис работает с балансом пользователей и имеет ручки:*

-   Получение информации о балансе</br>
Request:
[GET] /account/balance/info/{id:[0-9]+}?currency=CUR

Responce:
<pre>
200
{
	"Id": 1,
	"Balance": 500.00,
	"Currency": "RUB"
}
500
{
    "Message": "Не удалось предоставить информацию для выбранного курса валюты",
    "ErrCode": 500
}
</pre>

-   Изменение баланса</br>
Request:
[POST] /account/balance/change
<pre>
Body:
{
    "Id":1,
    "Delta":155
}
</pre>

Responce:
<pre>
200
{
    "Message": "Аккаунт 1 успешно пополнен на сумму 155.00 руб."
}
403
{
    "Message": "Недостаточно средств на счету",
    "ErrCode": 403
}
</pre>

-   Перевод между счетами</br>
Request:
[POST] /account/balance/transfer
<pre>
Body:
{
    "Id1":1,
    "Id2":2,
    "Delta":120
}
</pre>

Responce:
<pre>
200
{
    "Message": "Перевод на сумму 120.00 руб. с аккаунта 1 на аккаунт 2 выполнен успешно."
}
403
{
    "Message": "Недостаточно средств на счету",
    "ErrCode": 403
}
</pre>

-   История транзакций</br>
Request:
[POST] /account/balance/history

<pre>
Body:
{   "Id":1,
    "SortedBy":"transaction_time",  //необязательное поле, параметры: "transaction_time", "transaction_sum"
    "SortedByDesc":true             //необязательное поле
}
</pre>

Responce:
<pre>
200
[
    {
        "AccountId": 1,
        "Delta": 1000,
        "RemainingBalance": 1000,
        "TransactionMessage": "Аккаунт 1 успешно пополнен на сумму 1000.00 руб.",
        "CreatedAt": "2020-09-21T18:45:15.278878Z"
    },...
]
400
{
    "Message": "Некорректные входные данные",
    "ErrCode": 400
}
404
{
    "Message": "Отсутсвуют записи по выбранным условиям поиска",
    "ErrCode": 404
}
</pre>

*Сервис развертывается, используя базу данных Postgres. Для развертывания сервиса с использованием docker-compose необходимо создать образ базы данных с настроенными таблицами*

Порядок развертывания сервиса через docker-compose:
-   docker-compose up

*Предполагается, что ручки используются из-за firewall, и недоступны простому пользователю.*
