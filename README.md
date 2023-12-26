# HTTP-Server

Developed an HTTP server with the ability to interact with the database and perform CRUD (Create, Read, Update, Delete)
operations. The server is designed to handle incoming HTTP requests and perform CRUD operations. Integrated Kafka into
project, while during CRUD operations it records events using methods in Kafka and saves the time of occurrence, type and raw
request and displays them on the screen. Covered methods with integration and unit tests;
## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [Get](#get)
  - [Create](#create)
  - [Delete](#delete)
  - [Update](#update)
- [Testing](#testing)
- [Example](#example)

## Features

- Get: Retrieve data from the database based on the provided ID.
- Create: Add new data to the database.
- Delete: Remove data from the database based on the provided ID.
- Update: Update existing data in the database based on the provided ID.

## Prerequisites

Before running this application, ensure that you have the following prerequisites installed:

- Go: [Install Go](https://go.dev/doc/install/)
- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. Clone the repository:
  ```bash
    git clone https://github.com/NRKA/HTTP-Server.git
  ```

2. Navigate to the project directory:
  ```
    cd HTTP-Server
  ```
3. Build the Docker image:
  ```
    docker-compose build
  ```

## Usage
1. Start the Docker containers:
  ```
    docker-compose up
  ```
2. The application will be accessible at:
  ```
    localhost:9000
  ```
## API Endpoints

### Get

- Method: GET
- Endpoint: /entity
- Query Parameter: id=[entity_id]

Retrieve data from the database based on the provided ID.

- Response:
  - 200 OK: Returns the data in the response body.
  - 400 Bad Request: If the `id` query parameter is missing.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

### Create

- Method: POST
- Endpoint: /entity
- Request Body: JSON payload containing the ID and data.

Add new data to the database.

- Response:
  - 200 OK: If the request is successful.
  - 500 Internal Server Error: If there is an internal server error.

### Delete

- Method: DELETE
- Endpoint: /entity
- Query Parameter: id=[entity_id]

Remove data from the database based on the provided ID.

- Response:
  - 200 OK: If the request is successful.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

### Update

- Method: PUT
- Endpoint: /entity
- Request Body: JSON payload containing the ID and updated data.

Update existing data in the database based on the provided ID.

- Response:
  - 200 OK: If the request is successful.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

## Testing

- To run unit tests, use the following command:
  ```bash
    go test ./... -cover
  ```
- To run integration tests, use the following command:
  ```bash
    go test -tags integration ./... -cover
  ```

## Example 
1) **POST**

- **curl -X POST localhost:9000/article -d '{"name":"test","rating":15}' -i**
```
HTTP/1.1 200 OK

Date: Sun, 08 Oct 2023 15:40:20 GMT

Content-Length: 70

Content-Type: text/plain; charset=utf-8

{"ID":16,"Name":"test","Rating":15,"CreatedAt":"0001-01-01T00:00:00Z"}
```

3) **GET**

- **curl -X GET localhost:9000/article/1 -i**
```
HTTP/1.1 200 OK

Date: Sun, 08 Oct 2023 15:41:35 GMT

Content-Length: 85

Content-Type: text/plain; charset=utf-8

{"ID":1,"Name":"testname","Rating":44,"CreatedAt":"2023-10-05T18:12:38.011325+05:00"}
```


- **curl -X GET localhost:9000/article/10000 -i**
```
HTTP/1.1 404 Not Found

Date: Sun, 08 Oct 2023 15:42:34 GMT

Content-Length: 59

Content-Type: text/plain; charset=utf-8

Failed to find article:article not found
```
3) **DELETE**
- **curl -X DELETE localhost:9000/article/14 -i** ----->14 ID has been removed, please replace it with an existing ID
```
HTTP/1.1 200 OK

Date: Sun, 08 Oct 2023 15:43:22 GMT

Content-Length: 25

Content-Type: text/plain; charset=utf-8
```

- **curl -X DELETE localhost:9000/article/100000 -i**
```
HTTP/1.1 404 Not Found

Date: Sun, 08 Oct 2023 15:46:18 GMT

Content-Length: 45

Content-Type: text/plain; charset=utf-8

Failed to find article:article not found
```
4) **PUT**
- **curl -X PUT localhost:9000/article -d '{"id":1,"name":"testname","rating":44}' -i**
```
HTTP/1.1 200 OK

Date: Sun, 08 Oct 2023 15:46:54 GMT

Content-Length: 25

Content-Type: text/plain; charset=utf-8
```
- **curl -X PUT localhost:9000/article -d '{"id":100000,"name":"testname","rating":44}' -i**
```
HTTP/1.1 404 Not Found

Date: Sun, 08 Oct 2023 11:59:18 GMT

Content-Length: 45

Content-Type: text/plain; charset=utf-8

Failed to find article:article not found
```
