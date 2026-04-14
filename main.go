package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// WebhookEvent represents a captured webhook event
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Headers   map[string]string      `json:"headers"`
	Body      map[string]interface{} `json:"body"`
	RawBody   string                 `json:"raw_body"`
}

// Store for captured webhooks
var events []WebhookEvent
var storePath = ".webhookcatcher"

// Save event to disk
func saveEvent(event WebhookEvent) error {
	filename := filepath.Join(storePath, event.ID+".json")
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Load events from disk
func loadEvents() error {
	entries, err := os.ReadDir(storePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			data, err := os.ReadFile(filepath.Join(storePath, entry.Name()))
			if err != nil {
				continue
			}

			var event WebhookEvent
			if json.Unmarshal(data, &event) == nil {
				events = append(events, event)
			}
		}
	}
	return nil
}

// Webhook catcher endpoint
func catchWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", 500)
		return
	}

	// Parse headers
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, ", ")
	}

	// Parse body based on content type
	var bodyParsed map[string]interface{}
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "json") {
		json.Unmarshal(body, &bodyParsed)
	} else {
		bodyParsed = map[string]interface{}{
			"raw": string(body),
		}
	}

	// Create event
	event := WebhookEvent{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Source:    r.RemoteAddr,
		Headers:   headers,
		Body:      bodyParsed,
		RawBody:   string(body),
	}

	events = append(events, event)
	saveEvent(event)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "captured",
		"id":     event.ID,
		"total":  len(events),
	})
}

// List all events
func listEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Support filtering by source
	source := r.URL.Query().Get("source")
	
	var filtered []WebhookEvent
	if source != "" {
		for _, e := range events {
			if strings.Contains(e.Source, source) {
				filtered = append(filtered, e)
			}
		}
	} else {
		filtered = events
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(filtered),
		"events": filtered,
	})
}

// Get single event
func getEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, e := range events {
		if e.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(e)
			return
		}
	}

	http.Error(w, "Event not found", 404)
}

// Replay webhook to URL
func replayEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	targetURL := r.URL.Query().Get("url")

	if targetURL == "" {
		http.Error(w, "url parameter required", 400)
		return
	}

	var event WebhookEvent
	for _, e := range events {
		if e.ID == id {
			event = e
			break
		}
	}

	if event.ID == "" {
		http.Error(w, "Event not found", 404)
		return
	}

	// Forward the request
	req, _ := http.NewRequest("POST", targetURL, strings.NewReader(event.RawBody))
	for name, value := range event.Headers {
		if name != "Host" && name != "Content-Length" {
			req.Header.Set(name, value)
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Replay failed: %v", err), 500)
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, _ := io.ReadAll(resp.Body)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        resp.StatusCode,
		"response_body": string(respBody),
	})
}

// Delete event
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, e := range events {
		if e.ID == id {
			events = append(events[:i], events[i+1:]...)
			os.Remove(filepath.Join(storePath, id+".json"))
			w.WriteHeader(204)
			return
		}
	}

	http.Error(w, "Event not found", 404)
}

// Clear all events
func clearEvents(w http.ResponseWriter, r *http.Request) {
	events = nil
	os.RemoveAll(storePath)
	os.MkdirAll(storePath, 0755)
	w.WriteHeader(204)
}

// Dashboard HTML
func dashboard(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>WebhookCatcher</title>
	<style>
		body { font-family: -apple-system, sans-serif; background: #1a1a2e; color: #eee; margin: 0; padding: 20px; }
		.container { max-width: 1200px; margin: 0 auto; }
		h1 { color: #00d9ff; }
		.endpoint { background: #16213e; padding: 15px; margin: 10px 0; border-radius: 8px; }
		.endpoint code { background: #0f3460; padding: 2px 8px; border-radius: 4px; color: #00d9ff; }
		.event { background: #16213e; padding: 15px; margin: 10px 0; border-radius: 8px; cursor: pointer; }
		.event:hover { background: #1f2f50; }
		.pre { background: #0f3460; padding: 10px; border-radius: 4px; overflow-x: auto; white-space: pre-wrap; }
		.btn { background: #00d9ff; color: #1a1a2e; padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer; }
		.badge { background: #e94560; padding: 2px 8px; border-radius: 4px; font-size: 0.8em; }
	</style>
</head>
<body>
	<div class="container">
		<h1>🪝 WebhookCatcher</h1>
		<p>Your local webhook debugging proxy</p>

		<h2>Quick Endpoint</h2>
		<div class="endpoint">
			<code>POST http://localhost:%s/catch/{source}</code>
			<p>Catch webhooks and view them in the dashboard</p>
		</div>

		<h2>Recent Events <span class="badge">%d</span></h2>
		<div id="events">Loading...</div>
	</div>
	<script>
		async function loadEvents() {
			const resp = await fetch('/api/events');
			const data = await resp.json();
			const container = document.getElementById('events');
			if (data.events.length === 0) {
				container.innerHTML = '<p>No events captured yet!</p>';
				return;
			}
			container.innerHTML = data.events.reverse().map(e => '<div class="event"><code>' + e.id + '</code></div>').join('');
		}
		loadEvents();
		setInterval(loadEvents, 5000);
	</script>
</body>
</html>`, port, len(events))

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

var port = "9876"

func main() {
	// Create store directory
	os.MkdirAll(storePath, 0755)
	loadEvents()

	// Handle custom port
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	r := mux.NewRouter()

	// Dashboard
	r.HandleFunc("/", dashboard)
	r.HandleFunc("/catch/{source}", catchWebhook)
	r.HandleFunc("/api/events", listEvents)
	r.HandleFunc("/api/events/{id}", getEvent)
	r.HandleFunc("/api/events/{id}/replay", replayEvent)
	r.HandleFunc("/api/events/{id}", deleteEvent)
	r.HandleFunc("/api/clear", clearEvents)

	fmt.Printf(`
🪝 WebhookCatcher is running!

   Local:   http://localhost:%s
   Catch:   POST http://localhost:%s/catch/{source}

Examples:
   curl -X POST http://localhost:%s/catch/github -d '{"action":"push"}'
   curl -X POST http://localhost:%s/catch/stripe -d '{"type":"invoice.paid"}'
   curl -X POST http://localhost:%s/catch/slack -d '{"text":"Hello"}'

`, port, port, port, port, port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}