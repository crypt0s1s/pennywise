# Terrible Ideas 2025 - Go Game Server

This repository contains the backend Go server for a collection of simple games, developed as part of the 'Terrible Ideas 2025' hackathon. The server handles game logic, state management, and provides API endpoints for game interaction.

## Project Structure

- `main.go`: The main entry point of the server, setting up routes, middleware, and handling WebSocket connections.
- `models/`: Contains Go struct definitions for shared data models like `Battle` and `Product`.
- `golf/`: Implements the Golf game logic, including game state, moves, and API handlers.
- `scissorsPaperRock/`: Implements the Scissors, Paper, Rock game logic, including game state, moves, and API handlers.
- `vendor/`: Go module dependencies.

## Features

- **Multi-game Support**: Currently supports Golf and Scissors, Paper, Rock.
- **WebSocket Communication**: Real-time updates for game states.
- **CORS Enabled**: Configured to allow requests from `http://localhost:3000` and `https://www.amazon.com.au`.
- **Game State Management**: Manages game instances and their states in memory.
- **Recent Games Endpoint**: `/games` endpoint to list active games created within the last two minutes.

## Setup Instructions

To set up and run the project locally, follow these steps:

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/terrible-ideas-2025.git
    cd terrible-ideas-2025
    ```

2.  **Install Go dependencies:**
    ```bash
    go mod tidy
    ```

## How to Run

To start the Go server, run the following command from the root directory of the project:

```bash
go run main.go
```

The server will typically start on `http://localhost:8080` (or the port configured in `main.go`).

## API Endpoints

### General

- `GET /games`: Returns a list of active game IDs and their types created within the last two minutes.

### Scissors, Paper, Rock (SPR)

- `POST /spr/initiate`: Initiates a new SPR game.
- `POST /spr/accept`: Accepts an invitation to an SPR game.
- `POST /spr/:game_id`: Logs a move for a specific SPR game.
- `GET /spr/:game_id`: Retrieves the state of a specific SPR game.

### Golf

- `POST /golf/initiate`: Initiates a new Golf game.
- `POST /golf/accept`: Accepts an invitation to a Golf game.
- `POST /golf/:game_id`: Logs a move for a specific Golf game.
- `GET /golf/:game_id`: Retrieves the state of a specific Golf game.

## Contributing

Why would you want to build on this crap 😂
