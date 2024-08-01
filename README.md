> How to run development mode
1. Initiate MySQL container as dev database for develpoment.
```sh
$ docker run \
--name dev-mysql \
-p 3306:3306 \
-e MYSQL_DATABASE=go-test \
-e MYSQL_ROOT_PASSWORD=dev \
-d mysql:8.0
```
```sh
$ cd src
$ gin run main.go
```

> The backend platform service

[GCP Asset](https://console.cloud.google.com/home/dashboard?hl=en&project=asset-430906)

> Debug commands
1. If there is a requirement to access into container DB to investage data. Run the command to enter mysql interactive mode. Enter MySQL root password as terminal prompt a question.
```sh
$docker exec -it dev-mysql mysql -uroot -p
```
2. Command to initiate MySQL database schema by a pre-define script.
```sh
$ docker exec -i dev-mysql sh -c 'exec mysql -uroot -p"dev"' < /Users/paul/Study/asset-go/mysql-init.sql
```