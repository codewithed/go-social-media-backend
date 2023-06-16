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

// CRUD OPERATIONS FOR USERS
func (s *PostgresStore) GetAllUsers() ([]*User, error) {
	rows, err := s.db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for rows.Next() {
		user, err := ScanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *PostgresStore) GetUser(id int) (*User, error) {
	rows, err := s.db.Query(`SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoUser(rows)
	}

	return nil, nil
}

func (s *PostgresStore) CreateUser(user *User) error {
	_, err := s.db.Exec(`INSERT INTO users (firstName, lastName, email, bio, passwordHash)
	 VALUES ($1, $2, $3, $4, $5)`,
		user.Firstname, user.Lastname, user.Email, user.Bio, user.PasswordHash)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteUser(id int) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateUser(id int, req CreateUserRequest) error {
	_, err := s.db.Exec(`INSERT INTO users (firstName, lastName, email, bio, passwordHash)
	 VALUES ($1, $2, $3, $4, $5)`,
		req.Firstname, req.Lastname, req.Email, req.Bio, req.PasswordHash)

	if err != nil {
		return err
	}

	return nil
}

// CRUD OPERATIONS FOR POSTS
func (s *PostgresStore) GetAllPosts() ([]*Post, error) {
	return nil, nil
}

func (s *PostgresStore) GetPost(id int) (*Post, error) {
	return nil, nil
}

func (s *PostgresStore) CreatePost(req *CreatePostRequest) error {
	return nil
}

func (s *PostgresStore) DeletePost(id int) error {
	return nil
}

func (s *PostgresStore) UpdatePost(req int) error {
	return nil
}

// CRUD OPERATIONS FOR COMMENTS
func (s *PostgresStore) GetAllComments() ([]*Comment, error) {
	return nil, nil
}

func (s *PostgresStore) GetComment(id int) (*Comment, error) {
	return nil, nil
}

func (s *PostgresStore) CreateComment(req *CreateCommentRequest) error {
	return nil
}

func (s *PostgresStore) DeleteComment(id int) error {
	return nil
}

func (s *PostgresStore) UpdateComment(id int) error {
	return nil
}

// CRUD OPERATIONS FOR FOLLOWS
func (s *PostgresStore) CreateFollow(req *FollowRequest) error {
	return nil
}

func (s *PostgresStore) GetFollow(id int) (*Follow, error) {
	return nil, nil
}

func (s *PostgresStore) GetAllFollows() ([]*Follow, error) {
	return nil, nil
}

func (s *PostgresStore) DeleteFollow(id int) error {
	return nil
}

// FUNCTIONS FOR CREATING STRUCTS FROM SQL ROWS
func ScanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Bio,
		&user.PasswordHash,
		&user.Created_at,
	)

	return user, err
}

func ScanIntoPost(rows *sql.Rows) (*Post, error) {
	post := new(Post)
	err := rows.Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.MediaUrl,
		&post.Created_at,
	)

	return post, err
}
