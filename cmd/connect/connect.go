package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/d5avard/factory/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionRec struct {
	ID       int
	Question string
}

type AnswerRec struct {
	ID     int
	Answer string
}

type TagsRec struct {
	ID   int
	Tags []string
}

func main() {
	filename, err := internal.GetConfigFilename()
	if err != nil {
		log.Fatalf("Error getting config filename: %v", err)
	}

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

	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("connect to db error: %w", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	q := QuestionRec{Question: "What is the capital of France?"}
	a := AnswerRec{Answer: "Paris"}
	ts := TagsRec{Tags: []string{"geography", "capital", "France"}}
	var id int
	id, err = AddQuestionAnswerTags(ctx, conn, q, a, ts)
	if err != nil {
		return fmt.Errorf("add question answer tags error: %w", err)
	}
	log.Println("Added question, answer, and tags to database", id)

	return nil
}

func AddQuestionAnswerTags(ctx context.Context, conn *pgxpool.Pool, q QuestionRec, a AnswerRec, ts TagsRec) (int, error) {
	const INSERT_QUERY = `SELECT insert_question_answer_tags($1, $2, $3)`
	var questionId int
	err := conn.QueryRow(ctx, INSERT_QUERY, q.Question, a.Answer, ts.Tags).Scan(&questionId)
	if err != nil {
		return -1, err
	}

	return questionId, nil
}
