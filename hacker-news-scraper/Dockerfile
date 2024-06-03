FROM golang:1.22 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download

WORKDIR /app/cmd/webserver
# Copy local code to the container image.
COPY .  ./

# Build the binary.
# -mod=readonly switch is used to force failure of the build if the go.mod file needs changes to succeed.
# see https://golang.org/cmd/go/#hdr-Preliminary_module_support for more details.
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -v -o hacker-news-scraper

# Use the Google Distroless image for a minimal container.
FROM gcr.io/distroless/static

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/cmd/webserver/hacker-news-scraper /hacker-news-scraper

# Run the web service on container startup.
CMD ["/hacker-news-scraper"]