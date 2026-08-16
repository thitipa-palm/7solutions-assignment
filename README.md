# User Management API

A RESTful user management API built with Go, Fiber, MongoDB, and JWT authentication.

## Features

- User registration and login
- JWT authentication using HS256
- Create, list, get, update, and delete users
- Password hashing using bcrypt
- Request validation
- Unique email enforcement
- Request logging
- Background user count logging every 10 seconds
- Graceful shutdown
- Unit tests with a mocked repository
- Docker and Docker Compose support

## Technologies

- Go 1.25
- Fiber
- MongoDB
- MongoDB Go Driver
- JWT
- bcrypt
- go-playground/validator
- Docker

## Project Structure

```text
.
├── cmd/api                 # Application entry point
├── internal
│   ├── config              # Environment configuration
│   ├── handler             # HTTP handlers and request DTOs
│   ├── middleware          # Authentication and logging middleware
│   ├── model               # Domain models
│   ├── repository          # Repository interface and MongoDB implementation
│   ├── router              # HTTP route configuration
│   └── service             # Business logic
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

The service depends on the `UserRepository` interface instead of a concrete MongoDB implementation. This keeps the business logic separated from the database and makes unit testing easier.

## Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Example configuration:

```env
PORT=8080

MONGO_USERNAME=admin
MONGO_PASSWORD=dev-password
MONGO_URI=mongodb://admin:dev-password@localhost:27017/?authSource=admin
MONGO_DATABASE=user_management

JWT_SECRET=change-this-secret
JWT_EXPIRATION=24h
```

Do not commit the real `.env` file because it may contain sensitive values.

## Run with Docker Compose

Build and start the API and MongoDB:

```bash
docker compose up --build
```

The API will be available at:

```text
http://localhost:8080
```

Check the application health:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

Stop the containers:

```bash
docker compose down
```

MongoDB data is stored in a named Docker volume and remains available after the containers are stopped.

## Run Locally

Start MongoDB:

```bash
docker compose up -d mongodb
```

Install dependencies:

```bash
go mod download
```

Run the API:

```bash
go run ./cmd/api
```

## API Endpoints

| Method | Endpoint                | Authentication | Description             |
| ------ | ----------------------- | -------------: | ----------------------- |
| GET    | `/health`               |             No | Health check            |
| POST   | `/api/v1/auth/register` |             No | Register a user         |
| POST   | `/api/v1/auth/login`    |             No | Login and receive a JWT |
| POST   | `/api/v1/users/`        |            Yes | Create a user           |
| GET    | `/api/v1/users/`        |            Yes | List users              |
| GET    | `/api/v1/users/:id`     |            Yes | Get a user by ID        |
| PATCH  | `/api/v1/users/:id`     |            Yes | Update a user           |
| DELETE | `/api/v1/users/:id`     |            Yes | Delete a user           |

Protected routes require the following header:

```text
Authorization: Bearer <access-token>
```

## API Examples

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Palm",
    "email": "palm@example.com",
    "password": "password123"
  }'
```

Expected status:

```text
201 Created
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "palm@example.com",
    "password": "password123"
  }'
```

Example response:

```json
{
  "accessToken": "<JWT>",
  "tokenType": "Bearer"
}
```

### List Users

```bash
curl http://localhost:8080/api/v1/users/ \
  -H "Authorization: Bearer <access-token>"
```

### Get User by ID

```bash
curl http://localhost:8080/api/v1/users/<user-id> \
  -H "Authorization: Bearer <access-token>"
```

### Update User

```bash
curl -X PATCH http://localhost:8080/api/v1/users/<user-id> \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Palm Updated",
    "email": "palm.updated@example.com"
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/v1/users/<user-id> \
  -H "Authorization: Bearer <access-token>"
```

Expected status:

```text
204 No Content
```

## Validation and Error Responses

The API validates required fields, email format, and password length.

Example error response:

```json
{
  "error": "invalid registration data"
}
```

Common HTTP status codes:

| Status | Description                       |
| -----: | --------------------------------- |
|    200 | Request succeeded                 |
|    201 | Resource created                  |
|    204 | Resource deleted                  |
|    400 | Invalid request                   |
|    401 | Missing or invalid authentication |
|    404 | User not found                    |
|    409 | Email already exists              |
|    500 | Internal server error             |

## Tests

Run all tests:

```bash
go test ./...
```

Run service tests with detailed output and coverage:

```bash
go test ./internal/service -v -cover
```

The service tests use a mocked repository and do not require a real MongoDB connection.

## Background Task

The application runs a background goroutine every 10 seconds to count users and write the result to the application log.

The goroutine is stopped during graceful shutdown before the MongoDB connection is closed.

## Graceful Shutdown

The application listens for `SIGINT` and `SIGTERM`. When a shutdown signal is received, it:

1. Stops accepting new HTTP requests.
2. Waits for active requests to finish.
3. Stops the background user count monitor.
4. Disconnects from MongoDB.


## Lottery Search System Design

Part 2 of this assignment is a system design proposal only. No implementation code is required.

- [English Version](./docs/lottery-search-system.en.md)
- [Thai Version](./docs/lottery-search-system.th.md)