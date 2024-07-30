package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"messagio/internal/config"
	"messagio/internal/http"
	"messagio/internal/storage/postgresql"

	"net/http"
)

type Server struct {
	Router *mux.Router
	DB     *postgresql.Storage
	//Kafka  *kafka.Producer
	Config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	db, err := postgresql.New(cfg.PostgreSQL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	//kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	//if err != nil {
	//	log.Fatalf("Could not create Kafka producer: %v", err)
	//}

	router := mux.NewRouter()
	api.RegisterRoutes(router, db)

	return &Server{Router: router, DB: db, Config: cfg}
}

func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Config.HTTP.Port), s.Router)
}
