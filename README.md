# Arithfi periphery

## (Production) How to install?

1. Docker, (Recommend)
    ```shell
    docker-compose up
    ```
2. Run by bin file
    ```shell
   go mod download
   CGO_ENABLED=0 GOOS=linux go build -v -o bin/main cmd/main.go
   ./bin/main
   ```

## (Development) How to start?

    ```shell
    go run cmd/main.go
    ```
