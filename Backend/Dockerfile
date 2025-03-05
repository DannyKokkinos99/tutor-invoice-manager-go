# Use Golang as the builder image
FROM golang:1.24-alpine AS builder

# Install SQLite dependency
# RUN apk add --no-cache gcc musl-dev

# Install CompileDaemon
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:${PATH}"
ENV PROJECT_TYPE=development
# Set the working directory inside the container
WORKDIR /app

COPY . .
RUN go mod download
RUN swag init

# Expose the application port
EXPOSE 4200

# Start the application
CMD ["air", "-c", ".air.toml"]

