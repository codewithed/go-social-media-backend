package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=go-soc password=go-soc999 sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateTables()
}

func (s *PostgresStore) CreateTables() error {
	query := `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		firstName VARCHAR(255) NOT NULL,
		lastName VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		bio VARCHAR(255),
		passwordHash VARCHAR(1000) NOT NULL,
		created_at timestamptz NOT NULL
	);
	
	CREATE TABLE posts (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		content VARCHAR(255) NOT NULL,
		mediaUrl VARCHAR(10000),
		created_at timestamptz NOT NULL,
		FOREIGN KEY (userID) REFERENCES users (id)
	);

	CREATE TABLE comments (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		postID BIGINT NOT NULL,
		content VARCHAR(255) NOT NULL,
		created_at timestamptz NOT NULL,
		CONSTRAINT fk_userID FOREIGN KEY (userID) REFERENCES users (id),
		CONSTRAINT fk_postID FOREIGN KEY (postID) REFERENCES posts (id)
	);

	CREATE TABLE follows (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		followed_by_ID BIGINT NOT NULL,
		created_at timestamptz NOT NULL,
		CONSTRAINT fk_userID FOREIGN KEY (userID) REFERENCES users (id),
		CONSTRAINT fk_followed_by_ID FOREIGN KEY (followed_by_ID) REFERENCES users (id)
	);
	`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
