# 🪝 WebhookCatcher

**Debug, replay, and analyze webhooks locally**

A sleek local webhook debugging proxy that catches, stores, and lets you replay HTTP webhooks with a beautiful dashboard. Perfect for testing integrations with Stripe, GitHub, Slack, and any service that sends webhooks.

## Features

- **Catch Any Webhook** - One endpoint for all your webhooks
- **Beautiful Dashboard** - View, filter, and inspect captured events
- **Replay Support** - Forward captured webhooks to any URL
- **Persistent Storage** - Events saved to disk, survive restarts
- **Source Tagging** - Organize webhooks by source (github, stripe, slack, etc.)
- **Auto-refresh** - Dashboard updates in real-time

## Installation

```bash
# Download binary for your platform
curl -sSL https://github.com/YOUR_HANDLE/webhookcatcher/releases/latest | sh

# Or build from source
git clone https://github.com/YOUR_HANDLE/webhookcatcher.git
cd webhookcatcher
go build -o webhookcatcher .
```

## Usage

```bash
# Start the server
./webhookcatcher

# Or with custom port
PORT=8080 ./webhookcatcher
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/catch/{source}` | Catch incoming webhooks |
| GET | `/` | Dashboard |
| GET | `/api/events` | List all events |
| GET | `/api/events/{id}` | Get single event |
| POST | `/api/events/{id}/replay?url=...` | Replay to URL |
| DELETE | `/api/events/{id}` | Delete event |
| DELETE | `/api/clear` | Clear all events |

## Examples

```bash
# Catch GitHub webhooks
curl -X POST http://localhost:9876/catch/github \
  -H "Content-Type: application/json" \
  -d '{"action": "push", "repository": "my-app"}'

# Catch Stripe webhooks  
curl -X POST http://localhost:9876/catch/stripe \
  -d '{"type": "invoice.paid", "id": "in_123"}'

# Catch Slack webhooks
curl -X POST http://localhost:9876/catch/slack \
  -d '{"text": "Deployment completed!"}'

# Replay a webhook to your local server
curl -X POST "http://localhost:9876/api/events/{id}/replay?url=http://localhost:3000/webhook"

# List all events
curl http://localhost:9876/api/events | jq
```

## Dashboard

Open http://localhost:9876 to see the dashboard with:
- Real-time event count
- Click to expand event details
- Auto-refresh every 5 seconds
- Header and body inspection

## Use Cases

- **Stripe Development** - Test webhooks without Stripe CLI
- **GitHub Actions** - Debug CI/CD webhook payloads
- **Slack Apps** - Test Bolt/Slackkit integrations
- **Custom Integrations** - Debug any third-party webhook

## Tech Stack

- Go 1.19+
- gorilla/mux for routing
- UUID for event IDs
- Local JSON file storage

## License

MIT - Do whatever you want with it.
