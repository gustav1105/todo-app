# Todo App - gRPC with Go

This is a simple Todo List application built using Go, gRPC, Docker, and MySQL. It allows users to add, retrieve, and manage tasks via a gRPC server and a REST API gateway.

## Features

- Add tasks
- Retrieve tasks
- MySQL database for task persistence
- gRPC server for efficient communication
- REST API gateway built with Gin for easy interaction

## Technologies Used

- [Go](https://golang.org/) (gRPC and REST API)
- [gRPC](https://grpc.io/)
- [MySQL](https://www.mysql.com/)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- [Gin](https://gin-gonic.com/) for REST API
- [Uber FX](https://github.com/uber-go/fx) for dependency injection

## Prerequisites

- [Docker](https://www.docker.com/get-started) and Docker Compose
- [Go](https://golang.org/doc/install) (for local development)

## Getting Started

To run this project locally using Docker, follow these steps:

### 1. Clone the repository

git clone https://github.com/your-username/todo-app.git
cd todo-app

### 2. build and up

docker-compose up --build

### 3. http request

curl -X POST http://localhost:8080/tasks -d '{"title": "Buy groceries", "description": "Buy groceries for the week"}' -H "Content-Type: application/json"

curl http://localhost:8080/tasks

### 3. grpc requests

 ./todo-client get
 ./todo-client add --title "Buy planks" --description "nails no falls"

### 4. env vars needed
DB_USER=user
DB_PASSWORD=password
DB_HOST=mysql
DB_PORT=3306
DB_NAME=todo_db

### 5 spin up grpc server without env
./todo-client serve --db-user="your_db_user" --db-password="your_db_password" --db-host="your_db_host" --db-port="your_db_port" --db-name="your_db_name"
