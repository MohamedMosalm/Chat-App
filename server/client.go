package main

import "github.com/gofiber/websocket/v2"

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
}
