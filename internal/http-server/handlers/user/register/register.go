package register

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rigbyel/ad-market/internal/lib/jwt"
	"github.com/rigbyel/ad-market/internal/lib/request"
	"github.com/rigbyel/ad-market/internal/lib/response"
	"github.com/rigbyel/ad-market/internal/lib/validate"
	"github.com/rigbyel/ad-market/internal/models"
	"github.com/rigbyel/ad-market/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	response.Response
	ID int64 `json:"id"`
}

type UserSaver interface {
	SaveUser(u *models.User) (*models.User, error)
}

// New creates a new HandlerFunc to handle user registration
func New(log *slog.Logger, userSaver UserSaver, authSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.register.New"

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

		// validating login and password
		var validationErrs []string

		validationErrs = validate.ValidatePassword(req.Password)
		validationErrs = append(validationErrs, validate.ValidateLogin(req.Login)...)

		if len(validationErrs) != 0 {
			log.Error("invalid request data")

			render.JSON(w, r, response.Error(strings.Join(validationErrs, ", ")))

			return
		}

		// hashing password
		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password hash", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		// creating user
		user := &models.User{
			Login:    req.Login,
			PassHash: passHash,
		}

		// saving user in the storage
		user, err = userSaver.SaveUser(user)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("user", req.Login))

			render.JSON(w, r, response.Error("user already exists"))

			return
		}
		if err != nil {
			log.Error("error creating user", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("error creating user"))

			return
		}

		log.Info("user added", slog.Int64("id", user.Id))

		render.JSON(w, r,
			Response{
				Response: response.OK(),
				ID:       user.Id,
			},
		)
	}
}
