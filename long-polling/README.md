# Long polling application

## Preparation

```bash
chmod +x client.sh
go run .
```

## Send messages

```bash
CLIENT_ID=`curl -X POST -H 'Content-Type: application/json' -d '{ "roomID": "hello" }' http://localhost:1323/join 2>/dev/null`
curl -X POST \
       -d '{ "message": "Hello World!!" }' \
       -H 'Content-Type: application/json' \
       "http://localhost:1323/message?client_id=${CLIENT_ID}"
```

## Recieve messages

```bash
./client.sh
```

## TODO

- Make it possible for clients to leave thier rooms