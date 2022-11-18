# Simple CRUD

A simple CRUD application using standard http and sql libraries.

## How to run

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
export DB_USER=root
export DB_PASS=my-secret-pw
export DB_HOST=127.0.0.1
export DB_PORT=3307
export DB_NAME=testdb
export SERVER_PORT=8080

go run .
```

### Examples

```bash
curl -X POST -H 'Content-Type: application/json' \
	-d '{ "name": "Bob", "startDate": "2022-11-04" }' \
	http://localhost:8080/employees/ 

# {"id":5,"name":"Bob","startDate":"2022-11-04"}
```

```bash
curl http://localhost:8080/employees/5
```