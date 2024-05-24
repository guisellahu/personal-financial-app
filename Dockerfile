# Use the official Golang image as a base
FROM golang:1.22.0

# Install dockerize
RUN apt-get update && apt-get install -y wget && \
    wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz && \
    tar -C /usr/local/bin -xzvf dockerize-linux-amd64-v0.6.1.tar.gz && \
    rm dockerize-linux-amd64-v0.6.1.tar.gz

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .



# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable, waiting for PostgreSQL to be ready
CMD dockerize -wait tcp://db:5432 -timeout 60s ./main
