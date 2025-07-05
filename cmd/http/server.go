package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"order-service/internal/usecase"
)

func Start(ctx context.Context, uc *usecase.OrderUC, addr string) error {
	log.Printf("HTTP API listening on %s  –  press Ctrl-C to stop", addr)

	r := chi.NewRouter()

	// Serve static files
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	r.Get("/order/{uid}", func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, "uid")
		log.Printf("REQUEST: GET /order/%s", uid)
		
		if uid == "" {
			log.Println("MISSING UID")
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		obj, err := uc.Get(r.Context(), uid)
		if err != nil {
			log.Printf("INTERNAL ERROR: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if obj == nil {
			log.Printf("NOT FOUND: %s", uid)
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		pretty, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			log.Printf("ENCODE ERROR: %v", err)
			http.Error(w, "encode error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Printf("RESPONSE OK: %s", uid)
		w.Write(pretty)
	})

	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return ctx.Err()
}
