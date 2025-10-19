FROM golang:1.25.3-alpine AS build
WORKDIR /app

# Copy source and resolve dependencies (ensures go.sum is present)
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

FROM alpine:3.20
RUN adduser -D -g '' appuser \
    && apk add --no-cache ca-certificates
USER appuser
WORKDIR /home/appuser
COPY --from=build /out/server /usr/local/bin/server
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/server"]
