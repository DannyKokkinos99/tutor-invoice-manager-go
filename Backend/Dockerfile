# Use Golang as the builder image
FROM golang:1.24-alpine AS builder

# Install SQLite dependency
# RUN apk add --no-cache gcc musl-dev

# Install CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:${PATH}"

# Set the working directory inside the container
WORKDIR /app

COPY . .
RUN go mod download
RUN swag init

# Build the Go application
RUN go build -o /bin/main .

# Expose the application port
EXPOSE 4200

# Start the application
CMD ["sh", "-c", "CompileDaemon -log-prefix=false -directory=/app -polling=true -exclude=.*tmp.* -command='go run .'"]

