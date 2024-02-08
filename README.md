# RESTful-API Users And TodoList

Users And TodoList

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)

## Features

- TodoList
  - GetAll with Pagination
  - GetByID
  - Create TodoList
  - Update TodoList
  - Delete TodoList
  - Create Many TodoList
  - Delete Many TodoList
- Users
  - Login & Logout
  - Register and Register Many User For Role ADMIN & MODERATOR
  - Auhtentication & Authorization
## Getting Started

### Prerequisites

- Install External Package
    - GORM & MySQL Driver
    - Validator
    - Gin Web Framework
    - Viper
    - Goose DB Migration

### Installation

1. Run in CMD or Bash
    ```bash 
    go mod tidy
    ```
2. If you want to change the configuration, it's in the .env file
3. Run in CMD or Bash
    ```bash 
    go run cmd/app/main.go
    ```
   OR
    ```bash
    make build
   ./bin/main
    ```
4. Run Migration 
    ```bash 
   goose -dir internal/db/migration postgres "user=postgres dbname=test sslmode=disable" up
   make migration
   make migrationUp 
   make migrationDown
    ```
## Usage

Go to http://localhost:3500
See Result

## API Documentation

At Another Time Im So Lazy 

## Contributing

- Tirta Hakim Pambudhi

## License

This project is licensed under the Tirta Hakim Pambudhi - see the [LICENSE.md](LICENSE.md) file for details.

## Next Features
- Reset Password For Users