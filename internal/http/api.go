package api

import (
	"encoding/json"
	"fmt"
	"log"
	"messagio/internal/domain/models"
	"messagio/internal/kafka"
	"messagio/internal/storage/postgresql"

	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, db *postgresql.Storage, kafka *kafka.Producer) {
	router.HandleFunc("/messages", createMessageHandler(db, kafka)).Methods("POST")
	router.HandleFunc("/statistics", getStatisticsHandler(db)).Methods("GET")
}

func createMessageHandler(db *postgresql.Storage, kafka *kafka.Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg models.Message

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Printf("Failed to decode request payload: %v", err)
			return
		}

		id, err := db.SaveMessage(msg.Content)
		if err != nil {
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			log.Printf("Failed to save message: %v", err)
			return
		}

		if err := kafka.WriteMessage(nil, []byte(msg.Content)); err != nil {
			http.Error(w, "Failed to send message to Kafka", http.StatusInternalServerError)
			return
		}

		response := models.Message{
			ID:      id,
			Content: msg.Content,
		}

		messageData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to encode message", http.StatusInternalServerError)
			log.Printf("Failed to encode message: %v", err)
			return
		}

		if err := kafka.WriteMessage([]byte(fmt.Sprintf("%d", id)), messageData); err != nil {
			http.Error(w, "Failed to send message to Kafka", http.StatusInternalServerError)
			log.Printf("Failed to send message to Kafka: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Failed to encode response: %v", err)
		}
	}
}

func getStatisticsHandler(db *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := db.GetProcessedMessagesCount()
		if err != nil {
			http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
			log.Printf("Failed to get statistics: %v", err)
			return
		}

		response := map[string]int64{"processed_message_count": count}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
