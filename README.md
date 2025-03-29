# Lost & Found - Auth Service

## Overview
The User Service handles user authentication, profile management, and access control.

## Features
- User registration & authentication (JWT/OAuth2)
- Profile management
- Role-based access control (RBAC)

## Tech Stack
- Language: Golang
- Database: PostgreSQL
- Authentication: JWT / OAuth2
- API: REST + GraphQL

## Installation
```bash
git clone https://github.com/your-org/lostfound-user-service.git
cd lostfound-user-service
go mod tidy
go run main.go
```

## Usage
API runs on http://localhost:8080
GraphQL Playground: http://localhost:8080/graphql


## API Endpoints
Method	Endpoint	Description
POST	/users/signup	Register a new user
POST	/users/login	User login
GET	/users/:id	Get user details