
# FROM golang:1.20.6

# WORKDIR /app

# COPY go.mod .
# COPY go.sum .
# COPY . .
# RUN go mod download

# RUN go build -o main

# EXPOSE 8888
# CMD ["./main"]




# Start from the latest golang base image
FROM golang:1.20.6



# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
ARG HOST=shop-postgres
ENV DB_HOST=$HOST
# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8888

# Command to run the executable
# CMD ["./main"] 

CMD ["go", "run", "main.go"]

