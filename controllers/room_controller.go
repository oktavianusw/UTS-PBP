package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"week2/models"

	"github.com/gorilla/mux"
)

func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, room_name FROM rooms")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rooms []models.RoomResponseItem
	for rows.Next() {
		var room models.RoomResponseItem
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			log.Fatal(err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	response := models.RoomsResponse{
		Status: 200,
		Data: struct {
			Rooms []models.RoomResponseItem `json:"rooms"`
		}{
			Rooms: rooms,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDetailRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	row := db.QueryRow("SELECT * FROM rooms WHERE id = ?", id)

	var room models.Room
	err := row.Scan(&room.ID, &room.RoomName, &room.GameID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		} else {
			log.Fatal(err)
		}
	}
	rows, err := db.Query(`
	SELECT p.id, p.id_account, a.username 
	FROM participants p 
	INNER JOIN accounts a ON p.id_account = a.id 
	WHERE p.id_room = ?
    `, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var participants []models.GetDetailRoomParticipant
	for rows.Next() {
		var participant models.GetDetailRoomParticipant
		err := rows.Scan(&participant.ID, &participant.AccountID, &participant.Username)
		if err != nil {
			log.Fatal(err)
		}
		participants = append(participants, participant)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	response := models.RoomDetailResponse{
		Status: 200,
		Data: struct {
			Room struct {
				ID           int                               `json:"id"`
				RoomName     string                            `json:"room_name"`
				Participants []models.GetDetailRoomParticipant `json:"participants"`
			} `json:"room"`
		}{
			Room: struct {
				ID           int                               `json:"id"`
				RoomName     string                            `json:"room_name"`
				Participants []models.GetDetailRoomParticipant `json:"participants"`
			}{
				ID:           room.ID,
				RoomName:     room.RoomName,
				Participants: participants,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	accountID := vars["account_id"]

	row := db.QueryRow("SELECT max_players FROM games g INNER JOIN rooms r ON g.id = r.id_game WHERE r.id = ?", id)

	var maxPlayers int
	err := row.Scan(&maxPlayers)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		} else {
			log.Fatal(err)
		}
	}

	row = db.QueryRow("SELECT id FROM accounts WHERE id = ?", accountID)

	var accountExists int
	err = row.Scan(&accountExists)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		} else {
			log.Fatal(err)
		}
	}

	row = db.QueryRow("SELECT COUNT(*) FROM participants WHERE id_room = ?", id)

	var participantCount int
	err = row.Scan(&participantCount)
	if err != nil {
		log.Fatal(err)
	}

	if participantCount < maxPlayers {
		_, err = db.Exec("INSERT INTO participants (id_room, id_account) VALUES (?, ?)", id, accountID)
		if err != nil {
			log.Fatal(err)
		}

		response := models.InsertRoomResponse{
			Status:  200,
			Message: "Successfully joined the room",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := models.InsertRoomResponse{
			Status:  400,
			Message: "The room is full",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type LeaveRoomResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	accountID := vars["account_id"]

	row := db.QueryRow("SELECT id FROM rooms WHERE id = ?", id)

	var roomExists int
	err := row.Scan(&roomExists)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		} else {
			log.Fatal(err)
		}
	}

	row = db.QueryRow("SELECT id FROM accounts WHERE id = ?", accountID)

	var accountExists int
	err = row.Scan(&accountExists)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		} else {
			log.Fatal(err)
		}
	}

	result, err := db.Exec("DELETE FROM participants WHERE id_room = ? AND id_account = ?", id, accountID)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		http.Error(w, "Account not in room", http.StatusNotFound)
		return
	}

	response := LeaveRoomResponse{
		Status:  200,
		Message: "Successfully left the room",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
