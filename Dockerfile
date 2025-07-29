# Build stage
FROM golang:1.23-alpine AS gobuilder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mcp-slack-webhook .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=gobuilder /app/mcp-slack-webhook .
CMD ["./mcp-slack-webhook"] 