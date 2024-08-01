package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/fsnotify/fsnotify"
	"nhooyr.io/websocket"
)

func usage() error {
	return errors.New("usage: ./talk <talk number>")
}

func loadDirName(num string) (string, error) {
	targetDirNum, err := strconv.Atoi(num)
	if err != nil {
		return "", fmt.Errorf("input %q was not a valid number: %w", num, err)
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		return "", fmt.Errorf("failed to read current directory: %w", err)
	}

	for _, e := range entries {
		dirNumStr, _, found := strings.Cut(e.Name(), "-")
		if !found {
			continue
		}
		dirNum, err := strconv.Atoi(dirNumStr)
		if err != nil {
			continue
		}
		if dirNum == targetDirNum {
			return e.Name(), nil
		}
	}

	return "", fmt.Errorf("no dir found matching %q", num)
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type Metadata struct {
	Title string `json:"title"`
}

func loadMetadata(fn string) (*Metadata, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer f.Close() // Best-effort

	var md Metadata
	if err := json.NewDecoder(f).Decode(&md); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &md, nil
}

func run(args []string) error {
	if len(args) < 2 {
		return usage()
	}

	dir, err := loadDirName(os.Args[1])
	if err != nil {
		return fmt.Errorf("failed to find slides directory: %w", err)
	}

	md, err := loadMetadata(filepath.Join(dir, "metadata.json"))
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		return fmt.Errorf("failed to load index.html template: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errC := make(chan error)
	updateC := make(chan update)
	go func() {
		if err := reloadOnUpdate(ctx, updateC, dir); err != nil {
			errC <- fmt.Errorf("error during compilation: %w", err)
		}
	}()

	nextID := 1
	subs := make(map[int]chan update)
	go func() {
		for {
			// Wait for an update or cancellation
			select {
			case update := <-updateC:

				if update.Filename == "index.html" {
					time.Sleep(100 * time.Millisecond)
					tmpTmpl, err := template.ParseFiles("index.html")
					if err != nil {
						log.Printf("failed to reload index.html template: %v", err)
						continue
					}
					tmpl = tmpTmpl
				} else if update.Filename == "out.css" {
					time.Sleep(100 * time.Millisecond)
				}

				// Notify all active websockets
				for _, sub := range subs {
					select {
					case sub <- update:
						// Good
					case <-time.After(5 * time.Second):
						log.Println("failed to notify subscription after 5 seconds")
					}
				}
			case <-ctx.Done():
				errC <- ctx.Err()
				continue
			}

		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if err := tmpl.Execute(w, md); err != nil {
			log.Printf("failed to execute template: %v", err)
			return
		}
	})
	slidesPath := filepath.Join(dir, "slides.md")
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(dir, "assets")))))
	mux.HandleFunc("GET /slides.md", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		http.ServeFile(w, r, slidesPath)
	})
	mux.HandleFunc("GET /out.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		http.ServeFile(w, r, "out.css")
	})
	mux.HandleFunc("GET /remark.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "remark-latest.min.js")
	})
	mux.HandleFunc("GET /live", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Printf("failed to accept websocket conn: %v", err)
			return
		}
		sub := make(chan update)
		id := nextID
		subs[id] = sub
		nextID++

		defer c.CloseNow()
		defer func() {
			delete(subs, id)
			close(sub)
		}()

		ctx := c.CloseRead(r.Context())
	outer:
		for {
			select {
			case update := <-sub:
				if err := writeWithTimeout(ctx, time.Second*5, c, update); err != nil {
					log.Printf("write failed, closing websocket conn: %v", err)
					break outer
				}
			case <-ctx.Done():
				log.Printf("context done, closing websocket conn: %v", ctx.Err())
				break outer
			}
		}

		c.Close(websocket.StatusNormalClosure, "")
	})

	s := &http.Server{
		Addr:         ":8000",
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	return s.ListenAndServe()
}

func writeWithTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, ud update) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	dat, err := json.Marshal(ud)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %w", err)
	}

	return c.Write(ctx, websocket.MessageText, dat)
}

type update struct {
	Filename string `json:"filename"`
}

func reloadOnUpdate(ctx context.Context, updateC chan<- update, dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err = watcher.Add(dir); err != nil {
		return fmt.Errorf("failed to watch current dir: %w", err)
	}
	if err = watcher.Add("index.html"); err != nil {
		return fmt.Errorf("failed to watch index file: %w", err)
	}
	if err = watcher.Add("out.css"); err != nil {
		return fmt.Errorf("failed to watch CSS file: %w", err)
	}
	slidesPath := filepath.Join(dir, "slides.md")
	mdPath := filepath.Join(dir, "metadata.json")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("event channel closed")
			}
			if !(event.Has(fsnotify.Write)) {
				continue
			}
			switch n := event.Name; n {
			case "index.html", "out.css", slidesPath, mdPath:
				updateC <- update{Filename: event.Name}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("error channel closed")
			}
			log.Printf("error while watching: %v", err)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
