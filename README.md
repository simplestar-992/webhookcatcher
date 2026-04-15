# WebhookCatcher | Debug & Inspect Webhooks Locally

![Webhook Debugger](https://img.shields.io/badge/Purpose-Webhook%20Debugger-purple?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

---

## The Problem

You're integrating a third-party API that sends webhooks. How do you debug it locally? ngrok is complicated, cloud services cost money, and testing is a nightmare.

## The Solution

WebhookCatcher gives you an instant local endpoint to receive, inspect, and debug webhooks from any service.

### Perfect For

- **API Developers** - Debug webhook integrations locally
- **Third-party Integrations** - Stripe, GitHub, Slack, SendGrid, etc.
- **Testing** - Inspect payloads without complicated setup
- **Learning** - Understand how webhooks work

---

## Features

| Feature | Description |
|---------|-------------|
| **Instant Setup** | One command, get your endpoint |
| **Pretty JSON** | Formatted, syntax-highlighted payloads |
| **Request History** | See all received requests |
| **Response Customization** | Return whatever response you need |
| **Local Dashboard** | Web UI at localhost:8080 |
| **No Dependencies** | Single binary |

---

## Installation

```bash
git clone https://github.com/simplestar-992/webhookcatcher.git
cd webhookcatcher
go build -o webhookcatcher -ldflags="-s -w"
```

---

## Usage

### Start the Server

```bash
# Default (http://localhost:8080)
./webhookcatcher

# Custom port
./webhookcatcher -port 9090

# With custom response
./webhookcatcher -response '{"status":"received"}'
```

### Point Your Webhook

```
http://localhost:8080/any-endpoint-you-want
```

### Examples

```bash
# GitHub webhooks
./webhookcatcher -port 8080
# Point GitHub to: http://localhost:8080/github

# Stripe webhooks
./webhookcatcher -port 8080
# Point Stripe to: http://localhost:8080/stripe

# Slack webhooks
./webhookcatcher -port 8080
# Point Slack to: http://localhost:8080/slack
```

---

## Dashboard

Open http://localhost:8080 to see:
- All received requests
- Headers, body, method
- Timestamps
- Request count per endpoint

---

## Use Cases

### Debugging GitHub Actions
```bash
./webhookcatcher -port 8080
# In GitHub: Settings > Webhooks > Add webhook
# Payload URL: http://your-ip:8080/github
```

### Testing Stripe Integration
```bash
./webhookcatcher -port 8080
# In Stripe: Developers > Webhooks > Add endpoint
# URL: http://localhost:8080/stripe
```

---

## License

MIT © 2024 [simplestar-992](https://github.com/simplestar-992)
