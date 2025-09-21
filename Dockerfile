FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/main ./cmd/server/main.go

FROM alpine:3.22 AS production

WORKDIR /app

COPY --from=builder /build/main /app/main

EXPOSE 5000

CMD [ "./main" ]

FROM builder AS development

WORKDIR /app

COPY --from=builder /app /app

RUN go install github.com/air-verse/air@v1.52.3

EXPOSE 5000

CMD ["air", "-c", ".air.toml"]
