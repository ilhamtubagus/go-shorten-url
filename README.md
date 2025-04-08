# Go URL Shortener

This is a URL shortener application built with Go. It allows users to create shortened URLs, manage them, and redirect to the original URLs.

## Tech Stack

- **Go**: The main programming language used for the backend.
- **MongoDB**: Used as the primary database for storing shortened URL data.
- **Redis**: Utilized for caching to improve performance.
- **HTML Templates**: For rendering the frontend pages.
- **httprouter**: A lightweight HTTP router for handling routes.

### Key Libraries

- `github.com/joho/godotenv`: For loading environment variables from a .env file.
- `github.com/ilhamtubagus/goenv`: Unmarshal environment variables into struct
- `github.com/julienschmidt/httprouter`: A high-performance HTTP request router.
- `github.com/redis/go-redis/v9`: Redis client for Go.
- `go.mongodb.org/mongo-driver`: Official MongoDB driver for Go.
- `github.com/stretchr/testify`: Used for writing and running tests.

## Project Structure

- `constants`: Contains custom error definitions.
- `entity`: Defines the data models.
- `repository`: Implements data access layer (MongoDB and Redis).
- `routes`: Defines HTTP routes and handlers.
- `server`: Manages server connections (MongoDB and Redis).
- `services`: Contains business logic.
- `templates`: HTML templates for the frontend.
- `util`: Utility functions.

## Features

- Create shortened URLs
- List all shortened URLs
- Delete a shortened URL
- Update a shortened URL
- Redirect to original URL using the short code

## Setup and Running

1. Ensure you have Go installed on your system.
2. Clone the repository.
3. Create a `.env` file in the root directory with the following variables:
```
SERVICE_HOST=
SERVICE_PORT=
SERVICE_PROTOCOL=
REDIS_HOST=
REDIS_PORT=
REDIS_PASSWORD=
REDIS_TTL=
MONGODB_USER=
MONGODB_PASSWORD=
MONGODB_HOST=
MONGODB_DATABASE_NAME=
MONGODB_OPTIONS=
```
4. Run `go mod download` to install dependencies.
5. Start the server with `go run main.go`.

The server will start on the host and port specified in your `.env` file.

## API Endpoints

- `GET /`: Home page
- `POST /shorten-url`: Create a new shortened URL
- `GET /shorten-url`: List all shortened URLs
- `DELETE /:shortCode`: Delete a shortened URL
- `PATCH /:shortCode`: Update a shortened URL
- `GET /s/:shortCode`: Redirect to the original URL

## Testing

- Run tests using the `make test` command in the project root directory to run all test.
- Run tests using the `make test/cover` command in the project root directory to run all test with coverage report.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open-source and available under the [MIT License](LICENSE).
