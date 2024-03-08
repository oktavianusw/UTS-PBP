package models

type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Game struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	MaxPlayers int    `json:"max_players"`
}

type Room struct {
	ID       int    `json:"id"`
	RoomName string `json:"room_name"`
	GameID   int    `json:"game_id"`
}

type Participant struct {
	ID        int `json:"id"`
	RoomID    int `json:"room_id"`
	AccountID int `json:"account_id"`
}

type RoomResponseItem struct {
	ID       int    `json:"id"`
	RoomName string `json:"room_name"`
}

type RoomsResponse struct {
	Status int `json:"status"`
	Data   struct {
		Rooms []RoomResponseItem `json:"rooms"`
	} `json:"data"`
}

type GetDetailRoomParticipant struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Username  string `json:"username"`
}

type RoomDetailResponse struct {
	Status int `json:"status"`
	Data   struct {
		Room struct {
			ID           int                        `json:"id"`
			RoomName     string                     `json:"room_name"`
			Participants []GetDetailRoomParticipant `json:"participants"`
		} `json:"room"`
	} `json:"data"`
}

type InsertRoomResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type LeaveRoomResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
