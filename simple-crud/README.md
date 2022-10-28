# Simple CRUD

A simple CRUD application using [Echo](https://echo.labstack.com/) and [GORM](https://gorm.io/).


## How to run

### Initialize project

```bash
cp sample.env .env
```

### Set up DB

```bash
docker run --rm -p 3307:3306 --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:latest

mysql -u root -h 127.0.0.1 --port 3307 -p

mysql> CREATE DATABASE IF NOT EXISTS testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs;
mysql> USE testdb;
mysql> CREATE TABLE IF NOT EXISTS employees (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	start_date DATE NOT NULL
);
```

### Run the application

```bash
go run .
```

### Examples

```bash
curl -X POST -H 'Content-Type: application/json' \
	-d '{ "name": "Bob", "startDate": "2022-11-04" }' \
	http://localhost:1323/employees

# {"id":5,"name":"Bob","startDate":"2022-11-04"}
```

```bash
curl http://localhost:1323/employees/5

```
## References

- https://scrapbox.io/tsuchinaga/Golang%E3%81%AEJSON%E5%A4%89%E6%8F%9B%E3%81%A7%E4%BB%BB%E6%84%8F%E3%81%AE%E6%97%A5%E4%BB%98%E5%BD%A2%E5%BC%8F%E3%82%92%E5%A4%89%E6%8F%9B%E3%81%99%E3%82%8B
- https://stackoverflow.com/questions/23796543/go-checking-for-the-type-of-a-custom-error
- https://stackoverflow.com/questions/60954794/how-to-define-date-in-gorm
