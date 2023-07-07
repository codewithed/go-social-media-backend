package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetUser(username string) (*User, error)
	GetUserProfile(username string) (*UserProfile, error)
	CreateUser(user *User) error
	DeleteUser(username string) error
	UpdateUser(username string, user *UpdateUserRequest) error
	GetUserPosts(username string) ([]*Post, error)
	GetPost(id int) (*Post, error)
	CreatePost(req *CreatePostRequest) error
	DeletePost(id int) error
	UpdatePost(id int, req *CreatePostRequest) error
	GetCommentsFromPost(postID int) ([]*Comment, error)
	GetComment(id int) (*Comment, error)
	CreateComment(req *CreateCommentRequest) error
	DeleteComment(id int) error
	UpdateComment(id int, req *CreateCommentRequest) error
	GetFollowers(username string) ([]string, error)
	GetFollowing(username string) ([]string, error)
	CreateFollow(req *FollowRequest) error
	DeleteFollow(req *FollowRequest) error
	LikePost(req *LikeRequest) error
	UnlikePost(req *LikeRequest) error
	LikeComment(req *LikeRequest) error
	UnlikeComment(req *LikeRequest) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASS"))
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
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		userName VARCHAR(25) NOT NULL UNIQUE,
		name VARCHAR(225),
		email VARCHAR(255) NOT NULL,
		bio VARCHAR(255),
		passwordHash VARCHAR(1000) NOT NULL,
		created_at timestamptz NOT NULL DEFAULT timezone('UTC', now())
	);
	
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		content VARCHAR(255) NOT NULL,
		mediaUrl VARCHAR(10000),
		created_at timestamptz NOT NULL DEFAULT timezone('UTC', now()),
		FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		postID BIGINT NOT NULL,
		content VARCHAR(255) NOT NULL,
		created_at timestamptz NOT NULL DEFAULT timezone('UTC', now()),
		FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (postID) REFERENCES posts (id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS follows (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		followerID BIGINT NOT NULL,
		created_at timestamptz NOT NULL DEFAULT timezone('UTC', now()),
		FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (followerID) REFERENCES users (id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS post_likes (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		postID BIGINT NOT NULL,
		created_at timestamptz NOT NULL DEFAULT timezone('UTC', now()),
		FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (postID) REFERENCES posts (id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS comment_likes (
		id SERIAL PRIMARY KEY,
		userID BIGINT NOT NULL,
		commentID BIGINT NOT NULL,
		created_at timestamptz NOT NULL DEFAULT timezone('UTC', now()),
		FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (commentID) REFERENCES comments (id) ON DELETE CASCADE
	);`

	_, err := s.db.Exec(query)
	return err
}

// CRUD OPERATIONS FOR USERS
func (s *PostgresStore) GetUser(name string) (*User, error) {
	user_id, err := s.getUserIDFromUserName(name)
	if err != nil || user_id == 0 {
		return nil, fmt.Errorf("user %s not found", name)
	}

	rows, err := s.db.Query(`SELECT * FROM users WHERE id = $1`, user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoUser(rows)
	}

	return nil, nil
}

func (s *PostgresStore) GetUserProfile(username string) (*UserProfile, error) {
	user_id, err := s.getUserIDFromUserName(username)
	if err != nil || user_id == 0 {
		return nil, fmt.Errorf("user %s not found", username)
	}

	// get user info
	user_info, err := s.db.Query(`SELECT id, userName, name, bio FROM users WHERE id = $1`, user_id)
	if err != nil {
		return nil, err
	}
	defer user_info.Close()

	profile := new(UserProfile)
	if user_info.Next() {
		if err := user_info.Scan(
			&profile.UserID,
			&profile.UserName,
			&profile.Name,
			&profile.Bio,
		); err != nil {
			return nil, err
		}
	}

	// get user stats
	user_stats, err := s.db.Query(`SELECT
    (SELECT COUNT(*) FROM posts WHERE userID = $1) AS post_count,
    (SELECT COUNT(*) FROM follows WHERE userID = $1) AS follower_count,
    (SELECT COUNT(*) FROM follows WHERE followerID = $1) AS following_count;`, profile.UserID)
	if err != nil {
		return nil, err
	}
	defer user_stats.Close()

	if user_stats.Next() {
		if err := user_stats.Scan(
			&profile.Posts,
			&profile.Followers,
			&profile.Following,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user stats")
		}
	}

	return profile, nil
}

func (s *PostgresStore) CreateUser(user *User) error {
	_, err := s.db.Exec(`INSERT INTO users (userName, name, email, bio, passwordHash)
	 VALUES ($1, $2, $3, $4, $5)`,
		user.UserName, user.Name, user.Email, user.Bio, user.PasswordHash)

	return err
}

func (s *PostgresStore) DeleteUser(username string) error {
	user_id, err := s.getUserIDFromUserName(username)
	if err != nil || user_id == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	_, err = s.db.Exec(`DELETE FROM users WHERE id = $1`, user_id)
	return err
}

func (s *PostgresStore) UpdateUser(username string, user *UpdateUserRequest) error {
	user_id, err := s.getUserIDFromUserName(username)
	if err != nil {
		return err
	}

	if user.UserName != "" {
		_, err := s.db.Exec(`UPDATE users SET userName = $1 WHERE id = $2 AND userName != $1`,
			user.UserName, user_id)
		if err != nil {
			return err
		}
	}

	if user.Name != "" {
		_, err := s.db.Exec(`UPDATE users SET name = $1 WHERE id = $2 AND name != $1`,
			user.Name, user_id)
		if err != nil {
			return err
		}
	}

	if user.Email != "" {
		_, err := s.db.Exec(`UPDATE users SET email = $1 WHERE id = $2 AND email != $1`,
			user.Email, user_id)
		if err != nil {
			return err
		}
	}

	if user.Bio != "" {
		_, err := s.db.Exec(`UPDATE users SET bio = $1 WHERE id = $2 AND bio != $1`,
			user.Bio, user_id)
		if err != nil {
			return err
		}
	}

	if user.PasswordHash != "" {
		_, err := s.db.Exec(`UPDATE users SET passwordHash = $1 WHERE id = $2 AND passwordHash != $1`,
			user.PasswordHash, user_id)
		if err != nil {
			return err
		}
	}

	return nil
}

// CRUD OPERATIONS FOR POSTS
func (s *PostgresStore) GetUserPosts(username string) ([]*Post, error) {
	user_id, err := s.getUserIDFromUserName(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID from username")
	}

	rows, err := s.db.Query(`SELECT * FROM posts WHERE userID = $1`, user_id)
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
func (s *PostgresStore) GetCommentsFromPost(postID int) ([]*Comment, error) {
	rows, err := s.db.Query(`SELECT * FROM comments WHERE postID = $1`, postID)
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
		req.UserName,
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
func (s *PostgresStore) GetFollowers(username string) ([]string, error) {
	id, err := s.getUserIDFromUserName(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_id from username: %v", err)
	}

	rows, err := s.db.Query(`SELECT userName FROM follows WHERE userID = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers row")
	}

	followers := []string{}
	for rows.Next() {
		follower := ""
		err := rows.Scan(&follower)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user")
		}

		followers = append(followers, follower)
	}

	return followers, nil
}

func (s *PostgresStore) GetFollowing(username string) ([]string, error) {
	id, err := s.getUserIDFromUserName(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_id from username: %v", err)
	}

	rows, err := s.db.Query(`SELECT userName FROM follows WHERE followerID = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get following row")
	}

	following := []string{}
	for rows.Next() {
		follow := ""
		err := rows.Scan(&follow)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user")
		}

		following = append(following, follow)
	}

	return following, nil
}

func (s *PostgresStore) CreateFollow(req *FollowRequest) error {
	_, err := s.db.Exec(`INSERT INTO follows (userID, followerID) 
	VALUES ($1, $2)`,
		req.UserID,
		req.FollowingID)

	return err
}

func (s *PostgresStore) DeleteFollow(req *FollowRequest) error {
	_, err := s.db.Exec(`DELETE FROM follows WHERE userID = $1 AND followerID = $2`,
		req.UserID,
		req.FollowingID)
	return err
}

func (s *PostgresStore) LikePost(req *LikeRequest) error {
	_, err := s.db.Exec(`INSERT INTO post_likes (userID, postID) VALUES($1, $2)`, req.UserID, req.ResourceID)
	return err
}

func (s *PostgresStore) UnlikePost(req *LikeRequest) error {
	_, err := s.db.Exec(`DELETE FROM post_likes WHERE userID = $1 AND postID = $2`, req.UserID, req.ResourceID)
	return err
}

func (s *PostgresStore) LikeComment(req *LikeRequest) error {
	_, err := s.db.Exec(`INSERT INTO comment_likes (userID, commentID) VALUES($1, $2)`, req.UserID, req.ResourceID)
	return err
}

func (s *PostgresStore) UnlikeComment(req *LikeRequest) error {
	_, err := s.db.Exec(`DELETE FROM comment_likes WHERE userID = $1 AND commentID = $2`, req.UserID, req.ResourceID)
	return err
}

// FUNCTIONS FOR CREATING STRUCTS FROM SQL ROWS
func ScanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.UserName,
		&user.Name,
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

// HELPER FUNCTIONS
func (s *PostgresStore) getUserIDFromUserName(username string) (int, error) {
	var id int
	idrow, err := s.db.Query(`SELECT id FROM users WHERE username = $1`, username)
	if err != nil {
		return id, fmt.Errorf("failed to get userID from username: %v", err)
	}
	defer idrow.Close()

	if idrow.Next() {
		err := idrow.Scan(&id)
		if err != nil {
			return id, err
		}
	}

	return id, nil
}
