# URL Shortener

This is a simple URL shortener API built using Go and Gin framework. It uses GORM as the ORM and a SQLite database to store the shortened URLs.

## Running the API using Docker

### Prerequisites
- Docker

### Steps

1. Clone the repository and navigate to the project directory.
```bash
git clone https://github.com/C-m3-Codin/url_shortner.git
cd repo
```

2. Build the Docker image using the provided Dockerfile.
```bash
docker build -t url-shortener .
```

3. Run the Docker container using the provided docker-compose.yml file.
```bash
docker-compose up -d
```

## Usage

The API provides two endpoints:

1. `POST /create` - Creates a new shortened URL.
2. `GET /:shortenedUrl` - Redirects to the original URL associated with the given shortened URL.

### POST /create

To create a new shortened URL, send a JSON POST request to the `/create` endpoint with the following format:
```json
{
  "original_url": "https://www.google.com"
}
```

### GET /:shortenedUrl

To redirect to the original URL associated with a shortened URL, send a GET request to the `/:shortenedUrl` endpoint with the shortened URL as a parameter. For example, if the shortened URL is `abc123`, the request URL would be `http://localhost:8000/abc123`.

## Code Overview

The main Go code is contained in `main.go`. It creates a database connection and defines the endpoints for the API.

The `services` package contains the code for creating the database connection.

The `ShortLink` struct represents a shortened URL and is defined in `main.go`. The `generateShortenedURL` function generates a random 8-character string to use as the shortened URL.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.