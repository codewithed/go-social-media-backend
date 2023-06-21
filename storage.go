package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAllUsers() ([]*User, error)
	GetUser(id int) (*User, error)
	CreateUser(user *User) error
	DeleteUser(id int) error
	UpdateUser(id int, user *User) error
	GetAllPosts() ([]*Post, error)
	GetPost(id int) (*Post, error)
	CreatePost(req *CreatePostRequest) error
	DeletePost(id int) error
	UpdatePost(id int, req *CreatePostRequest) error
	GetAllComments() ([]*Comment, error)
	GetComment(id int) (*Comment, error)
	CreateComment(req *CreateCommentRequest) error
	DeleteComment(id int) error
	UpdateComment(id int, req *CreateCommentRequest) error
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
		FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE
	);

	CREATE TABLE comments (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		postID BIGINT NOT NULL,
		content VARCHAR(255) NOT NULL,
		created_at timestamptz NOT NULL,
		CONSTRAINT fk_userID FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE,
		CONSTRAINT fk_postID FOREIGN KEY (postID) REFERENCES posts (id) ON DELETE CASCADE
	);

	CREATE TABLE follows (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		followed_by_ID BIGINT NOT NULL,
		created_at timestamptz NOT NULL,
		CONSTRAINT fk_userID FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE,
		CONSTRAINT fk_followed_by_ID FOREIGN KEY (followed_by_ID) REFERENCES users (id) ON DELETE CASCADE
	);
	`

	_, err := s.db.Exec(query)
	return err
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

	return err
}

func (s *PostgresStore) DeleteUser(id int) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

func (s *PostgresStore) UpdateUser(id int, user *User) error {
	if user.Firstname != "" {
		_, err := s.db.Exec(`UPDATE users SET firstName = $1 WHERE id = $2`, user.Firstname, id)
		if err != nil {
			return err
		}
	}

	if user.Lastname != "" {
		_, err := s.db.Exec(`UPDATE users SET lastName = $1 WHERE id = $2`, user.Lastname, id)
		if err != nil {
			return err
		}
	}

	if user.Email != "" {
		_, err := s.db.Exec(`UPDATE users SET email = $1 WHERE id = $2`, user.Email, id)
		if err != nil {
			return err
		}
	}

	if user.Bio != "" {
		_, err := s.db.Exec(`UPDATE users SET bio = $1 WHERE id = $2`, user.Bio, id)
		if err != nil {
			return err
		}
	}

	if user.PasswordHash != "" {
		_, err := s.db.Exec(`UPDATE users SET firstname = $1 WHERE id = $2`, user.PasswordHash, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// CRUD OPERATIONS FOR POSTS
func (s *PostgresStore) GetAllPosts() ([]*Post, error) {
	rows, err := s.db.Query(`SELECT * FROM posts`)
	if err != nil {
		return nil, err
	}

	posts := []*Post{}
	for rows.Next() {
		post, err := ScanIntoPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostgresStore) GetPost(id int) (*Post, error) {
	rows, err := s.db.Query(`SELECT * FROM posts WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoPost(rows)
	}

	return nil, nil
}

func (s *PostgresStore) CreatePost(req *CreatePostRequest) error {
	_, err := s.db.Exec(`INSERT INTO posts (userID, mediaUrl, content) 
	VALUES ($1, $2, $3)`,
		req.UserID,
		req.MediaUrl,
		req.Content)

	return err
}

func (s *PostgresStore) DeletePost(id int) error {
	_, err := s.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	return err
}

func (s *PostgresStore) UpdatePost(id int, req *CreatePostRequest) error {
	if req.Content != "" {
		_, err := s.db.Exec(`UPDATE users SET content = $1 WHERE id = $2`, req.Content, id)
		if err != nil {
			return err
		}
	}

	if req.MediaUrl != "" {
		_, err := s.db.Exec(`UPDATE users SET mediaUrl = $1 WHERE id = $2`, req.MediaUrl, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// CRUD OPERATIONS FOR COMMENTS
func (s *PostgresStore) GetAllComments() ([]*Comment, error) {
	rows, err := s.db.Query(`SELECT * FROM comments`)
	if err != nil {
		return nil, err
	}

	comments := []*Comment{}
	for rows.Next() {
		comment, err := ScanIntoComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *PostgresStore) GetComment(id int) (*Comment, error) {
	rows, err := s.db.Query(`SELECT * FROM comments WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoComment(rows)
	}

	return nil, nil
}

func (s *PostgresStore) CreateComment(req *CreateCommentRequest) error {
	_, err := s.db.Exec(`INSERT INTO comments (text, userID, postID) 
	VALUES ($1, $2, $3)`,
		req.Text,
		req.UserID,
		req.PostID)

	return err
}

func (s *PostgresStore) DeleteComment(id int) error {
	_, err := s.db.Exec(`DELETE FROM comments WHERE id = $1`, id)
	return err
}

func (s *PostgresStore) UpdateComment(id int, req *CreateCommentRequest) error {
	if req.Text != "" {
		_, err := s.db.Exec(`UPDATE comments SET text = $1 WHERE id = $2`,
			req.Text, id)

		if err != nil {
			return err
		}
	}

	return nil
}

// CRUD OPERATIONS FOR FOLLOWS
func (s *PostgresStore) CreateFollow(req *FollowRequest) error {
	_, err := s.db.Exec(`INSERT INTO follows (userID, followerID) 
	VALUES ($1, $2)`,
		req.UserID,
		req.FollowerID)

	return err
}

func (s *PostgresStore) GetFollow(id int) (*Follow, error) {
	rows, err := s.db.Query(`SELECT * FROM follows WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoFollow(rows)
	}

	return nil, nil
}

func (s *PostgresStore) GetAllFollows() ([]*Follow, error) {
	rows, err := s.db.Query(`SELECT * FROM follows`)
	if err != nil {
		return nil, err
	}

	follows := []*Follow{}
	for rows.Next() {
		follow, err := ScanIntoFollow(rows)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	return follows, nil
}

func (s *PostgresStore) DeleteFollow(id int) error {
	_, err := s.db.Exec(`DELETE FROM comments WHERE id = $1`, id)
	return err
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

func ScanIntoComment(rows *sql.Rows) (*Comment, error) {
	comment := new(Comment)
	err := rows.Scan(
		&comment.ID,
		&comment.Text,
		&comment.UserID,
		&comment.PostID,
		&comment.Created_at,
	)

	return comment, err
}

func ScanIntoFollow(rows *sql.Rows) (*Follow, error) {
	follow := new(Follow)
	err := rows.Scan(
		&follow.ID,
		&follow.UserID,
		&follow.FollowerID,
		&follow.Created_at,
	)

	return follow, err
}
