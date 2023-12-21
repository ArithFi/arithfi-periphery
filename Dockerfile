FROM golang:1.21.5 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o bin/server/main cmd/server/main.go

FROM scratch

COPY --from=builder /app/bin/server/main .

ENV PORT=8080

EXPOSE 8080

CMD ["./main"]