# Go Social Media Backend

A robust, RESTful social media backend API built with Go, featuring user authentication, social interactions, and comprehensive CRUD operations.

## 🚀 Features

### Core Functionality
- **User Management**: Complete signup and login system with JWT authentication
- **Post Operations**: Create, read, update, and delete posts with proper authorization
- **Social Interactions**: 
  - Like/unlike posts and comments
  - Follow/unfollow users
  - Comment on posts with full CRUD operations
- **User Profiles**: View comprehensive user statistics (posts, followers, following counts)
- **Authorization & Security**: Role-based access control with JWT tokens and BCrypt password hashing

### API Endpoints
- **Users**: Registration, authentication, profile management
- **Posts**: Full CRUD operations with author-based permissions  
- **Comments**: Create and delete comments on posts
- **Likes**: Toggle likes on posts and comments
- **Following**: Manage user relationships and social connections

## 🛠️ Technology Stack

- **Language**: Go (Golang)
- **Router**: Chi - lightweight, idiomatic HTTP router
- **Database**: PostgreSQL with optimized queries
- **Authentication**: JWT (JSON Web Tokens)
- **Password Security**: BCrypt hashing algorithm
- **Containerization**: Docker for easy deployment
- **Architecture**: RESTful API design principles

## 📋 Prerequisites

Before running this application, ensure you have the following installed:

- **Docker** (v20.0 or higher)
- **Docker Compose** (optional, for easier setup)
- **Git** for cloning the repository

## 🔧 Local Setup

### 1. Clone the Repository
```bash
git clone https://github.com/codewithed/go-social-media-backend.git
cd go-social-media-backend
```

### 2. Environment Setup
Create a `.env` file in the root directory (if not already present):
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=social_media_db

# JWT Configuration  
JWT_SECRET=your_jwt_secret_key

# Server Configuration
PORT=8080
```

### 3. Database Setup
Start the PostgreSQL database using Docker:
```bash
make postgres
```

### 4. Build the Application
```bash
make build
```

### 5. Run the Application
```bash
make run
```

The API will be available at `http://localhost:8080`

## 📖 API Documentation

### Authentication Endpoints
- `POST /signup` - Register a new user
- `POST /login` - Authenticate user and receive JWT token

### User Profile Endpoints
- `GET /{username}` - Get user profile information
- `PUT /{username}` - Update user profile (owner only)
- `PATCH /{username}` - Partially update user profile (owner only)
- `DELETE /{username}` - Delete user account (owner only)

### User Social Endpoints
- `GET /{username}/followers` - Get user's followers list
- `GET /{username}/following` - Get users that this user follows
- `POST /{username}/follow` - Follow a user (authenticated)
- `DELETE /{username}/unfollow` - Unfollow a user (authenticated)

### User Posts Endpoints
- `GET /{username}/posts` - Get all posts by a specific user (authenticated)
- `POST /{username}/posts` - Create a new post (owner only)

### Post Management Endpoints
- `GET /posts/{id}` - Get a specific post by ID
- `PUT /posts/{id}` - Update a post (author only)
- `PATCH /posts/{id}` - Partially update a post (author only)
- `DELETE /posts/{id}` - Delete a post (author only)

### Post Interaction Endpoints
- `GET /posts/{id}/likes` - Get list of users who liked the post (authenticated)
- `POST /posts/{id}/like` - Like a post (authenticated, requires ?userID= query param)
- `DELETE /posts/{id}/unlike` - Unlike a post (authenticated, requires ?userID= query param)

### Comment Endpoints
- `GET /posts/{id}/comments` - Get all comments for a post (authenticated)
- `POST /posts/{id}/comments` - Add a comment to a post (authenticated)
- `GET /comments/{id}` - Get a specific comment by ID
- `PUT /comments/{id}` - Update a comment (author only)
- `PATCH /comments/{id}` - Partially update a comment (author only)
- `DELETE /comments/{id}` - Delete a comment (author only)

### Comment Interaction Endpoints
- `POST /comments/{id}/like` - Like a comment (authenticated, requires ?userID= query param)
- `DELETE /comments/{id}/unlike` - Unlike a comment (authenticated, requires ?userID= query param)

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication system
- **Password Hashing**: BCrypt algorithm for secure password storage
- **Authorization Middleware**: Route-level permission checking
- **Owner-based Permissions**: Users can only modify their own content
- **Input Validation**: Comprehensive request validation and sanitization

---

**Built with ❤️ using Go**
