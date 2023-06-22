package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	UserName     string    `json:"userName"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Bio          string    `json:"bio"`
	PasswordHash string    `json:"-"`
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
	UserName string `json:"userName"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Password string `json:"password"`
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

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserName string `json:"userName"`
	Token    string `json:"token"`
}

func NewUser(req *CreateUserRequest) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MaxCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           int64(rand.Intn(10000)),
		UserName:     req.UserName,
		Name:         req.Name,
		Email:        req.Email,
		Bio:          req.Bio,
		PasswordHash: string(passwordHash),
		Created_at:   time.Now().UTC(),
	}, nil
}

func (user *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pw)) == nil
}
