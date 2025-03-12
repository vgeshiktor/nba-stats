# Use the official Golang image for building
FROM golang:1.23 AS builder

WORKDIR /app

# Copy Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o nba-stats cmd/server/main.go

# Use a lightweight image for the final runtime
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/nba-stats .

# Copy the migrations folder into the image
COPY migrations/ migrations/

# Install SQLite (if needed)
RUN apk add --no-cache sqlite

# TODO: debug only, remove later
RUN apk add --no-cache bash curl vim strace net-tools procps
RUN apk add --no-cache gdb delve

# Expose the port and run the application
EXPOSE 8080
CMD ["./nba-stats"]
