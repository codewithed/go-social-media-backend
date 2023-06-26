package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ApiServer struct {
	ListenAddr string
	Store      Storage
}

func NewApiServer(addr string, store Storage) (*ApiServer, error) {
	return &ApiServer{
		ListenAddr: addr,
		Store:      store,
	}, nil
}

func (s *ApiServer) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.HandleFunc("/signup", makeHttpHandlerFunc(s.handleSignUp))
	r.HandleFunc("/login", makeHttpHandlerFunc(s.handleLogin))
	r.HandleFunc("/{username}", makeHttpHandlerFunc(s.handleUsersByName))
	r.HandleFunc("/{username}/posts", makeHttpHandlerFunc(s.handleGetUserPosts))
	r.HandleFunc("/{username}/followers", makeHttpHandlerFunc(s.handleGetFollowers))
	r.HandleFunc("/{username}/following", makeHttpHandlerFunc(s.handleGetFollowing))
	r.HandleFunc("/{username}/follow", makeHttpHandlerFunc(s.handleFollow))
	r.HandleFunc("/posts", makeHttpHandlerFunc(s.handleCreatePost))
	r.HandleFunc("/posts/{id}", makeHttpHandlerFunc(s.handlePostsByID))
	r.HandleFunc("/posts/{id}/likes", makeHttpHandlerFunc(s.handlePostlikes))
	r.HandleFunc("/posts/{id}/comments", makeHttpHandlerFunc(s.handleGetCommentsFromPost))
	r.HandleFunc("/posts/{postID}/comments/{commentID}", makeHttpHandlerFunc(s.handleGetCommentsFromPost))
	http.ListenAndServe(s.ListenAddr, r)
}

func (s *ApiServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJson(w, http.StatusBadRequest, &ApiError{Error: "Unexpected method"})
	}
	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	user, err := s.Store.GetUser(req.UserName)
	if err != nil {
		return err
	}

	if !user.ValidPassword(req.Password) {
		return WriteJson(w, http.StatusBadRequest, &ApiError{Error: "Access denied"})
	}

	token, err := CreateJWT(user)
	if err != nil {
		return err
	}

	res := LoginResponse{
		UserName: user.UserName,
		Token:    token,
	}

	return WriteJson(w, http.StatusOK, res)
}

func (s *ApiServer) handleUsersByName(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUserProfile(w, r)
	}

	if r.Method == "PUT" || r.Method == "PATCH" {
		return s.handleUpdateUser(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteUser(w, r)
	}
	return nil
}

func (s *ApiServer) handleGetFollowers(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("Unexpected method: %v", r.Method)
	}

	username := getUserName(r)
	followers, err := s.Store.GetFollowers(username)
	if err != nil {
		return fmt.Errorf("Couldn't get followers")
	}

	return WriteJson(w, http.StatusOK, followers)
}

func (s *ApiServer) handleGetFollowing(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("Unexpected method: %v", r.Method)
	}

	username := getUserName(r)
	following, err := s.Store.GetFollowing(username)
	if err != nil {
		return fmt.Errorf("Couldn't get followers")
	}

	return WriteJson(w, http.StatusOK, following)
}

func (s *ApiServer) handleFollow(w http.ResponseWriter, r *http.Request) error {

	req := new(FollowRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if r.Method == "POST" {
		err := s.Store.CreateFollow(req)
		if err != nil {
			return fmt.Errorf("Failed to follow user with id: %v", req.FollowingID)
		}

		return WriteJson(w, http.StatusOK, fmt.Sprintf("Followed user with id: %v", req.FollowingID))
	}

	if r.Method == "DELETE" {
		err := s.Store.DeleteFollow(req)
		if err != nil {
			return fmt.Errorf("Failed to unfollow user with id: %v", req.FollowingID)
		}

		return WriteJson(w, http.StatusOK, fmt.Sprintf("Unfollowed user with id: %v", req.FollowingID))
	}

	return nil
}

func (s *ApiServer) handleGetUserPosts(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("Unexpected method")
	}

	username := getUserName(r)
	posts, err := s.Store.GetUserPosts(username)
	if err != nil {
		return fmt.Errorf("Error getting user posts: %v", err)
	}

	return WriteJson(w, http.StatusOK, posts)
}

