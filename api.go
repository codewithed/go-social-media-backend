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
	r.HandleFunc("/users", makeHttpHandlerFunc(s.handleUsers))
	r.HandleFunc("/users/{id}", makeHttpHandlerFunc(s.handleUsersByID))
	r.HandleFunc("/posts", makeHttpHandlerFunc(s.handlePosts))
	r.HandleFunc("/posts/{id}", makeHttpHandlerFunc(s.handlePostsByID))
	http.ListenAndServe(s.ListenAddr, r)
}

func (s *ApiServer) handleUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		users, err := s.Store.GetAllUsers()
		if err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, users)
	}

	if r.Method == "POST" {
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
	return nil
}

func (s *ApiServer) handleUsersByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		user, err := s.Store.GetUser(id)
		if err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, user)
	}

	if r.Method == "PUT" || r.Method == "PATCH" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		req := new(CreateUserRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return err
		}

		user, err := NewUser(req)
		if err != nil {
			return err
		}

		if err := s.Store.UpdateUser(id, user); err != nil {
			return err
		}
		return WriteJson(w, http.StatusOK, user)
	}

	if r.Method == "DELETE" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		if err := s.Store.DeleteUser(id); err != nil {
			return err
		}
		deletedMsg := fmt.Sprintf("User with id: %d deleted successfully", id)
		return WriteJson(w, http.StatusOK, deletedMsg)
	}
	return nil
}

func (s *ApiServer) handlePosts(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		users, err := s.Store.GetAllPosts()
		if err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, users)
	}

	if r.Method == "POST" {
		req := new(CreatePostRequest)

		if err := s.Store.CreatePost(req); err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, req)
	}
	return nil
}

func (s *ApiServer) handlePostsByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
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

	if r.Method == "PUT" || r.Method == "PATCH" {
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

	if r.Method == "DELETE" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		if err := s.Store.DeletePost(id); err != nil {
			return err
		}
		deleteMsg := fmt.Sprintf("Post with id: %d deleted successfully", id)
		return WriteJson(w, http.StatusOK, deleteMsg)
	}
	return nil
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

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}
