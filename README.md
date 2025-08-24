# Trading Platform Simulator - Go Backend

This repository contains the primary backend service for the Full-Stack Trading Platform Simulator. What started as a simple trading simulator has evolved into a data-driven market analysis tool.

This service, written in Go (Golang), orchestrates a polyglot microservices architecture. It handles secure user authentication, portfolio management, and simulated trading via a REST API. Additionally, it runs a real-time data pipeline to fetch market news from Finnhub, calls a separate Python service for AI sentiment analysis, and broadcasts all data updates to clients via WebSockets.

[Try the Live Application](https://trade-sim-liard.vercel.app/)

## Key Features

- **AI-Powered Sentiment Analysis:** A background worker in Go polls the Finnhub API for market news, then calls a custom Python/FastAPI microservice to perform sentiment analysis using a fine-tuned Hugging Face transformer model.
- **Advanced Data Aggregation:** Continuously calculates and stores a 14-day Exponential Moving Average (EMA) of news sentiment for each stock, providing a dynamic market mood indicator.
- **Secure Authentication:** User registration with **bcrypt** for secure password hashing and a login flow that issues **JSON Web Tokens (JWTs)** for session management.
- **Real-time Data & Notifications:** A concurrent WebSocket server pushes live stock prices, new market news, and trade confirmation events to clients, enabling a fully dynamic UI.
- **Transactional Integrity:** All financial operations (buy/sell orders) are handled within atomic PostgreSQL transactions to ensure data integrity.
- **Microservices Architecture:** Orchestrates communication between the user, the database, and a separate AI microservice for specialized tasks.

## System Architecture

This project uses a microservices approach:
1.  **Next.js Frontend:** The user interface, deployed on Vercel.
2.  **Go Backend (This Repo):** The core service for user logic, trading, and data orchestration. Deployed on Render.
3.  **Python Sentiment Service:** A separate FastAPI service deployed on Hugging Face Spaces that exposes the sentiment analysis model.
4.  **PostgreSQL Database:** The central data store, hosted on Supabase.

## Tech Stack

- **Language:** Go (Golang)
- **API:** REST & WebSockets
- **Key Libraries:**
  - `net/http` (for REST API server)
  - `gorilla/websocket` (for WebSocket communication)
  - `golang-jwt/jwt/v5` (for JWT creation & validation)
  - `golang.org/x/crypto/bcrypt` (for password hashing)
  - `database/sql` with `lib/pq` (for PostgreSQL interaction)
- **Database:** PostgreSQL
- **Environment:** Docker & Docker Compose (for local database setup)

---

## Getting Started

Follow these instructions to get the backend server running on your local machine.

### Prerequisites

- [Go](https://go.dev/dl/) (version 1.22+ recommended)
- [Docker](https://www.docker.com/products/docker-desktop/) and Docker Compose
- The [Next.js frontend repository](https://github.com/ashutoshDhopte/trading_platform_frontend)
- (Optional for local run) The [Python Sentiment Analysis service](https://huggingface.co/spaces/ashutoshDhopte/sentiment_analysis_service) or your own local version.

### Installation & Setup

1.  **Clone the Repository**
    ```bash
    git clone [https://github.com/your-username/trading-platform-backend.git](https://github.com/your-username/trading-platform-backend.git)
    cd trading-platform-backend
    ```

    **Other repositories**

    Frontend
    ```bash
    https://github.com/ashutoshDhopte/trading_platform_frontend
    ```

    Sentiment Analysis
    ```bash
    https://huggingface.co/spaces/ashudhopte123/trading_platform_ml_huggingface/tree/main
    ```

2.  **Set Up the Database with Docker**
    A `docker-compose.yml` file is included to easily spin up a PostgreSQL container.
    ```bash
    # From the project root directory
    docker-compose up -d
    ```
    This will start a PostgreSQL server on `localhost:5432`.

    Update Aug 23 2025: Added new services Kafka, Zookeeper and Debezium Connect for CDC (Change Data Capture) implementation.


3.  **Run Database Schema Setup**
    Connect to the running PostgreSQL instance using a database tool (like DBeaver or `psql`) and run the SQL script found at `/postgres/create.sql` to set up the necessary tables and initial data.


4. **Run the following command to set up Debezium CDC,to capture changes in "orders" table.**

    ```bash
   curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '{
    "name": "trading-transactions-connector",
    "config": {
      "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
      "database.hostname": "db",
      "database.port": "5432",
      "database.user": "trading_user",
      "database.password": "trading_password",
      "database.dbname": "trading_db",
      "database.server.name": "tradingplatform",
      "topic.prefix": "tradingplatform",
      "table.include.list": "public.orders",
      "plugin.name": "pgoutput"
    }
    }'
   ```

4.  **Configure Environment Variables**
    Create a `.env` file in the project's root directory. You will need API keys from [Finnhub.io](https://finnhub.io/) and a self-generated JWT secret (`openssl rand -base64 32`).

    ```env
    # .env file
    DB_USER=trading_user
    DB_PASSWORD=trading_password
    DB_NAME=trading_db
    DB_HOST=localhost
    DB_PORT=5432
    DB_SSL_MODE=disable
    
    JWT_SECRET=your-super-long-randomly-generated-string-goes-here
    
    FINNHUB_API_KEY=your-finnhub-api-key
    SENTIMENT_API_URL=[https://ashutoshdhopte-sentiment-analysis-service.hf.space/sentiment](https://ashutoshdhopte-sentiment-analysis-service.hf.space/sentiment)
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
-   `POST /login`: Authenticates a user and returns a JWT.

### Trading & Data API (Protected - Requires `Authorization: Bearer <JWT>`)

-   `GET /dashboard`: Fetches the user's main dashboard data.
-   `GET /markets/{ticker}`: Fetches detailed market and news analysis for a specific stock ticker.
-   `POST /buy-stocks`: Executes a buy order.
-   `POST /sell-stocks`: Executes a sell order.

### WebSocket API

-   **Endpoint:** `ws://localhost:8080/trade-sim/ws/dashboard`
-   **Functionality:** Establishes a WebSocket connection. Pushes real-time stock price updates, new market news, and trade confirmations to the client.

