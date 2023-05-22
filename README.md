## Subscription Service

### Administrative Side

The system provides administrative endpoints where operations like <br />
* Setting up Products
* Setting up vouchers
* Setting up subscription plans 

### Seeding Data
* RUN `cp .env.example .env` and fill out the database connection environment variables
* RUN `make seed`


### Start Up
#### Local
 * RUN `cp .env.example .env` and setup the database values
 * RUN `make start`

#### Docker
* RUN `make docker-start`


## Usage
> Some endpoints need authorization. <br />
> Assumption is that the server runs on `localhost:8080` <br />
> `{BEARER_TOKEN}`/`{ACCESS_TOKEN}` are placeholders.

### Create User
Request:
```bash
curl --location 'localhost:8080/user/create' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {BEARER_TOKEN}' \
--data-raw '{
    "name": "Angelina Emmerich",
    "email": "terrillrolfson@anderson.org",
    "password": "password",
    "user_type": "admin"
}'
```
Response: 201
```json
{
    "id": "646a8ed27e31426e275edeb6",
    "name": "Angelina Emmerich",
    "email": "terrillrolfson@anderson.org",
    "token": "{ACCESS_TOKEN}",
    "created_at": "2023-05-21T23:36:18.820031+02:00"
}
```

### Authenticate User
Request:
```bash
curl --location 'localhost:8080/login' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {BEARER_TOKEN}' \
--data-raw '{
    "email": "terrillrolfson@anderson.org",
    "password": "password"
}'
```

Response:
```json
{
    "id": "646a8ed27e31426e275edeb6",
    "name": "Angelina Emmerich",
    "email": "terrillrolfson@anderson.org",
    "token": "{ACCESS_TOKEN}",
    "created_at": "2023-05-21T23:36:18.82+02:00"
}
```

### Create new Product
Request:
```bash
curl --location 'localhost:8080/products' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {BEARER_TOKEN}' \
--data '{
    "name": "encryption",
    "description": "facere",
    "tax": 16.27,
    "trial_exists": true
}'
```
Response:
```json
{
    "id": "646bb222f7752a37417c6482",
    "name": "encryption",
    "description": "facere",
    "tax": 16.27,
    "created_at": "2023-05-22T20:19:14.378662+02:00",
    "updated_at": "2023-05-22T20:19:14.38+02:00"
}
```

### Fetch all Products
```bash
curl --location 'localhost:8080/products'
```

### Fetch Single Product
```bash
curl --location 'localhost:8080/products/:product_id'
```
>`:product_id` is the product ID

### Create Subscription Plan
```bash
curl --location 'localhost:8080/plans' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {BEARER_TOKEN}' \
--data '{
    "product_id": "646a854f586097c74343f315",
    "amount": 200,
    "duration": 3,
    "trial_duration": 0
}'
```

### Create Trial Subscription Plan
> Setting the `trial_duration` to the number of desired months
```bash
curl --location 'localhost:8080/plans' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {BEARER_TOKEN}' \
--data '{
    "product_id": "646a854f586097c74343f315",
    "amount": 200,
    "duration": 3,
    "trial_duration": 1
}'
```

### Buy Single Product
```bash
curl --location 'localhost:8080/subscriptions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {BEARER_TOKEN}' \
--data '{
    "product_id": "646a8c4f0220f41221076550",
    "plan_id": "646a8c4f0220f41221076552"
}'
```

### Buy Single Product with a voucher
> A `voucher` code can be provided. If the voucher is invalid, an error is returned
```bash
curl --location 'localhost:8080/subscriptions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {BEARER_TOKEN}' \
--data '{
    "product_id": "646a8c4f0220f41221076550",
    "plan_id": "646a8c4f0220f41221076552",
    "voucher": "WELCOME_STAFF_ONLY"
}'
```

### Fetch all subscriptions
>Returns all the user subscriptions
```bash
curl --location 'localhost:8080/subscriptions' \
--header 'Authorization: Bearer {BEARER_TOKEN}'
```

Response:
```json
[
    {
        "id": "646a8d577e31426e275edeb5",
        "start_date": "2023-05-21T23:29:59.452+02:00",
        "end_date": "2023-09-18T23:29:59.452+02:00",
        "duration": 4,
        "trial_duration": 0,
        "total_duration": 4,
        "price": 1881.49,
        "tax": 444.28,
        "discount": 282.22,
        "total": 2043.54,
        "status": "active"
    }
]
```

### Pause Subscription
> You cannot pause during trial period. If you try to, it returns an error
```bash
curl --location --request PATCH 'localhost:8080/subscriptions/:subscription_id/pause' \
--header 'Authorization: Bearer {BEARER_TOKEN}'
```
>`:subscription_id` is the unique subscription ID

Response: 
```json
"subscription paused successfully"
```

### Unpause Subscription
```bash
curl --location --request PATCH 'localhost:8080/subscriptions/:subscription_id/unpause' \
--header 'Authorization: Bearer {BEARER_TOKEN}'
```
>`:subscription_id` is the unique subscription ID

Response: 
```json
"subscription unpaused successfully"
```

### Cancel Subscription
```bash
curl --location --request DELETE 'localhost:8080/subscriptions/:subscription_id' \
--header 'Authorization: Bearer {BEARER_TOKEN}'
```
>`:subscription_id` is the unique subscription ID

Response: 
```json
"subscription cancelled successfully"
```