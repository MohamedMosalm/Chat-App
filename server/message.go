package main

type Message struct {
    SenderID  string `json:"sender_id"`
    Content   string `json:"content"`
    Timestamp string `json:"timestamp"`
}
