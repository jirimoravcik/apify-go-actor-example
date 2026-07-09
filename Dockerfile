# --- Build stage ---
FROM golang:1.23 AS build

WORKDIR /app

# Download dependencies first to leverage layer caching.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a fully static binary so it runs on scratch (no libc, no dynamic linker).
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /goactor

# --- Final stage ---
FROM scratch

# CA certificates are needed for HTTPS calls to the Apify API.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /goactor /goactor

ENTRYPOINT ["/goactor"]
