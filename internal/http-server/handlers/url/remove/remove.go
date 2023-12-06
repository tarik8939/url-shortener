package remove

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/storage"
)

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLDelete interface {
	DeleteURL(alias string) (string, error)
}

func New(log *slog.Logger, urlSaver URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		msg, err := urlSaver.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get url"))

			return
		}

		log.Info("url added", slog.String("msg", msg))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(
		w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		},
	)
}
