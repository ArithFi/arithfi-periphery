FROM golang:1.21.5 as builder

RUN apt-get update && apt-get install -y --no-install-recommends tzdata

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o bin/rest-api/main cmd/rest-api/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/transferLogsRetrieve cmd/workers/transferLogsRetrieve.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/transferLogsAnalysis cmd/workers/transferLogsAnalysis.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/tokenholderSnapshotCreate cmd/workers/tokenholderSnapshotCreate.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/pancakeswapSnapshotCreate cmd/workers/pancakeswapSnapshotCreate.go && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin/workers/getForexKline cmd/workers/getForexKline.go

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin .

ENV TZ=Asia/Shanghai
ENV PORT=8080

EXPOSE 8080

CMD ["./rest-api/main"]