package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	// Middleware untuk mengecek apakah klien mendukung WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Endpoint WebSocket
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer c.Close()
		for {
			// Ambil waktu sekarang
			currentTime := time.Now().Format("15:04:05")
			// Kirim waktu ke klien
			if err := c.WriteMessage(websocket.TextMessage, []byte(currentTime)); err != nil {
				log.Println("Error mengirim waktu:", err)
				break
			}

			// Tunggu 1 detik
			time.Sleep(1 * time.Second)
		}
	}))

	// Serve file HTML di root
	app.Static("/", "./public")

	log.Fatal(app.Listen(":3000"))
}
