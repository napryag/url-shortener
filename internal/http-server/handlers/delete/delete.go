package delete

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/napryag/url-shortener/internal/lib/api/response"
	"github.com/napryag/url-shortener/internal/lib/logger/sl"
	"github.com/napryag/url-shortener/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2 --name=URLDeleter

type URLDeleter interface {
	DeleteURL(alias string) error
}

type Request struct {
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil && errors.Is(err, io.EOF) {
			alias := chi.URLParam(r, "alias")

			if alias == "" {
				log.Info("alias is empty")

				render.JSON(w, r, resp.Error("invalid request"))

				return
			}

			err := urlDeleter.DeleteURL(alias)
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", alias)

				render.JSON(w, r, resp.Error("not found"))

				return
			}
			if err != nil {
				log.Error("failed to delete url", sl.Err(err))

				render.JSON(w, r, resp.Error("internal error"))

				return
			}

			log.Info("url deleted correctly", slog.String("alias", alias))

			responseOK(w, r, alias)

		} else if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		} else {

			log.Info("request body decoded", slog.Any("req", req))

			err = urlDeleter.DeleteURL(req.Alias)
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", req.Alias)

				render.JSON(w, r, resp.Error("not found"))

				return
			}
			if err != nil {
				log.Error("failed to delete url", sl.Err(err))

				render.JSON(w, r, resp.Error("internal error"))

				return
			}

			log.Info("url deleted correctly", slog.String("alias", req.Alias))

			responseOK(w, r, req.Alias)
		}
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
