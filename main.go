package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var addr = flag.String("addr", ":8080", "The address to listen on")

func main() {
	flag.Parse()

	server := &http.Server{
		Addr:    *addr,
		Handler: http.DefaultServeMux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Server started")
		
		http.HandleFunc("/", homeOriginal)
		http.HandleFunc("/events", eventsOriginal)
		
		http.HandleFunc("/search-app", homeSearch)
		http.HandleFunc("/search-events", eventsSearch)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	sig := <-done
	slog.Info("Received signal, shutting down", slog.String("signal", sig.String()))

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error during server shutdown", slog.String("error", err.Error()))
	} else {
		slog.Info("Server shut down gracefully")
	}
}

func homeOriginal(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "typer.html")
}

func eventsOriginal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	alert := []string{"Alert:", "Someone", "just", "forgot", "a", "semicolon", "in", "production!"}

	for _, word := range alert {
		content := fmt.Sprintf("data: %s\n\n", word)
		w.Write([]byte(content))
		w.(http.Flusher).Flush()
		time.Sleep(time.Millisecond * 500)
	}

	reaction := []string{"I", "just", "deployed", "a", "week's", "work.", "Wait!", "Whattttt???", "Nooooo!"}

	for _, word := range reaction {
		content := fmt.Sprintf("data: %s\n\n", word)
		w.Write([]byte(content))
		w.(http.Flusher).Flush()
		time.Sleep(time.Millisecond * 500)
	}

	panicMsg := []string{"ERROR:", "The", "server", "is", "now", "in", "panic", "mode...", "Please", "consider", "rebooting", "your", "life...", "and", "the", "server 🖥️⚡💔!"}

	for _, word := range panicMsg {
		content := fmt.Sprintf("data: %s\n\n", word)
		w.Write([]byte(content))
		w.(http.Flusher).Flush()
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Println("Connection lost due to panic!")
}

func homeSearch(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "search.html")
}

func eventsSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		query = "nothing"
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	var response []string
	switch query {
	case "golang":
		response = []string{"Golang", "is", "a", "statically", "typed,", "compiled", "programming", "language", "designed", "at", "Google.", "It's", "great", "for", "concurrency!"}
	case "sse":
		response = []string{"Server-Sent", "Events", "(SSE)", "is", "a", "standard", "describing", "how", "servers", "can", "initiate", "data", "transmission", "towards", "clients", "once", "an", "initial", "client", "connection", "has", "been", "established."}
	case "hello":
		response = []string{"Hello", "there!", "I", "am", "a", "live", "search", "assistant", "streaming", "responses", "to", "you", "using", "Go", "and", "Server-Sent", "Events."}
	default:
		response = append([]string{"I", "searched", "my", "knowledge", "base", "for", fmt.Sprintf("'%s',", query), "but", "I", "couldn't", "find", "anything", "specific."}, "However,", "this", "demonstrates", "how", "I", "can", "stream", "any", "text", "back", "to", "you", "live!")
	}

	for _, word := range response {
		select {
		case <-ctx.Done():
			slog.Info("Client disconnected, stopping stream")
			return
		default:
			content := fmt.Sprintf("data: {\"word\": \"%s\"}\n\n", word)
			w.Write([]byte(content))
			w.(http.Flusher).Flush()
			time.Sleep(time.Millisecond * 200)
		}
	}

	finishMsg := "data: {\"done\": true}\n\n"
	w.Write([]byte(finishMsg))
	w.(http.Flusher).Flush()
}
