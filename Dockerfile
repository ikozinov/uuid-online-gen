# Use a Go base image
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the WebAssembly binary
# We build it into web/app.wasm so the server can find it,
# and the static generator can copy it to dist/
RUN GOOS=js GOARCH=wasm go build -o web/app.wasm

# Generate the static site (creates dist/ directory)
RUN go run main.go dist

# Build the server binary
RUN go build -o server

# Expose the port
EXPOSE 8000

# Run the server
CMD ["./server"]
