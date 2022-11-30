# Long polling application

## Preparation

```bash
chmod +x client.sh
go run .
```

## Send messages

```bash
ROOM_ID=room_1
CLIENT_ID=$(curl -X POST -H 'Content-Type: application/json' -d "{ \"roomID\": \"${ROOM_ID}\" }" http://localhost:1323/enter 2>/dev/null)

curl -X POST \
       -d '{ "message": "Hello World!!" }' \
       -H 'Content-Type: application/json' \
       "http://localhost:1323/message?client_id=${CLIENT_ID}"
```

## Receive messages

```bash
./client.sh
```

## Leave chat rooms

```bash
curl -X DELETE "http://localhost:1323/leave?client_id=${CLIENT_ID}"
```

## TODO

- Make it possible for clients to leave thier rooms