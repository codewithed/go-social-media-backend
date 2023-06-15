package main

import "time"

type User struct {
	ID           int64     `json:"id"`
	Firstname    string    `json:"firstName"`
	Lastname     string    `json:"lastName"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordhash"`
	Created_at   time.Time `json:"createdAt"`
}

type Post struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userID"`
	MediaUrl   string    `json:"mediaUrl"`
	Content    string    `json:"content"`
	Created_at time.Time `json:"createdAt"`
}

type Comments struct {
	ID         int64     `json:"id"`
	Text       string    `json:"text"`
	UserID     int64     `json:"userID"`
	PostID     int64     `json:"postID"`
	Created_at time.Time `json:"createdAt"`
}

type Follow struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userID"`
	FollowerID int64     `json:"followerID"`
	Created_at time.Time `json:"createdAt"`
}

type CreateUserRequest struct {
	Firstname  string    `json:"firstName"`
	Lastname   string    `json:"lastName"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_at time.Time `json:"createdAt"`
}

type CreatePostRequest struct {
	UserID   int64  `json:"userID"`
	MediaUrl string `json:"mediaUrl"`
	Content  string `json:"content"`
}

type CreateCommentRequest struct {
	Text   string `json:"text"`
	UserID int64  `json:"userID"`
}

type FollowRequest struct {
	UserID     int64 `json:"userID"`
	FollowerID int64 `json:"followerID"`
}
