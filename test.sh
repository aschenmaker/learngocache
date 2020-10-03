#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8801 &
./server -port=8802 &
./server -port=8803 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &



wait