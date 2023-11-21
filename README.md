# go-social-media-backend

A social media backend built using Go

## Technologies

- Go standard library
- Chi router
- PostgreSQL database
- Docker
- JWT tokens

## Features include:
- CRUD operations for users, posts, comments, follows, likes, etc
- Getting user profiles
- Certain actions require authorization to perform: 
 For example:
  - A user has to be the author of a post to delete or update it.
  - Users have to create an account before they can like or comment on a post.
