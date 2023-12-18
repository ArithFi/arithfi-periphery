FROM golang:1.21.5 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o bin/echo-server cmd/server/server.go

FROM scratch

COPY --from=builder /app/bin/echo-server .

ENV PORT=8080

EXPOSE 8080

CMD ["./bin/echo-server"]