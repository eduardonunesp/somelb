# Load Balancer

This project implements a load balancer in Go. It includes functionality for reading configuration, managing backend hosts, and handling HTTP requests.

## Installation

To install the dependencies, run:

```sh
go mod tidy
```

## Usage

To run the load balancer, use:

```sh
go run main.go
```

## Configuration

The configuration is provided in YAML format. Below is an example configuration:

```yaml
servers:
  - frontend:
      port: 8080
    backend:
      hosts:
        - url: "http://localhost:8081"
```

### Configuration Structure

- `servers`: List of server configurations.
    - `frontend`: Configuration for the frontend.
        - `port`: Port number for the frontend.
    - `backend`: Configuration for the backend.
        - `hosts`: List of backend hosts.
            - `url`: URL of the backend host.

## Testing

To run the tests, use:

```sh
go test ./...
```

## Manual test

Start two nginx instances using the ports 8080 and 8081.

```sh
docker run -d -p 8081:80 nginx
docker run -d -p 8080:80 nginx
```

And the following configuration

```yaml
servers:
  - backend:
      hosts:
        - url: "http://localhost:8080"
        - url: "http://localhost:8081"
    frontend:
      port: 1337
```

Build and start project

```yaml
go build && ./somelb
```

After the logs showing healthy instances 

```log
2024/09/02 15:36:38 INFO starting server
2024/09/02 15:36:48 INFO Host given signs of alive server.backend.host_1.url=http://localhost:8081
2024/09/02 15:36:48 INFO Host given signs of alive server.backend.host_0.url=http://localhost:8080
2024/09/02 15:36:58 INFO Host given signs of alive server.backend.host_1.url=http://localhost:8081
2024/09/02 15:36:58 INFO Host given signs of alive server.backend.host_0.url=http://localhost:8080
2024/09/02 15:37:08 INFO Host is alive server.backend.host_1.url=http://localhost:8081
2024/09/02 15:37:08 INFO Host given signs of alive server.backend.host_1.url=http://localhost:8081
2024/09/02 15:37:08 INFO Host is alive server.backend.host_0.url=http://localhost:8080
2024/09/02 15:37:08 INFO Host given signs of alive server.backend.host_0.url=http://localhost:8080
```

And start a curl request

```log
curl -vv localhost:1337
```



















