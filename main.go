package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Global variable untuk menyimpan kata-kata agar tidak perlu split string setiap request (Hemat CPU & RAM)
var (
	textWords       []string
	textWordsNdjson []string
	loopText        []string
)

// Struct untuk format JSON response
type Chunk struct {
	Content string `json:"content"`
	Finish  bool   `json:"finish"` // Menggunakan bool true/false agar konsisten dengan JSON standar
}

func init() {
	// Pre-process text saat aplikasi start
	paragraph := "Large Language Models, often abbreviated as LLMs, are a type of artificial intelligence model trained on vast amounts of text data. They are designed to understand, generate, and respond to human language in a coherent and contextually relevant manner. This streaming demonstration mimics how an LLM might deliver its response token by token, providing a more interactive user experience rather than waiting for the entire output to be generated. Each word you see is a separate chunk of data sent from the server. This technique is crucial for applications that require real-time feedback, such as chatbots and live content generation."
	textWords = strings.Split(paragraph, " ")

	// Pre-process text NDJSON
	paragraphNdjson := "This NDJSON demonstration mimics how Ollama or other LLM APIs send responses as newline-delimited JSON objects. Each line you see is a chunk of data, separated by a newline, until the final finish signal is sent. Large Language Models, often abbreviated as LLMs, are a type of artificial intelligence model trained on vast amounts of text data. They are designed to understand, generate, and respond to human language in a coherent and contextually relevant manner. This streaming demonstration mimics how an LLM might deliver its response token by token, providing a more interactive user experience rather than waiting for the entire output to be generated."
	textWordsNdjson = strings.Split(paragraphNdjson, " ")

	loop := "The quick brown fox jumped over the lazy dog. The sun was shining brightly in the clear blue sky. A gentle breeze rustled the leaves of the trees as the birds sang their sweet melodies. In the distance, the sound of children's laughter echoed through the air."
	loopText = strings.Split(loop, " ")
}

// Middleware CORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// 1. STREAMING SSE
func streamSSE(w http.ResponseWriter, r *http.Request) {
	// Pastikan client mendukung streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i <= len(textWords); i++ {
		select {
		case <-r.Context().Done(): // Handle client disconnect
			fmt.Println("Client closed SSE connection")
			return
		case <-ticker.C:
			var chunk Chunk
			if i < len(textWords) {
				chunk = Chunk{Content: textWords[i] + " ", Finish: false}
			} else {
				chunk = Chunk{Content: "", Finish: true} // Null diganti empty string agar type-safe, finish true
			}

			jsonData, _ := json.Marshal(chunk)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush() // Kirim data segera ke client

			if chunk.Finish {
				return
			}
		}
	}
}

// 2. NDJSON ENDPOINT
func streamNDJSON(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i <= len(textWordsNdjson); i++ {
		select {
		case <-r.Context().Done():
			fmt.Println("Client closed NDJSON connection")
			return
		case <-ticker.C:
			var chunk Chunk
			if i < len(textWordsNdjson) {
				chunk = Chunk{Content: textWordsNdjson[i] + " ", Finish: false}
			} else {
				chunk = Chunk{Content: "", Finish: true}
			}

			json.NewEncoder(w).Encode(chunk) // Encode otomatis menambahkan newline
			flusher.Flush()

			if chunk.Finish {
				return
			}
		}
	}
}

// 3. STANDARD FETCH ENDPOINT
func standardData(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"id":      "a1b2-c3d4-e5f6",
		"type":    "static_response",
		"title":   "Complete Static Data",
		"message": strings.Join(textWords, " "),
		"author":  "Practice API (Go)",
		"metadata": map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
			"source":    "server-generated",
		},
		"payload": []map[string]interface{}{
			{"point": 1, "value": "First item"},
			{"point": 2, "value": "Second item"},
			{"point": 3, "value": "Third item"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// 4. LOOPING STREAM
func streamLoop(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			word := loopText[i%len(loopText)]
			fmt.Fprintf(w, "data: %s\n\n", word)
			flusher.Flush()
			i++
		}
	}
}

// 5. HEALTH CHECK
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"service": "LLM Practice API (Golang)",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		Practice API server is running.<br>
		Try accessing:<br>
		- /stream-sse<br>
		- /stream-ndjson<br>
		- /api/data<br>
		- /stream-loop
	`)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3950"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", enableCORS(rootHandler))
	mux.HandleFunc("/stream-sse", enableCORS(streamSSE))
	mux.HandleFunc("/stream-ndjson", enableCORS(streamNDJSON))
	mux.HandleFunc("/stream-loop", enableCORS(streamLoop))
	mux.HandleFunc("/api/data", enableCORS(standardData))
	mux.HandleFunc("/health", enableCORS(healthCheck))

	fmt.Printf("Server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
