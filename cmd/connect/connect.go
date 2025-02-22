package main

/*
dbURL might look like:
"postgres://username:password@localhost:5432/database_name"
*/

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/d5avard/factory/internal"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	var filename string
	var err error

	filename, err = internal.GetConfigFilename()
	if err != nil {
		log.Fatalf("Error getting config filename: %v", err)
	}

	// Load config file
	config, err := internal.LoadConfig(filename)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := connect(config); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
}

func connect(config internal.Config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("connect to db error: %w", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}
	log.Println("Connected to database")

	s, err := NewStorage(ctx, conn)
	if err != nil {
		return fmt.Errorf("new storage error: %w", err)
	}
	defer s.Close()

	q, err := s.GetQuestion(ctx, 1)
	if err != nil {
		return fmt.Errorf("get question error: %w", err)
	}
	log.Printf("Retrieved question: %+v\n", q)

	q.Question = "What is the capital of France?"
	err = s.AddQuestion(ctx, q)
	if err != nil {
		return fmt.Errorf("add question error: %w", err)
	}
	log.Println("Added question to database")

	return nil
}

type QuestionRec struct {
	ID       int
	Question string
}

type Storage struct {
	conn            *sql.DB
	getQuestionStmt *sql.Stmt
	addQuestionStmt *sql.Stmt
}

func NewStorage(ctx context.Context, conn *sql.DB) (*Storage, error) {
	s := Storage{conn: conn}
	var err error

	s.getQuestionStmt, err = conn.PrepareContext(
		ctx,
		`SELECT "question" FROM "questions" WHERE "id" = $1`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	s.addQuestionStmt, err = conn.PrepareContext(
		ctx,
		`INSERT INTO "questions" ("question") VALUES ($1)`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	return &s, nil
}

func (s *Storage) Close() {
	if s.getQuestionStmt != nil {
		s.getQuestionStmt.Close()
	}
	if s.addQuestionStmt != nil {
		s.addQuestionStmt.Close()
	}
}

func (s *Storage) GetQuestion(ctx context.Context, id int) (QuestionRec, error) {
	q := QuestionRec{ID: id}
	err := s.getQuestionStmt.QueryRowContext(ctx, id).Scan(&q.Question)
	if err != nil {
		return q, fmt.Errorf("failed to execute get question query: %w", err)
	}
	return q, nil
}

func (s *Storage) AddQuestion(ctx context.Context, q QuestionRec) error {
	_, err := s.addQuestionStmt.ExecContext(ctx, q.Question)
	if err != nil {
		return fmt.Errorf("failed to execute add question query: %w", err)
	}
	return nil
}