func (s *ApiServer) handlePostsByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetPostByID(w, r)
	}

	if r.Method == "PUT" || r.Method == "PATCH" {
		return s.handleUpdatePostByID(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeletePostByID(w, r)
	}
	return nil
}

func (s *ApiServer) handleCommentsByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetCommentByID(w, r)
	}

	if r.Method == "UPDATE" {
		return s.handleUpdateCommentByID(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteCommentByID(w, r)
	}
	return nil
}

func (s *ApiServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("Unexpected method %s", r.Method)
	}
	req := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	user, err := NewUser(req)
	if err != nil {
		return err
	}

	if err := s.Store.CreateUser(user); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, user)
}

func (s *ApiServer) handleGetUserProfile(w http.ResponseWriter, r *http.Request) error {
	username := getUserName(r)

	user, err := s.Store.GetUserProfile(username)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, user)
}

func (s *ApiServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	username := getUserName(r)

	req := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	user, err := NewUser(req)
	if err != nil {
		return err
	}

	if err := s.Store.UpdateUser(username, user); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, user)
}

func (s *ApiServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	username := getUserName(r)

	if err := s.Store.DeleteUser(username); err != nil {
		return err
	}
	deletedMsg := fmt.Sprintf("User %s deleted successfully", username)
	return WriteJson(w, http.StatusOK, &ApiError{Error: deletedMsg})
}

func (s *ApiServer) handleCreatePost(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("Unexpected method %s", r.Method)
	}
	req := new(CreatePostRequest)

	if err := s.Store.CreatePost(req); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, req)
}

func (s *ApiServer) handleGetPostByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	post, err := s.Store.GetPost(id)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, post)
}

func (s *ApiServer) handleUpdatePostByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	req := new(CreatePostRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if err := s.Store.UpdatePost(id, req); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, req)
}

func (s *ApiServer) handleDeletePostByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.Store.DeletePost(id); err != nil {
		return err
	}
	deletedMsg := fmt.Sprintf("Post with id: %d deleted successfully", id)
	return WriteJson(w, http.StatusOK, &ApiError{Error: deletedMsg})
}

func (s *ApiServer) handleGetCommentsFromPost(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return fmt.Errorf("Invalid id")
	}

	comments, err := s.Store.GetCommentsFromPost(id)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, comments)
}

func (s *ApiServer) handleCreateComment(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateCommentRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	if err := s.Store.CreateComment(req); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, req)
}

func (s *ApiServer) handleUpdateCommentByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	req := new(CreateCommentRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	if err := s.Store.UpdateComment(id, req); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, req)
}

func (s *ApiServer) handleGetCommentByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	comment, err := s.Store.GetComment(id)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, comment)
}

func (s *ApiServer) handleDeleteCommentByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.Store.DeleteComment(id); err != nil {
		return err
	}
	deletedMsg := fmt.Sprintf("Deleted comment %v succesfully", id)
	return WriteJson(w, http.StatusOK, &ApiError{Error: deletedMsg})
}

func (s *ApiServer) handlePostlikes(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleLikePost(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleUnlikePost(w, r)
	}
	return nil
}

func (s *ApiServer) handleLikePost(w http.ResponseWriter, r *http.Request) error {
	req := new(LikeRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return fmt.Errorf("Error decoding like request")
	}

	err := s.Store.LikePost(LikeRequest)
	if err != nil {
		return fmt.Errorf("Failed to like post")
	}
	return WriteJson(w, http.StatusOK, fmt.Sprintf("Liked post:%s  successfully", req.ResourceID))
}

func (s *ApiServer) handleUnlikePost(w http.ResponseWriter, r *http.Request) error {
	req := new(LikeRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return fmt.Errorf("Error decoding like request")
	}

	err := s.Store.UnlikePost(LikeRequest)
	if err != nil {
		return fmt.Errorf("Failed to like post")
	}
	return WriteJson(w, http.StatusOK, fmt.Sprintf("Unliked post:%s successfully", req.ResourceID))
}

func makeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, &ApiError{Error: err.Error()})
		}
	}
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func getID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("Invalid id: %v", id)
	}
	return id, nil
}

func getUserName(r *http.Request) string {
	username := chi.URLParam(r, "username")
	return username
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}
