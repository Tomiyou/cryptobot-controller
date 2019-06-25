############################
# STEP 1 build executable binary
############################
FROM golang:alpine as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR $GOPATH/bitbucket.org/tomihrib/crypto-arbitrage/
COPY . .

# Fetch dependencies.
# Using go mod with go 1.11
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/crypto-arbitrage

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Use an unprivileged user.
USER appuser

# Set working directory to the home folder
WORKDIR /home/appuser/

# Copy our static executable
COPY --from=builder /go/bin/crypto-arbitrage .

# Run the hello binary.
ENTRYPOINT ["./crypto-arbitrage"]