package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	Firstname    string    `json:"firstName"`
	Lastname     string    `json:"lastName"`
	Email        string    `json:"email"`
	Bio          string    `json:"bio"`
	PasswordHash string    `json:"passwordhash"`
	Created_at   time.Time `json:"createdAt"`
}

type Post struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userID"`
	MediaUrl   string    `json:"mediaUrl"`
	Content    string    `json:"content"`
	Created_at time.Time `json:"createdAt"`
	Edited_at  time.Time `json:"editedAt"`
}

type Comment struct {
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
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Password  string `json:"password"`
}

type CreatePostRequest struct {
	UserID   int64  `json:"userID"`
	Content  string `json:"content"`
	MediaUrl string `json:"mediaUrl"`
}

type CreateCommentRequest struct {
	Text   string `json:"text"`
	UserID int64  `json:"userID"`
	PostID int64  `json:"postID"`
}

type FollowRequest struct {
	UserID     int64 `json:"userID"`
	FollowerID int64 `json:"followerID"`
}

func NewUser(req *CreateUserRequest) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MaxCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           int64(rand.Intn(10000)),
		Firstname:    req.Firstname,
		Lastname:     req.Lastname,
		Email:        req.Email,
		Bio:          req.Bio,
		PasswordHash: string(passwordHash),
		Created_at:   time.Now().UTC(),
	}, nil
}
