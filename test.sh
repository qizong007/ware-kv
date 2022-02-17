curl -X POST http://localhost:7777/str -H 'Content-Type: application/json' -d '{"k": "name","v": "qizong007"}'
sleep 1
curl -i http://localhost:7777/str/name
sleep 1
curl -X DELETE http://localhost:7777/str/name
sleep 1
curl -i http://localhost:7777/str/name