package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// SlackMessage represents the structure for a Slack webhook message
type SlackMessage struct {
	Text      string `json:"text"`
	Channel   string `json:"channel,omitempty"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	IconURL   string `json:"icon_url,omitempty"`
	Markdown  bool   `json:"mrkdwn,omitempty"`
}

// SlackMCPServer represents our MCP server for Slack webhooks
type SlackMCPServer struct {
	server *server.MCPServer
}

// NewSlackMCPServer creates a new Slack MCP server instance
func NewSlackMCPServer() *SlackMCPServer {
	s := &SlackMCPServer{}

	mcpServer := server.NewMCPServer(
		"slack-webhook-server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Create the send_message tool
	tool := mcp.NewTool("send_message",
		mcp.WithDescription("Send a message to Slack via webhook (webhook URL configured via SLACK_WEBHOOK_URL environment variable)"),
		mcp.WithString("text", mcp.Required(), mcp.Description("The message text to send")),
		mcp.WithString("channel", mcp.Description("Optional: The channel to send the message to (if not specified in webhook)")),
		mcp.WithString("username", mcp.Description("Optional: Username to display as the sender")),
		mcp.WithString("icon_emoji", mcp.Description("Optional: Emoji to use as the icon (e.g., ':robot_face:')")),
		mcp.WithString("icon_url", mcp.Description("Optional: URL to an image to use as the icon")),
		mcp.WithBoolean("markdown", mcp.Description("Optional: Whether to enable markdown formatting (default: true)")),
	)

	// Register the tool
	mcpServer.AddTool(tool, s.handleSendMessage)

	s.server = mcpServer
	return s
}

// handleSendMessage handles the send_message tool execution
func (s *SlackMCPServer) handleSendMessage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get webhook URL from environment variable
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return mcp.NewToolResultError("SLACK_WEBHOOK_URL environment variable is required"), nil
	}

	// Extract required parameters
	text, err := request.RequireString("text")
	if err != nil {
		return mcp.NewToolResultError("text is required"), nil
	}

	// Create Slack message
	message := SlackMessage{
		Text:     text,
		Markdown: true, // Default to true
	}

	// Optional parameters
	if channel := request.GetString("channel", ""); channel != "" {
		message.Channel = channel
	}
	if username := request.GetString("username", ""); username != "" {
		message.Username = username
	}
	if iconEmoji := request.GetString("icon_emoji", ""); iconEmoji != "" {
		message.IconEmoji = iconEmoji
	}
	if iconURL := request.GetString("icon_url", ""); iconURL != "" {
		message.IconURL = iconURL
	}
	message.Markdown = request.GetBool("markdown", true)

	// Send message to Slack
	err = s.sendSlackMessage(webhookURL, message)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send message: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message sent successfully to Slack")), nil
}

// sendSlackMessage sends a message to Slack via webhook
func (s *SlackMCPServer) sendSlackMessage(webhookURL string, message SlackMessage) error {
	// Marshal message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Slack API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Run starts the MCP server
func (s *SlackMCPServer) Run() error {
	return server.ServeStdio(s.server)
}

func main() {
	// Create and start the server
	server := NewSlackMCPServer()

	log.Println("Starting Slack Webhook MCP Server...")
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
		os.Exit(1)
	}
}
