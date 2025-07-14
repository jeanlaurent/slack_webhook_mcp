# Slack Webhook MCP Server

A Model Context Protocol (MCP) server that enables sending messages to Slack through webhooks. This server provides a single tool `send_message` for sending formatted messages to Slack channels.

## Features

- Send messages to Slack via webhook URLs
- Support for channel override, custom usernames, and icons
- Markdown formatting support
- Comprehensive error handling
- Docker support with multi-stage builds

## Prerequisites

- Go 1.23 or later
- A Slack webhook URL (see [Slack's documentation](https://api.slack.com/messaging/webhooks) for setup)
- Docker (optional, for containerized deployment)

## Installation

### From Source

```bash
git clone <repository-url>
cd mcp-slack-webhook
go mod download
go build -o mcp-slack-webhook .
```

### Using Docker

```bash
docker build -t mcp-slack-webhook .
docker run -e SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL" mcp-slack-webhook
```

## Usage

### Running the Server

First, set the required environment variable:

```bash
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
```

Then run the server:

```bash
./mcp-slack-webhook
```

The server will start and listen for MCP protocol messages on stdin/stdout.

### Tool: send_message

The server provides a single tool called `send_message` with the following parameters:

#### Required Parameters

- `text` (string): The message text to send

#### Required Environment Variable

- `SLACK_WEBHOOK_URL`: The Slack webhook URL to send messages to

#### Optional Parameters

- `channel` (string): The channel to send the message to (overrides webhook default)
- `username` (string): Username to display as the sender
- `icon_emoji` (string): Emoji to use as the icon (e.g., `:robot_face:`)
- `icon_url` (string): URL to an image to use as the icon
- `markdown` (boolean): Whether to enable markdown formatting (default: true)

#### Example Tool Call

First, set the environment variable:

```bash
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
```

Then call the tool:

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/call",
  "params": {
    "name": "send_message",
    "arguments": {
      "text": "Hello from MCP! This is a *test* message.",
      "channel": "#general",
      "username": "MCP Bot",
      "icon_emoji": ":robot_face:",
      "markdown": true
    }
  }
}
```

## Configuration

### Slack Webhook Setup

1. Go to [Slack API: Incoming Webhooks](https://api.slack.com/messaging/webhooks)
2. Create a new webhook for your workspace
3. Copy the webhook URL for use with this server

### Environment Variables

#### Required

- `SLACK_WEBHOOK_URL`: The Slack webhook URL for sending messages

#### Optional

You may also want to configure:

- Logging levels
- Authentication (not implemented in this basic version)

## Docker Deployment

The included Dockerfile uses a multi-stage build for optimal image size:

```bash
# Build the image
docker build -t mcp-slack-webhook .

# Run the container
docker run -i -e SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL" mcp-slack-webhook
```

For production use, you might want to:

```bash
# Run with environment file
docker run -i --env-file .env mcp-slack-webhook

# Or run with secrets management
docker run -i -e SLACK_WEBHOOK_URL="$(cat /path/to/webhook_url_secret)" mcp-slack-webhook
```

## Development

### Project Structure

```
mcp-slack-webhook/
├── main.go           # Main application code
├── go.mod            # Go module dependencies
├── go.sum            # Go dependency checksums
├── Dockerfile        # Multi-stage Docker build
└── README.md         # This file
```

### Dependencies

- [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) v0.32.0 - MCP protocol implementation

### Building

```bash
# Regular build
go build -o mcp-slack-webhook .

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o mcp-slack-webhook-linux .
```

### Testing

```bash
# Run tests (when available)
go test ./...

# Test Docker build
docker build -t mcp-slack-webhook-test .
```

## Slack Message Format

Messages sent to Slack follow the webhook format:

```json
{
  "text": "Your message text",
  "channel": "#channel-name",
  "username": "Bot Name",
  "icon_emoji": ":robot_face:",
  "icon_url": "https://example.com/icon.png",
  "mrkdwn": true
}
```

## Error Handling

The server handles various error conditions:

- Invalid webhook URLs
- Missing required parameters
- Slack API errors
- Network connectivity issues

All errors are returned as MCP tool result errors with descriptive messages.

## Security Considerations

- Webhook URLs are securely stored as environment variables and not exposed in tool calls
- Webhook URLs should never be logged or exposed in debug output
- Consider implementing authentication for production use
- Validate input parameters to prevent injection attacks
- Use HTTPS webhook URLs only
- In production, use secure secret management systems instead of plain environment variables

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
