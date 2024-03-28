package show

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rigbyel/ad-market/internal/lib/jwt"
	"github.com/rigbyel/ad-market/internal/lib/response"
	"github.com/rigbyel/ad-market/internal/models"
	"github.com/rigbyel/ad-market/internal/models/constraints"
)

type Advert struct {
	Header   string `json:"header"`
	Body     string `json:"body"`
	ImageURL string `json:"image_url,omitempty"`
	Price    int    `json:"price"`
	Author   string `json:"author"`
	IsAuthor bool   `json:"is_author,omitempty"`
}

type Response struct {
	response.Response
	Adverts *[]Advert `json:"adverts"`
}

type AdProvider interface {
	Adverts(minPrice, maxPrice int) (*[]models.Advert, error)
}

// New creates a new HandlerFunc for showing feed of adverts
func New(log *slog.Logger, adProv AdProvider, authSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.feed.show.New"

		// setting up logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var isAuthorized bool
		var login string

		// check if there's a valid token in the request
		tokenString := r.Header.Get("Authorization-access")
		tokenClaims, err := jwt.GetTokenClaims(tokenString, authSecret)
		if err == nil {
			log.Info("user authorized")

			login = tokenClaims.Login
			isAuthorized = true
		}

		// getting query parameters for min and max of advert prices
		// if there's no such parameters, set default values
		priceMin, _ := strconv.Atoi(r.URL.Query().Get("priceMin"))
		priceMax, _ := strconv.Atoi(r.URL.Query().Get("priceMax"))

		if priceMax == 0 {
			priceMax = constraints.MaxPrice
		}

		// get all adverts within given price range from storage
		adverts, err := adProv.Adverts(priceMin, priceMax)
		if err != nil {
			log.Error("failed to get adverts", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		// preparing adverts to show according to filters from query
		pageAdverts, err := prepareAdverts(r, adverts, isAuthorized, login)
		if err != nil {
			log.Error("no adverts found", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("nothing found"))

			return
		}

		log.Info("adverts accessed")

		render.JSON(w, r, Response{
			Response: response.OK(),
			Adverts:  &pageAdverts,
		})

	}
}

// prepares adverts according to filters from query parameters (sort type, sort order, page, etc)
func prepareAdverts(r *http.Request, adverts *[]models.Advert, isAuthorized bool, authorLogin string) ([]Advert, error) {
	const op = "handlers.feed.create.prepareAdverts"

	// get sorting type from query parameters
	sortType := r.URL.Query().Get("sort")
	if sortType == "" {
		sortType = "new"
	}

	// sort adverts with a given sorting type
	err := sortAdverts(*adverts, sortType)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// getting page of adverts feed from query parameters
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if page == 0 {
		page = 1
	}

	// check if there's any adverts on the given page
	if (page-1)*constraints.AdvertsOnPage >= len(*adverts) {
		return nil, fmt.Errorf("nothing found on page %d", page)
	}

	// extracting page from the whole feed of adverts
	var pageAdverts []Advert
	start := constraints.AdvertsOnPage * (page - 1)
	end := min(start+constraints.AdvertsOnPage, len(*adverts))

	for i := start; i < end; i++ {
		ad := (*adverts)[i]
		filteredAd := Advert{
			Header:   ad.Header,
			Body:     ad.Body,
			ImageURL: ad.ImageURL,
			Price:    ad.Price,
			Author:   ad.AuthorLogin,
		}

		if isAuthorized && ad.AuthorLogin == authorLogin {
			filteredAd.IsAuthor = true
		}

		pageAdverts = append(pageAdverts, filteredAd)
	}

	return pageAdverts, nil
}

// sorts adverts list with a given sorting type
func sortAdverts(adverts []models.Advert, sortType string) error {
	const op = "handlers.feed.show.sortAdverts"

	switch sortType {

	// ascending price
	case "priceUp":
		sort.Slice(adverts, func(i, j int) bool {
			return adverts[i].Price < adverts[j].Price
		})

	// descending price
	case "priceDown":
		sort.Slice(adverts, func(i, j int) bool {
			return adverts[i].Price > adverts[j].Price
		})

	// newest adverts in the beginning
	case "new":
		sort.Slice(adverts, func(i, j int) bool {
			return adverts[i].Date.Compare(adverts[j].Date) == 1
		})

	// oldest adverts in the begginning
	case "old":
		sort.Slice(adverts, func(i, j int) bool {
			return adverts[i].Date.Compare(adverts[j].Date) == -1
		})

	default:
		return fmt.Errorf("%s: wrong sorting parameter", op)
	}

	return nil
}
