# Step 1: Modules caching
FROM golang:1.24.3-alpine3.21 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.24.3-alpine3.21 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/wordka

# Step 3: Final
FROM scratch

WORKDIR /app

COPY --from=builder /app/.env.example /app/.env
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /bin/app /app/app

CMD ["/app/app"]