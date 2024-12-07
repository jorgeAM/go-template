# Go Template

A boilerplate template for initializing a Golang API with essential features. This template includes structured project organization, basic configurations, logging, routing, and database integration to help you quickly set up and start building your Golang-based APIs.

## Features

- **Structured Project Layout**: Organized folder structure for scalability and maintainability.
- **Configuration Management**: Easy configuration with `viper`.
- **Logging**: Integrated Zap logger for structured logging.
- **Routing**: HTTP routing with Chi.
- **Database Integration**: PostgreSQL integration using `sqlx` and `goqu`.
- **User CRUD**: Basic user CRUD operations.

## Getting Started

1. **Clone the repository:**

   ```sh
   git clone https://github.com/jorgeAM/go-template.git
   cd go-template
   ```
2. **Install dependencies:**
    ```sh
   go mod tidy
   ```
3. **Set up environment variables:**
Create a .env file in the root directory with the following variables (you can copy .env.example):

    ```env
    PORT=8080

    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_DB=coderhouse
    POSTGRES_USER=admin
    POSTGRES_PASSWORD=passwd123
    POSTGRES_MAX_IDLE_CONNECTIONS=10
    POSTGRES_MAX_OPEN_CONNECTIONS=30
    ```

4. **Run the application::**
    just run the following code:

    ```sh
    make run
    ```

    you can also run:
    ```sh
    go run cmd/go-template/main.go | jq '.'
    ```

## Usage

- Access the API at `http://localhost:8080`.
- Basic user CRUD endpoints are available under `/users`.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/jorgeAM/go-template/blob/main/LICENCE) file for details.


