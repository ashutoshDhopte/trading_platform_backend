# Trading Simulator - Go Backend

This repository contains the backend service for the Full-Stack Trading Platform Simulator. It's a REST API built with Go (Golang) that handles all core business logic, including user portfolio management, processing simulated buy/sell orders, and broadcasting real-time stock price updates via WebSockets.

## Key Features

- **REST API:** Provides endpoints for buying/selling mock stocks, and fetching portfolio, stock list, and transaction history.
- **Real-time Price Updates:** A concurrent WebSocket server pushes simulated live stock price updates to all connected clients.
- **Transactional Integrity:** All financial operations (buy/sell orders) are handled within atomic PostgreSQL transactions to ensure data integrity.
- **Database Management:** Connects to a PostgreSQL database to manage all persistent data like user balances, holdings, and transaction logs.
- **Environment-based Configuration:** Securely manages database connections and other settings via environment variables.

## Tech Stack

- **Language:** Go (Golang)
- **API:** REST & WebSockets
- **Primary Libraries:**
  - `net/http` (for REST API server)
  - `gorilla/websocket` (for WebSocket communication)
  - `database/sql` with `lib/pq` (for PostgreSQL interaction)
- **Database:** PostgreSQL
- **Environment:** Docker & Docker Compose (for local database setup)

---

## Getting Started

Follow these instructions to get the backend server running on your local machine for development and testing.

### Prerequisites

- [Go](https://go.dev/dl/) (version 1.22+ recommended)
- [Docker](https://www.docker.com/products/docker-desktop/) and Docker Compose
- A running PostgreSQL instance (instructions to run one with Docker are below)

### Installation & Setup

1.  **Clone the Repository**
    ```bash
    # If the backend is in its own repository
    git clone [https://github.com/your-username/trading-platform-backend.git](https://github.com/your-username/trading-platform-backend.git)
    cd trading-platform-backend
    ```

2.  **Set Up the Database with Docker**
    A `docker-compose.yml` file is included in the project root to easily spin up a PostgreSQL container.
    ```bash
    # From the project root directory
    docker-compose up -d
    ```
    This will start a PostgreSQL server on `localhost:5432`.

3.  **Run Database Migrations/Schema Setup**
    Connect to the running PostgreSQL instance using a database tool (like DBeaver, pgAdmin, or `psql`) and run the SQL script `/postgres/create.sql` to set up the necessary tables and initial data.

4.  **Configure Environment Variables**
    The application uses environment variables for configuration. You can use a `.env` file locally with a library like `godotenv` or set them in your terminal.
    ```env
    # .env file
    DB_USER=trading_user
    DB_PASSWORD=trading_password
    DB_NAME=trading_db
    DB_HOST=localhost
    DB_PORT=5432
    DB_SSL_MODE=disable
    ```

5.  **Run the Backend Server**
    ```bash
    # Navigate to the directory (replace <dir>) containing main.go
    cd <dir>
    
    # Run the application
    go run main.go
    ```
    The server should now be running on `http://localhost:8080`.

---

## API Endpoints

**Base URL:** `http://localhost:8080` (or as configured)
*(Note: If your router uses a prefix like `/trade-sim`, adjust the paths accordingly, e.g., `http://localhost:8080/trade-sim/api/stocks`)*

### REST API
  Swagger coming soon


### WebSocket API

-   **Endpoint:** `ws://localhost:8080/trade-sim/ws/dashboard`
-   **Functionality:** Establishes a WebSocket connection. The server will automatically start pushing real-time stock price updates to the client every few seconds.

