package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"messagio/internal/config"
	"messagio/internal/domain/models"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.PostgresConfig) (*Storage, error) {
	const op = "storage.postgresql.New"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveMessage(content string) (int64, error) {
	const op = "storage.postgresql.SaveMessage"

	query := "INSERT INTO messages (content, created_at, processed, processed_at) VALUES ($1, $2, $3, $4) RETURNING id"

	var id int64
	err := s.db.QueryRow(query, content, time.Now(), false, nil).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetMessages() ([]models.Message, error) {
	const op = "storage.postgresql.GetMessages"

	rows, err := s.db.Query("SELECT id, content, created_at FROM messages")
	if err != nil {
		return []models.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("%s: %w", op, cerr)
		}
	}()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, err
}

func (s *Storage) GetProcessedMessagesCount() (int64, error) {
	const op = "storage.postgresql.GetProcessedMessagesCount"

	var count int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM messages WHERE processed = TRUE").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}

func (s *Storage) Close() error {
	return s.db.Close()
}
