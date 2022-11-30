#!/bin/bash

# references
## https://askubuntu.com/questions/714458/bash-script-store-curl-output-in-variable-then-format-against-string-in-va
## https://www.cyberciti.biz/faq/bash-infinite-loop/
## https://opensource.com/article/17/1/getting-started-shell-scripting

echo "Press [CTRL+C] to stop.."

ROOM_ID=$1
echo "Enter the room: ${ROOM_ID}"

CLIENT_ID=$(curl -X POST -H 'Content-Type: application/json' -d "{ \"roomID\": \"${ROOM_ID}\" }" http://localhost:1323/enter 2>/dev/null)
echo "client_id=${CLIENT_ID}"

while :
do
	MESSAGE=`curl http://localhost:1323/message?client_id=${CLIENT_ID} 2>/dev/null`
	
	if [ "${MESSAGE}" != "" ]; then
		echo ${MESSAGE}
	fi
done