package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type ChatRoom struct {
	clients map[*websocket.Conn]bool
}
type JoinGroup struct {
	Group    string `json:"group"`
	DriverID string `json:"driver_id"`
}

func main() {
	app := fiber.New()
	// Optional middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if c.Get("host") == "localhost:3000" {
			c.Locals("Host", "Localhost:3001")
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})

	// Upgraded websocket request
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		fmt.Println(c.Locals("Host")) // "Localhost:3000"
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s type:%d", msg, mt)
			// err = c.WriteMessage(mt, msg)
			// if err != nil {
			// 	log.Println("write:", err)
			// 	break
			// }
		}
	}))
	chatRoom := &ChatRoom{
		clients: make(map[*websocket.Conn]bool),
	}
	room := map[string]*ChatRoom{
		"bkk": chatRoom,
	}
	app.Get("/join-group", websocket.New(func(c *websocket.Conn) {
		var req JoinGroup
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read: --->", err)
				break
			}
			log.Printf("recv: %s type:%d", msg, mt)
			err = json.Unmarshal(msg, &req)
			room[req.Group].clients[c] = true
		}
	}))
	app.Get("/broadcast", func(c *fiber.Ctx) error {
		bkkGroup := room["bkk"]
		message := "มีงานที่ตลาดพลูใหญ่"
		for client := range bkkGroup.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Error broadcasting message:", err)
				client.Close()
			}
		}
		return nil
	})
	// ws://localhost:3000/ws
	log.Fatal(app.Listen(":3000"))
}
