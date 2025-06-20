# Trading Platform Simulator - Go Backend

This repository contains the backend service for the Full-Stack Trading Platform Simulator. It's a REST API built with Go (Golang) that handles all core business logic, including secure user authentication, portfolio management, processing simulated buy/sell orders, and broadcasting real-time stock price updates via WebSockets.

[Try it](https://trade-sim-liard.vercel.app/)

## Key Features

- **Secure Authentication:** User registration with **bcrypt** for secure password hashing and a login flow that issues **JSON Web Tokens (JWTs)** for session management.
- **Protected REST API:** Provides endpoints for buying/selling mock stocks and fetching user-specific data. Sensitive endpoints are protected using JWT authentication middleware.
- **Real-time Price Updates & Notifications:** A concurrent WebSocket server pushes simulated live stock price updates and trade confirmation events to clients, enabling real-time UI updates and browser notifications.
- **Transactional Integrity:** All financial operations (buy/sell orders) are handled within atomic PostgreSQL transactions to ensure data integrity.
- **Database Management:** Connects to a PostgreSQL database to manage all persistent data like user credentials, balances, holdings, and transaction logs.
- **Environment-based Configuration:** Securely manages database connections and secrets via environment variables.

## Tech Stack

- **Language:** Go (Golang)
- **API:** REST & WebSockets
- **Primary Libraries:**
  - `net/http` (for REST API server)
  - `gorilla/websocket` (for WebSocket communication)
  - `golang-jwt/jwt/v5` (for JWT creation & validation)
  - `golang.org/x/crypto/bcrypt` (for password hashing)
  - `database/sql` with `lib/pq` (for PostgreSQL interaction)
- **Database:** PostgreSQL
- **Environment:** Docker & Docker Compose (for local database setup)

---

## Getting Started

Follow these instructions to get the backend server running on your local machine for development and testing.

### Prerequisites

- [Go](https://go.dev/dl/) (version 1.22+ recommended)
- [Docker](https://www.docker.com/products/docker-desktop/) and Docker Compose
- The [Next.js frontend repository](https://github.com/ashutoshDhopte/trading_platform_frontend) for a complete end-to-end experience.

### Installation & Setup

1.  **Clone the Repository**
    ```bash
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

3.  **Run Database Schema Setup**
    Connect to the running PostgreSQL instance using a database tool (like DBeaver, pgAdmin, or `psql`) and run the SQL script found at `/postgres/create.sql` to set up the necessary tables and initial data.

4.  **Configure Environment Variables**
    The application uses environment variables for configuration. Create a `.env` file in the project's root directory. You can generate a strong JWT secret by running `openssl rand -base64 32` in your terminal.

    ```env
    # .env file
    DB_USER=trading_user
    DB_PASSWORD=trading_password
    DB_NAME=trading_db
    DB_HOST=localhost
    DB_PORT=5432
    DB_SSL_MODE=disable
    JWT_SECRET=your-super-long-randomly-generated-string-goes-here
    ```

5.  **Run the Backend Server**
    ```bash
    # Navigate to the directory containing main.go
    cd cmd/server
    
    # Run the application
    go run main.go
    ```
    The server should now be running on `http://localhost:8080`.

---

## API Endpoints

**Base URL:** `http://localhost:8080/trade-sim`

### Authentication API (Public)

-   `POST /create-account`: Creates a new user account.
    -   **Body:** `{ "email": "user@example.com", "password": "secure_password" }`
-   `POST /login`: Authenticates a user and returns a JWT.
    -   **Body:** `{ "email": "user@example.com", "password": "secure_password" }`

### Trading & Data API (Protected - Requires `Authorization: Bearer <JWT>`)

-   `GET /dashboard`: Fetches the user's dashboard data, including portfolio.
-   `POST /buy-stocks`: Executes a buy order for the authenticated user.
    -   **Body:** `{ "stock_ticker": "FAKE_AAPL", "quantity": 10 }`
-   `POST /sell-stocks`: Executes a sell order for the authenticated user.
    -   **Body:** `{ "stock_ticker": "FAKE_AAPL", "quantity": 5 }`

### WebSocket API

-   **Endpoint:** `ws://localhost:8080/trade-sim/ws/dashboard`
-   **Functionality:** Establishes a WebSocket connection. The server will automatically start pushing real-time stock price updates to the client every few seconds. Trade confirmation events may also be sent over this channel.

