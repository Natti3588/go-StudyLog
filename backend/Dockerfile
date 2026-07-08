# ----- build stage -----
FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /studylog-api ./cmd/api

# ----- run stage -----
FROM alpine:latest

WORKDIR /root

COPY --from=builder /studylog-api .

EXPOSE 8080

CMD ["./studylog-api"]
