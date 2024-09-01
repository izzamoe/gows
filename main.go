package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	timeCh      = make(chan string)
	activeConns = 0
	connMutex   = &sync.Mutex{}
	wg          = &sync.WaitGroup{}
)

func main() {
	app := fiber.New()

	// Middleware WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Endpoint WebSocket
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		connMutex.Lock()
		activeConns++
		if activeConns == 1 {
			wg.Add(1)
			go timeSender()
		}
		connMutex.Unlock()

		defer func() {
			connMutex.Lock()
			activeConns--
			if activeConns == 0 {
				// Alih-alih menutup channel, kirim sinyal berhenti ke timeSender
				timeCh <- "stop"
			}
			connMutex.Unlock()
			c.Close()
		}()

		for currentTime := range timeCh {
			// Hentikan loop jika timeSender mengirim sinyal "stop"
			if currentTime == "stop" {
				break
			}

			// Pastikan hanya satu operasi WriteMessage yang berjalan pada satu waktu
			if err := c.WriteMessage(websocket.TextMessage, []byte(currentTime)); err != nil {
				// user disconnected
				fmt.Println("User disconnected" + err.Error())
				break
			}
		}
	}))

	// Serve file HTML di root
	app.Static("/", "./public")

	log.Fatal(app.Listen(":3000"))
}

func timeSender() {
	defer wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case timeCh <- time.Now().Format("15:04:05"):
			// Waktu terkirim
		case <-timeCh:
			fmt.Println("timeCh ditutup")
			// timeCh ditutup, keluar dari loop
			return
		}
	}
}
