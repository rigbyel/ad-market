package create

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rigbyel/ad-market/internal/lib/jwt"
	"github.com/rigbyel/ad-market/internal/lib/request"
	"github.com/rigbyel/ad-market/internal/lib/response"
	"github.com/rigbyel/ad-market/internal/lib/validate"
	"github.com/rigbyel/ad-market/internal/models"
)

type Response struct {
	response.Response
	Id          int64     `json:"id"`
	AuthorLogin string    `json:"author_id"`
	Date        time.Time `json:"date"`
}

type AdSaver interface {
	SaveAd(ad *models.Advert) (*models.Advert, error)
}

// New creates a new HandlerFunc for handling advert creation
func New(log *slog.Logger, adSaver AdSaver, authSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.advert.create.New"

		// setting up logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// authorization via jwt token
		tokenString := r.Header.Get("Authorization-access")
		tokenClaims, err := jwt.GetTokenClaims(tokenString, authSecret)
		if err != nil {
			log.Error("authorization failed", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("authorization failed"))

			return
		}

		// getting user's login
		login := tokenClaims.Login

		var req request.AdvertRequest

		// decoding request
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// validating advert data from user request
		validationErrs := validate.ValidateAdvert(req)
		if len(validationErrs) != 0 {
			log.Error("invalid request")

			render.JSON(w, r, response.Error(strings.Join(validationErrs, ", ")))

			return
		}

		// creating and saving advert
		ad := &models.Advert{
			Header:      req.Header,
			Body:        req.Body,
			ImageURL:    req.ImageURL,
			Price:       req.Price,
			Date:        time.Now(),
			AuthorLogin: login,
		}

		ad, err = adSaver.SaveAd(ad)
		if err != nil {
			log.Error("error saving advert", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("error saving advert"))

			return
		}

		log.Info("advert saved")

		render.JSON(w, r, Response{
			Response:    response.OK(),
			Id:          ad.Id,
			Date:        ad.Date,
			AuthorLogin: login,
		})
	}
}
