# go-social-media-backend

A social media backend built using Go

## Technologies used

- Go standard library
- Chi router
- PostgreSQL database
- Docker
- JWT tokens

## Features
- CRUD operations for users, posts, comments, follows, likes, etc
- Getting user profiles (including their bio, number of posts, followers, and following)
- Certain actions require authorization to perform: 
 For example:
  - A user has to be the author of a post to delete or update it.
  - Users have to create an account before they can like or comment on a post.

 ## Local setup
 ```
git clone https://github.com/codewithed/go-social-media-backend.git
cd go-social-media backend
```
### Build command
```
make build
```
### Run command
```
make run
```
