FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o automated-messageing-system /app/main.go

FROM golang:1.25.1-alpine3.22 AS worker

WORKDIR /app

COPY --from=builder /app/automated-messageing-system .

ENTRYPOINT [ "./automated-messageing-system" ]
