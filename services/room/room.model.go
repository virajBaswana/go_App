package room

import "time"

type Room struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  string    `json:"updated_at"`
}

type RoomMembers struct {
	Room_Id   int       `json:"room_id"`
	User_Id   int       `json:"user_id"`
	Joined_At time.Time `json:"joined_at"`
	Left_At   time.Time `json:"left_at"`
}
