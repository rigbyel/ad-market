package login

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rigbyel/ad-market/internal/lib/jwt"
	"github.com/rigbyel/ad-market/internal/lib/request"
	"github.com/rigbyel/ad-market/internal/lib/response"
	"github.com/rigbyel/ad-market/internal/models"
	"github.com/rigbyel/ad-market/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	response.Response
	Login string `json:"login"`
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type UserProvider interface {
	User(login string) (*models.User, error)
}

// New create a HandlerFunc to handle /login endpoint
func New(log *slog.Logger, userProvider UserProvider, authSecret string, tokenTL time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Login.New"

		// setting up logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// checking if user already authorized
		tokenString := r.Header.Get("Authorization-access")
		_, err := jwt.GetTokenClaims(tokenString, authSecret)
		if err == nil {
			log.Info("user already authorized")

			render.JSON(w, r, response.Error("user already authorized"))

			return
		}

		var req request.UserRequest

		// decoding request
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// getting user info from storage
		user, err := userProvider.User(req.Login)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found", slog.String("user", req.Login))

			render.JSON(w, r, response.Error("user not found"))

			return
		}
		if err != nil {
			log.Error("error finding user", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("error finding user"))

			return
		}

		// checking user's password
		if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(req.Password)); err != nil {
			log.Info("invalid credentials", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("invalid credentials"))

			return
		}

		log.Info("user logged in successfully")

		// creating new jwt token
		token, err := jwt.NewToken(*user, authSecret, tokenTL)
		if err != nil {
			log.Error("failed to generate token", slog.String("string", err.Error()))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		render.JSON(w, r,
			Response{
				Response: response.OK(),
				Login:    user.Login,
				Id:       user.Id,
				Token:    token,
			},
		)
	}
}
