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
