FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./go.* ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o search-service ./cmd/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/search-service .

EXPOSE 8086

CMD ["./search-service"]
