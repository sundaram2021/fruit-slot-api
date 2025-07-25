# Fruit Slot API

A simple Golang API that simulates a fruit slot machine. It randomly selects 3 fruits from a pool of 6. If 2 or more fruits are the same, the player wins.

## Project Purpose

This project implements a basic random-outcome API using Go and the Gin web framework. It demonstrates core backend concepts like handling HTTP requests, generating secure random numbers, implementing game logic, returning JSON responses, and writing tests.

## How to Run

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/sundaram2021/fruit-slot-api.git
    cd fruit-slot-api
    ```
2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
3.  **Run directly (for development):**
    ```bash
    go run main.go
    ```
4.  **Build the executable:**
    ```bash
    go build -o play .
    ```
5.  **Run the executable:**
    ```bash
    ./play
    ```
    The server will start on `http://localhost:8080` by default. You can change the port using the `PORT` environment variable (e.g., `PORT=3000 ./play`).

## API Endpoints

### Single Play

*   **URL:** `/play`
*   **Method:** `GET`
*   **Success Response:**
    *   **Code:** `200 OK`
    *   **Content:**
        ```json
        {
          "fruits": ["Lemon", "Lemon", "Grape"],
          "message": "You win!"
        }
        ```
        OR
        ```json
        {
          "fruits": ["Cherry", "Orange", "Watermelon"],
          "message": "Try again!"
        }
        ```

### Ten Plays

*   **URL:** `/play/10`
*   **Method:** `GET`
*   **Success Response:**
    *   **Code:** `200 OK`
    *   **Content:**
        ```json
        {
          "spins": [
            { "fruits": ["Cherry", "Lemon", "Orange"], "message": "Try again!" },
            { "fruits": ["Grape", "Grape", "Watermelon"], "message": "You win!" },
            // ... 8 more spins
          ],
          "win_count": 3
        }
        ```

## Running Tests


**Instructions to Run:**

1.  Create the directory structure and files as shown above.
2.  Initialize the module and install Gin: `go mod init fruit-slot-api && go get github.com/gin-gonic/gin github.com/stretchr/testify/assert`
3.  Run the application: `go run main.go`
4.  Build the executable: `go build -o play .`
5.  Run the executable: `./play`
6.  Test the endpoints using `curl` or Postman:
    *   `curl http://localhost:8080/play`
    *   `curl http://localhost:8080/play/10`
7.  Run tests: `go test ./...`
8.  Run coverage profile 
    - `go test ./... test -cover -coverageprofile=coverage.out`
    - `go tool cover -html=coverage.out` it will route you to html page where you can see coverage profile
9. Code Coverage % is ~81.8%

