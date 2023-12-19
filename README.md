# go-social-media-backend

A social media backend built using Go

## Technologies used

- Go standard library
- Chi router
- PostgreSQL database
- BCrypt
- Docker
- JWT tokens

## Features include

- User signup and login
- Creating, deleting and updating posts
- Liking and unliking posts
- Following and unfollowing other users
- Commenting and deleting comments
- Viewing a user's profile (number of posts, followers, following)
- Liking and unliking comments
- CRUD endpoints for users, posts, following, likes
- Certain actions require authorization to perform:
 For example:
  - A user has to be the author of a post to delete or update it.
  - Users have to create an account before they can like or comment on a post.

## Local setup

``` bash
git clone https://github.com/codewithed/go-social-media-backend.git
cd go-social-media backend
```

### Setup Postgres Database

``` bash
make postgres
```

### Build command

``` bash
make build
```

### Run command

``` bash
make run
```
