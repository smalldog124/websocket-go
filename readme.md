# golang websocket
code golag ทำ server websocket และ client connect to server
## server1
เปิด server port :3000

ใช้ gofiber websocket v2
```go
 go run server1/server1.go
```
## client
ต่อ websocket ไปต่อที่ `ws://localhost:3000/join-group`

ใช้ gorilla websocket
```go
go run client1/client1.go
```