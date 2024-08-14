# # Start from the official Go image
# FROM golang:1.22

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy go mod and sum files
# COPY go.mod go.sum ./

# # Download all dependencies, tidy up, and print go.mod for debugging
# RUN go mod download && go mod tidy && cat go.mod

# # Install Air
# RUN go install github.com/air-verse/air@latest
# # RUN go install github.com/cosmtrek/air@latest
# # Copy the source from the current directory to the Working Directory inside the container
# COPY . .

# # Run go mod tidy again after copying all files and print go.mod
# RUN go mod tidy && cat go.mod

# ARG HOST=shop-postgres
# ENV DB_HOST=$HOST

# # Expose port 8888 to the outside world
# EXPOSE 8888

# # Command to run Air for live reloading
# CMD ["air", "-c", ".air.toml"]



FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Run go mod tidy to update go.mod and go.sum
RUN go mod tidy

ARG HOST=shop-postgres
ENV DB_HOST=$HOST

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8888

# Command to run the executable
CMD ["go", "run", "cmd/server/main.go"]