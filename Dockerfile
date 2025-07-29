FROM golang:1.23-alpine AS gobuilder
WORKDIR /app
RUN --mount=type=cache,target=/var/cache/apk apk add --no-cache git
COPY --link go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY --link . .
RUN go build -o mcp-slack-webhook .

FROM alpine:latest
RUN --mount=type=cache,target=/var/cache/apk apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=gobuilder /app/mcp-slack-webhook .
CMD ["./mcp-slack-webhook"]