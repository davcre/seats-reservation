# Seats reservation
Seats reservation is a RESTful API server written in Golang for performing CRUD operations on reservations.

## Installation
Follow the instructions on [this](https://golang.org/doc/install) site to install go.

Initialize the Go modules with your GitHub repository address:

```bash
go mod init github.com/<your GitHub username>/<project name>
```

Fetch the Go modules:

```bash
go get github.com/gorilla/mux
```

```bash
go get github.com/jinzhu/gorm
```

```bash
go get github.com/go-sql-driver/mysql
```

```bash
go get github.com/joho/godotenv
```

Login to MySQL shell and create the database:

```mysql
CREATE DATABASE Dbname;
```

Create the .env file to pass the DB variables to the server.

Compile the server:

```bash
go build
```

## Usage

Start the server with the following command:

```bash
go run main.go
```

Server works on TCP port 1234.