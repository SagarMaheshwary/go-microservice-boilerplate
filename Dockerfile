FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/server/main.go

FROM alpine:3.22 AS production

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 5000
CMD ["./main"]

FROM golang:1.25 AS development

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@v1.52.3
RUN go mod download

EXPOSE 5000
CMD ["air", "-c", ".air.toml"]
