#########################################################
# Build Executable binary
#########################################################
FROM golang:1.22.2-alpine AS builder

# Install dependencies ('make')
RUN apk update && apk add --no-cache make

# Set the current working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. They will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the WORKDIR inside the container
COPY . .

# Build the Go app
RUN make all

#########################################################
# Build a small image with a scratch base and the binary
#########################################################
FROM scratch

# RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/miracummapper /app/miracummapper

# Copy files needed for the application
COPY cmd/miracummapper/default-config.yaml /app/default-config.yaml
# COPY db/migrations /app/migrations

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./miracummapper"]

# # alpine enables to use curl for healthcheck
# FROM alpine:latest

# RUN apk --update --no-cache add curl && rm -rf /var/cache/apk/*

# WORKDIR /app

# # Copy the pre-built binary from the previous stage
# COPY --from=builder /app/miracummapper /app/miracummapper

# # Copy files needed for the application
# COPY cmd/miracummapper/default-config.yaml /app/default-config.yaml
# COPY db/migrations /app/migrations

# # Expose port 8080 to the outside world
# EXPOSE 8080

# # Command to run the executable
# CMD ["./miracummapper"]