# Start from the official Go image, version 1.16
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the Go module dependency file to the container
COPY go.mod ./

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Copy the rest of the application files to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 3001 for the container to listen on
EXPOSE 3001

# Start the Go application
CMD ["./main"]