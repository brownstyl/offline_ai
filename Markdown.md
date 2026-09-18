# SMS-Query Backend Progress Log

## What This Project Is

SMS-Query is an offline AI search engine accessible via SMS. A user on any phone smartphone or feature phone sends a plain text message to a dedicated number. They need **zero internet data** to do this. Our backend receives that message, passes it to an AI engine, and sends a reply back as SMS. The user never touches the internet.

The backend is built in **Go**. The AI processing layer will be built in **Python**. The admin dashboard (showing live query logs) is built in **HTML/CSS/JS** by a separate team pair.

---

## Architecture Overview

```
[User's Phone —> no internet needed]
        │  plain SMS
        ▼
[Africa's Talking Gateway]
        │  converts SMS → HTTP POST (webhook)
        ▼
[ngrok tunnel —> development only]
        │  forwards public URL → localhost:8080
        ▼
[Go Backend —> :8080]
        │  extracts phone number + message
        │  calls Africa's Talking Send SMS API
        ▼
[Reply delivered back to user's phone]
```

In production, ngrok is replaced by a real server with a public IP. The Go backend and the gateway are the only parts that need internet the end user never does.

---

## File Structure

```
OFFLINE_AI/
└── backend/
    ├── main.go        — server entry point, routing, environment setup
    ├── handlers.go    — incoming SMS webhook handler
    ├── sendsms.go     — SmsSender interface + AfricasTalkingSender implementation
    ├── .env           — credentials (never committed to git)
    ├── go.mod         — Go module dependencies
    └── go.sum         — dependency checksums
```

---

## What Each File Does

### `main.go`
- Loads `.env` credentials using `godotenv.Load()` must run before anything else
- Reads `AT_USERNAME`, `AT_APIKEY`, `AT_URL` from environment via `os.Getenv`
- Creates an `AfricasTalkingSender` struct with those credentials
- Registers the `/incoming-sms` route using `http.NewServeMux()`
- Uses a **closure** to inject the sender into the handler without changing the handler's signature
- Starts the HTTP server on `:8080`

### `handlers.go`
- `HandleIncomingSms` receives the incoming webhook POST from Africa's Talking
- Guards against non-POST requests with an early return
- Extracts `from` (sender phone number) and `text` (message body) using `r.FormValue()`
- Calls `sender.Send()` to dispatch the reply — it only knows about the `SmsSender` interface, not Africa's Talking specifically
- Logs received data to terminal

### `sendsms.go`
- Defines the `SmsSender` interface with one method: `Send(to, message string) error`
- Defines `AfricasTalkingSender` struct holding `username`, `apikey`, and `url`
- `Send` method builds a form-encoded POST body using `url.Values{}`
- Wraps it with `strings.NewReader` so it satisfies `io.Reader` for `http.NewRequest`
- Sets required headers: `apiKey`, `Accept: application/json`, `Content-Type: application/x-www-form-urlencoded`
- Sends the request via `http.Client{}`
- Reads and logs the response body for debugging
- Returns descriptive errors for 401 (bad key), 400 (bad request), and other non-2xx statuses

### `.env`
```
AT_USERNAME=sandbox
AT_APIKEY=your_api_key_here
AT_URL=https://api.sandbox.africastalking.com/version1/messaging
```
This file is **gitignored** — never committed. Each teammate creates their own locally. In production, these values are injected as environment variables by the hosting platform — no file needed.

---

## Key Concepts Learned Building This

### Webhooks
Africa's Talking cannot reach `localhost` directly. When a user sends an SMS, AT converts it to an HTTP POST and sends it to a registered callback URL. During development, **ngrok** creates a public tunnel that forwards requests to `localhost:8080`. In production, a real server with a public IP replaces ngrok entirely.

### Why the Interface Exists
`handlers.go` accepts a `SmsSender` interface, not an `AfricasTalkingSender` directly. This means swapping Africa's Talking for a different provider (Android gateway, Twilio, etc.) requires:
- Writing a new struct that implements `Send(to, message string) error`
- Changing one line in `main.go` to pass the new struct

Zero changes to `handlers.go`. The handler doesn't know or care what's underneath — it only knows the behavior it needs.

### Environment Variables vs `.env` Files
`os.Getenv("KEY")` reads from the **shell environment**, not from files. A `.env` file on its own does nothing — `godotenv.Load()` bridges the gap by reading the file and loading each line into the shell environment at startup. In production, the hosting platform (Railway, Heroku, a VPS) injects these values directly into the process environment through a settings dashboard — no `.env` file ever touches the server.

### HTTP Status Codes Matter
Africa's Talking returns HTTP **201 Created** on a successful send, not 200 OK. The response body contains the real status. Our error handling checks for both `200` and `201` as success, and reads the body to surface meaningful error messages instead of silent failures.

### Form Encoding vs JSON
Africa's Talking's SMS API uses `application/x-www-form-urlencoded` POST bodies — the same format HTML forms use — not JSON. `url.Values{}.Encode()` handles encoding correctly, including safely escaping special characters like `+` in phone numbers as `%2B`. Never manually manipulate encoded strings with `strings.ReplaceAll` — it cannot distinguish between a `+` that means "space" and a `+` that is a real character.

---

## How to Run Locally

**Prerequisites:** Go installed, ngrok account (free), Africa's Talking sandbox account

```bash
# 1. Clone the repo
git clone <repo-url>
cd OFFLINE_AI/backend

# 2. Install dependencies
go mod tidy

# 3. Create your .env file (ask team lead for sandbox credentials)
cp .env.example .env
# fill in your values

# 4. Start the Go server
go run .

# 5. In a separate terminal, start ngrok
ngrok http 8080

# 6. Copy the ngrok forwarding URL (e.g. https://xxxx.ngrok-free.dev)
# Paste it into Africa's Talking sandbox:
# SMS → Callback URLs → Incoming Messages
# Full URL: https://xxxx.ngrok-free.dev/incoming-sms

# 7. Open the Africa's Talking simulator and send a test message to shortcode 90433
# Watch your Go terminal — you should see the phone number and message printed
# The simulator should receive a reply: "Message Received! Ai reply is coming soon..."
```

---

## What's Next

- [ ] Fix the 201 status check in `sendsms.go` so it doesn't log a false error on success
- [ ] Remove the debug `io.ReadAll` response body log from `sendsms.go` (or move it behind a debug flag)
- [ ] Wire up the Python AI engine — Go calls Python, Python queries the AI API, returns answer text
- [ ] Add a database layer to log every incoming message and outgoing reply (needed for the admin dashboard)
- [ ] Build the `/api/logs` endpoint that the frontend admin dashboard polls
- [ ] Replace hardcoded reply string with actual AI-generated response

---

## Environment Variables Reference

| Variable | Description | Sandbox Value |
|---|---|---|
| `AT_USERNAME` | Africa's Talking username | `sandbox` (literal) |
| `AT_APIKEY` | API key from AT dashboard Settings → API Key | generated per account |
| `AT_URL` | SMS send endpoint | `https://api.sandbox.africastalking.com/version1/messaging` |

> ⚠️ The sandbox and live environments use **separate API keys**. A sandbox key will always fail against the live URL and vice versa. Always confirm which environment your credentials belong to before debugging auth errors.