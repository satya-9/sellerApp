# Use an official Golang runtime as a parent image
FROM golang:1.17.2-alpine3.14

# Set the current working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Install any needed dependencies
RUN go mod download

# Build the Go app
RUN go build -o main .

# Expose port 5000
EXPOSE 5000

# Run the Go app
CMD ["./main"]