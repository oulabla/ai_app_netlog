package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/oulabla/ai_app_netlog/internal/assets"
	"github.com/oulabla/ai_app_netlog/internal/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

// internal/server/swagger.go

func StartSwaggerServer(ctx context.Context) {
	addr := config.GetString(ctx, config.K.ServerSwaggerPort) // например ":8081"
	if addr == "" {
		log.Warn().
			Msg("Swagger port not configured, , skipping")
		return
	}

	mux := http.NewServeMux()

	// 1. Отдаём объединённый OpenAPI JSON
	mux.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		serveMergedSwaggerJSON(ctx, w)
	})

	// 2. Swagger UI на корневом пути (или /swagger, как вам удобнее)
	mux.Handle("/", httpSwagger.Handler(
		httpSwagger.URL("/openapi.json"),
		httpSwagger.DeepLinking(true),
		// httpSwagger.Prefix("/swagger"),   // если хотите путь /swagger/
	))
	log.Info().
		Str("addr", addr).
		Msg("starting HTTP Swagger UI server")

	if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
		log.Err(err).
			Str("addr", addr).
			Msg("Swagger listen error")
	}
}

func serveMergedSwaggerJSON(ctx context.Context, w http.ResponseWriter) {
	data, _ := assets.SwaggerFS.ReadFile("openapi/all-apis.swagger.json")

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		http.Error(w, "Invalid merged OpenAPI JSON", http.StatusInternalServerError)
		return
	}

	doc["host"] = config.GetString(ctx, config.K.ServerSwaggerHost)
	if doc["host"] == "" {
		doc["host"] = "localhost:8080"
	}

	doc["schemes"] = []string{"http"}

	modified, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		http.Error(w, "Cannot encode OpenAPI", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(modified)
	if err != nil {
		log.Err(err).Msg("write swagger")
	}
}
