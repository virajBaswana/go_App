package message

import "time"

type Message struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

type SocketMessage struct {
	Sender_id   int    `json:"sender_id"`
	Receiver_id int    `json:"receiver_id"`
	Content     string `json:"content"`
}

type MessageExchange struct {
	Id           int       `json:"id"`
	Message_Id   int       `json:"message_id"`
	Sender_Id    int       `json:"sender_id"`
	Receiver_Id  int       `json:"receiver_id"`
	Sent_At      time.Time `json:"sent_at"`
	Delivered_At time.Time `json:"delivered_at"`
	Read_At      time.Time `json:"read_at"`
}
