package httpdelivery

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func HealthHandler(pool *pgxpool.Pool, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		status := "ok"
		httpStatus := http.StatusOK

		if err := pool.Ping(ctx); err != nil {
			status = "unavailable"
			httpStatus = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)

		resp := map[string]string{"status": status}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error().Err(err).Msg("failed to encode health response")
		}
	}
}
