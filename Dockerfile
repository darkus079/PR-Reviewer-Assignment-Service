FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM golang:1.25-alpine AS runtime

RUN apk --no-cache add ca-certificates make

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 8080

CMD ["./main"]
