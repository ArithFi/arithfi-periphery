FROM golang:1.21.5 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o bin/rest-api/main cmd/rest-api/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/transfer_logs_retrieve cmd/workers/transfer_logs_retrieve.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/transfer_logs_analysis cmd/workers/transfer_logs_analysis.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/tokenholder_snapshot_create cmd/workers/tokenholder_snapshot_create.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/pancakeswap_snapshot_create cmd/workers/pancakeswap_snapshot_create.go

FROM scratch

COPY --from=builder /app/bin .

ENV PORT=8080

EXPOSE 8080

CMD ["./rest-api/main"]